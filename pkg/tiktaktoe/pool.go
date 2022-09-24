package tiktaktoe

import "fmt"

type Pool struct {
	Register   chan *Room
	Unregister chan *Room
	rooms      map[string]*Room
}

func NewPool() *Pool {
	return &Pool{
		Register:   make(chan *Room),
		Unregister: make(chan *Room),
		rooms:      make(map[string]*Room),
	}
}

func (p *Pool) Start() {
	for {
		select {
		case room := <-p.Register:
			{
				fmt.Printf("Creating new room with id '%s' \n", room.Id)
				p.rooms[room.Id] = room
				go room.Start()
			}
		case room := <-p.Unregister:
			{
				fmt.Printf("Closing room with id '%s' \n", room.Id)
				delete(p.rooms, room.Id)
			}
		}
	}
}

func (p Pool) GetRoom(id string) *Room {
	return p.rooms[id]
}
