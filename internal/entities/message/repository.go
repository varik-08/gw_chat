package message

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type RepositoryInterface interface {
	GetMessagesByChatID(chatID int) ([]*Message, error)
	CreateChatMessage(message *Message) (int, error)
}

type Repository struct {
	db *pgxpool.Pool
}

func NewMessageRepository(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

func (r *Repository) GetMessagesByChatID(chatID int) ([]*Message, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	rows, err := r.db.Query(ctx, `
SELECT distinct m.id, m.chat_id, m.user_id, m.text, m.created_at, u.username 
FROM messages m
INNER JOIN chat_user cu on cu.chat_id = m.chat_id
INNER JOIN users u on u.id = m.user_id
WHERE cu.chat_id=$1
ORDER BY m.created_at
`, chatID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	messages := make([]*Message, 0)
	for rows.Next() {
		var message Message
		if err := rows.Scan(&message.ID, &message.ChatID, &message.UserID,
			&message.Content, &message.CreatedAt, &message.Username); err != nil {
			return nil, err
		}
		messages = append(messages, &message)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return messages, nil
}

func (r *Repository) CreateChatMessage(message *Message) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var messageID int

	query := `INSERT INTO messages (chat_id, user_id, text) VALUES ($1, $2, $3) RETURNING id`

	err := r.db.QueryRow(ctx, query, message.ChatID, message.UserID, message.Content).Scan(&messageID)

	return messageID, err
}
