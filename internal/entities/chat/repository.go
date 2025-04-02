package chat

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type RepositoryInterface interface {
	CreateChat(chat *Chat) (int, error)
	AddMember(member *Member) error
	GetChatByID(chatID int) (*Chat, error)
	CheckExistsUserInChat(userID int, chatID int) (bool, error)
	GetChatsByUserID(userID int) ([]Chat, error)
}

type Repository struct {
	db *pgxpool.Pool
}

func NewChatRepository(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

func (r *Repository) GetChatsByUserID(userID int) ([]Chat, error) {
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
		if err := rows.Scan(&chat.ID, &chat.Name, &chat.IsPublic, &chat.OwnerID); err != nil {
			return nil, err
		}
		chats = append(chats, chat)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return chats, nil
}

func (r *Repository) GetChatByID(chatID int) (*Chat, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	chat := &Chat{ID: chatID}
	err := r.db.QueryRow(ctx, "SELECT name, is_public, owner_id FROM chats WHERE id=$1", chatID).
		Scan(&chat.Name, &chat.IsPublic, &chat.OwnerID)
	if err != nil {
		return nil, err
	}

	return chat, nil
}

func (r *Repository) CheckExistsUserInChat(userID int, chatID int) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var exists bool

	err := r.db.QueryRow(ctx, "SELECT EXISTS (SELECT 1 FROM chat_user WHERE user_id=$1 and chat_id=$2)",
		userID, chatID).Scan(&exists)

	return exists, err
}

func (r *Repository) CreateChat(chat *Chat) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var chatID int

	query := `INSERT INTO chats (name, is_public, owner_id) VALUES ($1, $2, $3) RETURNING id`

	err := r.db.QueryRow(ctx, query, chat.Name, chat.IsPublic, chat.OwnerID).Scan(&chatID)

	return chatID, err
}

func (r *Repository) AddMember(member *Member) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `INSERT INTO chat_user (chat_id, user_id) VALUES ($1, $2)`
	_, err := r.db.Exec(ctx, query, member.ChatID, member.UserID)

	return err
}

func (r *Repository) GetMembersIDsByChatID(chatID int) ([]int, error) {
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
