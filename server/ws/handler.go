package ws

import (
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/varik-08/gw_chat/server/http/middlewares"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(_ *http.Request) bool { return true },
}

func WebSocketHandler(hub *Hub, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}

	client := &Client{
		Conn:         conn,
		UserID:       r.Context().Value(middlewares.UserIDKey).(int),
		Username:     r.Context().Value(middlewares.UsernameKey).(string),
		MessagesChan: make(chan wsMessageInterface),
	}
	hub.Register <- client

	go client.ReadMessages(hub)
	go client.WriteMessages()
}
