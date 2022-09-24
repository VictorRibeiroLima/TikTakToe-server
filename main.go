package main

import (
	"fmt"
	"net/http"
	"strings"
	"tiktaktoe/pkg/tiktaktoe"

	"github.com/gorilla/websocket"
)

func createRoom(pool *tiktaktoe.Pool, w http.ResponseWriter, r *http.Request) {
	room := tiktaktoe.NewRoom(pool)
	conn, err := websocket.Upgrade(w, r, nil, 1024, 1024)
	if err != nil {
		w.WriteHeader(http.StatusBadGateway)
		fmt.Fprint(w, "Unable to stablish connection")
	}
	client := tiktaktoe.Client{Conn: conn, Id: conn.RemoteAddr().String(), Room: room}
	pool.Register <- room
	room.Register <- &client
	client.Read()
}

func joinRoom(pool *tiktaktoe.Pool, w http.ResponseWriter, r *http.Request) {
	roomId := strings.TrimPrefix(r.URL.Path, "/ws/join-room/")
	room := pool.GetRoom(roomId)
	conn, err := websocket.Upgrade(w, r, nil, 1024, 1024)
	if err != nil {
		w.WriteHeader(http.StatusBadGateway)
		fmt.Fprint(w, "Unable to stablish connection")
	}
	if room == nil {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, "Room not found")
		return
	}
	if !room.CanAddPlayer() {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "Room already full")
		return
	}
	client := tiktaktoe.Client{Conn: conn, Id: conn.RemoteAddr().String(), Room: room}
	room.Register <- &client
	client.Read()
}

func startServer() {
	pool := tiktaktoe.NewPool()
	go pool.Start()
	http.HandleFunc("/ws/create-room/", func(w http.ResponseWriter, r *http.Request) {
		createRoom(pool, w, r)
	})

	http.HandleFunc("/ws/join-room/", func(w http.ResponseWriter, r *http.Request) {
		joinRoom(pool, w, r)
	})

	http.ListenAndServe(":3000", nil)
}

func main() {
	startServer()
}
