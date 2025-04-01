package message

import "time"

type Message struct {
	ID        int       `json:"id"`
	ChatID    int       `json:"chatId"`
	UserID    int       `json:"userId"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"createdAt"`
	Username  string    `json:"username"`
}

type DTO struct {
	ID        int       `json:"id"`
	ChatID    int       `json:"chatId"`
	UserID    int       `json:"userId"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"createdAt"`
	Username  string    `json:"username"`
}
