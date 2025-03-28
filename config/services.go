package config

import (
	"github.com/varik-08/gw_chat/internal/entities/auth"
	"github.com/varik-08/gw_chat/internal/entities/chat"
	"github.com/varik-08/gw_chat/internal/entities/message"
	"github.com/varik-08/gw_chat/internal/entities/users"
)

type Service struct {
	UserService    *users.UserService
	AuthService    *auth.AuthService
	ChatService    *chat.ChatService
	MessageService *message.MessageService
}

func newService(conf *Cfg, repositories *Repository) *Service {
	return &Service{
		UserService: users.NewUserService(repositories.UserRepository),
		AuthService: auth.NewAuthService(
			conf.SecretKey,
			repositories.UserRepository,
		),
		ChatService:    chat.NewChatService(repositories.ChatRepository),
		MessageService: message.NewChatService(repositories.MessageRepository, repositories.ChatRepository),
	}
}
