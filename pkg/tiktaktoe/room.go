package tiktaktoe

import (
	"errors"
	"fmt"
	"math/rand"
)

type Room struct {
	Id            string
	Register      chan *Client
	Unregister    chan *Client
	BroadCast     chan Payload
	pool          *Pool
	players       [2]*Client
	gameOn        bool
	currentPlayer *Client
	turn          int8
	game          *game
}

func NewRoom(pool *Pool) *Room {
	//id := uuid.New()
	return &Room{
		Id:         "test",
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		BroadCast:  make(chan Payload),
		pool:       pool,
		game:       &game{},
	}
}

func (r *Room) addPlayer(c *Client) error {
	if r.players[0] == nil {
		r.players[0] = c
		return nil
	} else if r.players[1] == nil {
		r.players[1] = c
		return nil
	}
	return errors.New("Room already full")
}

func (r *Room) removePlayer(c *Client) error {
	if r.players[0] == c {
		r.players[0] = nil
		return nil
	} else if r.players[1] == c {
		r.players[1] = nil
		return nil
	}
	return errors.New("Player not on this room")
}

func (r *Room) startGame() {
	if !r.CanAddPlayer() {
		randPlayer := rand.Intn(2)
		fmt.Println(randPlayer)
		r.currentPlayer = r.players[randPlayer]
		otherPlayer := r.getOtherPlayer(r.currentPlayer)
		r.currentPlayer.Conn.WriteJSON(message{Event: "GAME_START", Message: "Starting game"})
		otherPlayer.Conn.WriteJSON(message{Event: "GAME_START", Message: "Starting game"})
		r.emitTurn()
	}
}

func (r Room) getOtherPlayer(c *Client) *Client {
	if r.players[0] == c {
		return r.players[1]
	}
	return r.players[0]
}

func (r Room) CanAddPlayer() bool {
	return r.players[0] == nil || r.players[1] == nil
}

func (r Room) emitTurn() {
	current := r.currentPlayer
	otherPlayer := r.getOtherPlayer(current)
	current.Conn.WriteJSON(message{Event: "YOUR_TURN", Message: "Is your turn"})
	otherPlayer.Conn.WriteJSON(message{Event: "OPPONENT_TURN", Message: "Is your opponent turn"})
}

func (r *Room) changeTurn() {
	current := r.currentPlayer
	otherPlayer := r.getOtherPlayer(current)
	r.currentPlayer = otherPlayer
}

func (r *Room) play(payload Payload) {
	event := payload.Message.Event
	movement := payload.Message.Movement
	sender := payload.From
	otherPlayer := r.getOtherPlayer(sender)
	if event != "PLAY" {
		if sender != nil {
			sender.Conn.WriteJSON(message{Event: "ERROR", Message: "Event '" + event + "' not supported"})
		}
		return
	}
	if otherPlayer == nil {
		sender.Conn.WriteJSON(message{Event: "ERROR", Message: "Event not enough players"})
		return
	}
	if sender != r.currentPlayer {
		sender.Conn.WriteJSON(message{Event: "ERROR", Message: "Not your turn"})
		return
	}
	if movement.Row < 0 || movement.Row > 2 || movement.Column < 0 || movement.Column > 2 {
		sender.Conn.WriteJSON(message{Event: "ERROR", Message: "Invalid Movement"})
		return
	}
	result, err := r.game.MakePlay(movement.Row, movement.Column)
	if err != nil {
		sender.Conn.WriteJSON(message{Event: "ERROR", Message: err.Error()})
	}
	otherPlayer.Conn.WriteJSON(message{Event: "MOVEMENT", Message: "Opponent moved", Movement: movement})
	r.game.Draw()
	if result == 1 {
		sender.Conn.WriteJSON(message{Event: "WIN", Message: "YOU WIN"})
		otherPlayer.Conn.WriteJSON(message{Event: "LOST", Message: "YOU LOUSE"})
		r.startGame()
	} else if result == 2 {
		sender.Conn.WriteJSON(message{Event: "DRAW", Message: "YOU DRAW"})
		otherPlayer.Conn.WriteJSON(message{Event: "DRAW", Message: "YOU DRAW"})
		r.startGame()
	} else {
		r.changeTurn()
		r.emitTurn()
	}
}

func (r *Room) Start() {
	for {
		select {
		case player := <-r.Register:
			{
				err := r.addPlayer(player)
				if err != nil {
					player.Conn.WriteJSON(message{Event: "ERROR", Message: err.Error()})
				} else {
					player.Conn.WriteJSON(message{Event: "ROOM_CONNECTION", Message: r.Id})
					r.startGame()
				}
			}
		case player := <-r.Unregister:
			{
				r.removePlayer(player)
				r.game = &game{}
				if r.players[0] == nil && r.players[1] == nil {
					r.pool.Unregister <- r
					return
				}
				otherPlayer := r.getOtherPlayer(player)
				if otherPlayer != nil {
					otherPlayer.Conn.WriteJSON(message{Event: "GAME_STOP", Message: "Player disconnected"})
				}
			}
		case payload := <-r.BroadCast:
			{
				r.play(payload)
			}
		}
	}

}
