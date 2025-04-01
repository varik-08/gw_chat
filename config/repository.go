package config

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/varik-08/gw_chat/internal/entities/chat"
	"github.com/varik-08/gw_chat/internal/entities/message"
	"github.com/varik-08/gw_chat/internal/entities/users"
)

type Repository struct {
	UserRepository    *users.UserRepository
	ChatRepository    *chat.Repository
	MessageRepository *message.Repository
}

func newRepository(db *pgxpool.Pool) *Repository {
	return &Repository{
		UserRepository:    users.NewUserRepository(db),
		ChatRepository:    chat.NewChatRepository(db),
		MessageRepository: message.NewMessageRepository(db),
	}
}
