package tiktaktoe

import (
	"encoding/json"
	"fmt"

	"github.com/gorilla/websocket"
)

type Client struct {
	Id   string
	Conn *websocket.Conn
	Room *Room
}

func (c *Client) Read() {
	defer func() {
		c.Room.Unregister <- c
		c.Conn.Close()
	}()

	for {
		_, p, err := c.Conn.ReadMessage()
		if err != nil {
			fmt.Println(err)
			return
		}

		var msg message
		err = json.Unmarshal(p, &msg)
		if err != nil {
			c.Conn.WriteJSON(message{Event: "ERROR", Message: "Invalid payload"})
		}
		payload := Payload{From: c, Message: msg}
		c.Room.BroadCast <- payload
	}
}
