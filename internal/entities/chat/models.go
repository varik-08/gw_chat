package chat

type Chat struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	IsPublic bool   `json:"isPublic"`
	OwnerID  int    `json:"ownerId"`
}

type DTO struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	IsPublic bool   `json:"isPublic"`
	OwnerID  int    `json:"ownerId"`
}

type Member struct {
	ChatID int `json:"chatId"`
	UserID int `json:"userId"`
}

type MemberDTO struct {
	ChatID int `json:"chatId"`
	UserID int `json:"userId"`
}
