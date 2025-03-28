package ws

import "github.com/varik-08/gw_chat/internal/entities/message"

type wsMessageInterface interface {
	GetType() string
}

type wsChatMessage struct {
	Type    string          `json:"type"`
	Message message.Message `json:"message"`
}

func (m *wsChatMessage) GetType() string {
	return m.Type
}

type wsTypingMessage struct {
	Type     string `json:"type"`
	ChatID   int    `json:"chat_id"`
	UserID   int    `json:"user_id"`
	UserName string `json:"username"`
	IsTyping bool   `json:"is_typing"`
}

func (m *wsTypingMessage) GetType() string {
	return m.Type
}

type wsUserStatusMessage struct {
	Type        string       `json:"type"`
	ActiveUsers map[int]bool `json:"active_users"`
}

func (m *wsUserStatusMessage) GetType() string {
	return m.Type
}
