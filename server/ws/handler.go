package ws

import (
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func WebSocketHandler(hub *Hub, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}

	client := &Client{
		Conn:         conn,
		UserID:       r.Context().Value("userID").(int),
		Username:     r.Context().Value("username").(string),
		MessagesChan: make(chan wsMessageInterface),
	}
	hub.Register <- client

	go client.ReadMessages(hub)
	go client.WriteMessages()
}
