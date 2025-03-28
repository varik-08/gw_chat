package chat

type Chat struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	IsPublic bool   `json:"is_public"`
	OwnerId  int    `json:"owner_id"`
}

type ChatDTO struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	IsPublic bool   `json:"is_public"`
	OwnerId  int    `json:"owner_id"`
}

// ChatMember связывает пользователей с чатами
type ChatMember struct {
	ChatID int `json:"chat_id"`
	UserID int `json:"user_id"`
}

type ChatMemberDTO struct {
	ChatID int `json:"chat_id"`
	UserID int `json:"user_id"`
}
