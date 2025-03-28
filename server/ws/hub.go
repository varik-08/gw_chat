package ws

import (
	"fmt"
	"github.com/varik-08/gw_chat/config"
	"github.com/varik-08/gw_chat/internal/entities/message"
	"sync"
)

type Hub struct {
	Clients     map[*Client]bool
	UserClients map[int][]*Client
	Broadcast   chan wsMessageInterface
	Register    chan *Client
	Unregister  chan *Client
	ActiveUsers map[int]bool
	TypingUsers map[int]int
	mu          sync.Mutex
	app         *config.App
}

func NewHub(app *config.App) *Hub {
	return &Hub{
		Clients:     make(map[*Client]bool),
		UserClients: make(map[int][]*Client),
		Broadcast:   make(chan wsMessageInterface),
		Register:    make(chan *Client),
		Unregister:  make(chan *Client),
		ActiveUsers: make(map[int]bool),
		TypingUsers: make(map[int]int),
		app:         app,
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			h.mu.Lock()

			h.Clients[client] = true
			h.UserClients[client.UserID] = append(h.UserClients[client.UserID], client)
			h.ActiveUsers[client.UserID] = true

			h.mu.Unlock()

			h.BroadcastUserStatus()

		case client := <-h.Unregister:
			h.mu.Lock()

			delete(h.Clients, client)
			delete(h.ActiveUsers, client.UserID)

			h.mu.Unlock()

			h.BroadcastUserStatus()

		case mes := <-h.Broadcast:
			h.mu.Lock()
			for client := range h.Clients {
				client.MessagesChan <- mes
			}
			h.mu.Unlock()
		}
	}
}

func (h *Hub) BroadcastUserStatus() {
	h.mu.Lock()
	defer h.mu.Unlock()

	statusMessage := wsUserStatusMessage{
		Type:        "user_status",
		ActiveUsers: h.ActiveUsers,
	}

	// Отправляем статус всем клиентам
	for client := range h.Clients {
		client.MessagesChan <- &statusMessage
	}
}

func (h *Hub) BroadcastTyping(wsTypingMessage *wsTypingMessage) {
	h.mu.Lock()
	if wsTypingMessage.IsTyping {
		h.TypingUsers[wsTypingMessage.UserID] = wsTypingMessage.ChatID
	} else {
		delete(h.TypingUsers, wsTypingMessage.UserID)
	}
	h.mu.Unlock()

	userIDs, err := h.app.Repositories.ChatRepository.GetMembersIDsByChatID(wsTypingMessage.ChatID)
	if err != nil {
		fmt.Println("Ошибка при получении участников чата:", err)
		return
	}

	for _, userID := range userIDs {
		if clients, ok := h.UserClients[userID]; ok {
			for _, client := range clients {
				client.MessagesChan <- wsTypingMessage
			}
		}
	}
}

func (h *Hub) BroadcastMessage(msg *wsChatMessage) {
	var messageDTO message.MessageDTO
	messageDTO.ChatID = msg.Message.ChatID
	messageDTO.UserID = msg.Message.UserID
	messageDTO.Content = msg.Message.Content

	messageID, err := h.app.Services.MessageService.CreateChatMessage(&messageDTO)
	if err != nil {
		fmt.Println("Ошибка при создании сообщения:", err)
		return
	}

	msg.Message.ID = messageID

	userIDs, err := h.app.Repositories.ChatRepository.GetMembersIDsByChatID(msg.Message.ChatID)
	if err != nil {
		fmt.Println("Ошибка при получении участников чата:", err)
		return
	}

	for _, userID := range userIDs {
		if clients, ok := h.UserClients[userID]; ok {
			for _, client := range clients {
				client.MessagesChan <- msg
			}
		}
	}
}
