package config

import (
	"github.com/varik-08/gw_chat/internal/entities/auth"
	"github.com/varik-08/gw_chat/internal/entities/chat"
	"github.com/varik-08/gw_chat/internal/entities/message"
	"github.com/varik-08/gw_chat/internal/entities/users"
)

type Service struct {
	UserService    *users.UserService
	AuthService    *auth.Service
	ChatService    *chat.Service
	MessageService *message.Service
}

func newService(conf *Cfg, repositories *Repository) *Service {
	return &Service{
		UserService: users.NewUserService(repositories.UserRepository),
		AuthService: auth.NewAuthService(
			conf.SecretKey,
			repositories.UserRepository,
		),
		ChatService:    chat.NewChatService(repositories.ChatRepository),
		MessageService: message.NewMessageService(repositories.MessageRepository, repositories.ChatRepository),
	}
}
