package ws

import (
	"encoding/json"
	"github.com/gorilla/websocket"
)

type Client struct {
	Conn         *websocket.Conn
	UserID       int
	Username     string
	MessagesChan chan wsMessageInterface
}

func (c *Client) ReadMessages(hub *Hub) {
	defer func() {
		hub.Unregister <- c
		c.Conn.Close()
	}()

	for {
		var rawMsg map[string]interface{}
		err := c.Conn.ReadJSON(&rawMsg)
		if err != nil {
			break
		}

		jsonData, err := json.Marshal(rawMsg)
		if err != nil {
			continue
		}

		messageType, ok := rawMsg["type"].(string)
		if !ok {
			continue
		}

		switch messageType {
		case "message":
			var msg wsChatMessage
			msg.Type = messageType
			msg.Message.Username = c.Username
			msg.Message.UserID = c.UserID
			msg.Message.Content = rawMsg["content"].(string)
			msg.Message.ChatID = int(rawMsg["chat_id"].(float64))

			hub.BroadcastMessage(&msg)

		case "typing":
			var msg wsTypingMessage
			if err := json.Unmarshal(jsonData, &msg); err != nil {
				continue
			}

			typingMessage := wsTypingMessage{
				Type:     "typing",
				UserID:   c.UserID,
				UserName: c.Username,
				ChatID:   msg.ChatID,
				IsTyping: msg.IsTyping,
			}

			hub.BroadcastTyping(&typingMessage)
		default:
			continue
		}
	}
}

func (c *Client) WriteMessages() {
	for msg := range c.MessagesChan {
		c.Conn.WriteJSON(msg)
	}
}
