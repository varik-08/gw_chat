package chat

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"time"
)

type ChatRepository struct {
	db *pgxpool.Pool
}

func NewChatRepository(db *pgxpool.Pool) *ChatRepository {
	return &ChatRepository{db: db}
}

func (r *ChatRepository) GetChatsByUserID(userID int) ([]Chat, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	rows, err := r.db.Query(ctx, `
		SELECT c.id, c.name, c.is_public, c.owner_id 
		FROM chats c
		JOIN chat_user cu ON c.id = cu.chat_id
		WHERE cu.user_id = $1
		ORDER BY c.id DESC
		`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var chats []Chat
	for rows.Next() {
		var chat Chat
		if err := rows.Scan(&chat.ID, &chat.Name, &chat.IsPublic, &chat.OwnerId); err != nil {
			return nil, err
		}
		chats = append(chats, chat)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return chats, nil
}

func (r *ChatRepository) GetChatByID(chatID int) (*Chat, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	chat := &Chat{ID: chatID}
	err := r.db.QueryRow(ctx, "SELECT name, is_public, owner_id FROM chats WHERE id=$1", chatID).
		Scan(&chat.Name, &chat.IsPublic, &chat.OwnerId)

	if err != nil {
		return nil, err
	}

	return chat, nil
}

func (r *ChatRepository) CheckExistsUserInChat(userID int, chatID int) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var exists bool

	err := r.db.QueryRow(ctx, "SELECT EXISTS (SELECT 1 FROM chat_user WHERE user_id=$1 and chat_id=$2)", userID, chatID).Scan(&exists)

	return exists, err
}

func (r *ChatRepository) CreateChat(chat *Chat) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var chatID int

	query := `INSERT INTO chats (name, is_public, owner_id) VALUES ($1, $2, $3) RETURNING id`

	err := r.db.QueryRow(ctx, query, chat.Name, chat.IsPublic, chat.OwnerId).Scan(&chatID)

	return chatID, err
}

func (r *ChatRepository) AddMember(member *ChatMember) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `INSERT INTO chat_user (chat_id, user_id) VALUES ($1, $2)`
	_, err := r.db.Exec(ctx, query, member.ChatID, member.UserID)

	return err
}

func (r *ChatRepository) GetMembersIDsByChatID(chatID int) ([]int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	rows, err := r.db.Query(ctx, "SELECT user_id FROM chat_user WHERE chat_id=$1", chatID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ids []int
	for rows.Next() {
		var id int
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return ids, nil
}
