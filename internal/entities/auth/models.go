package auth

import "time"

type Token struct {
	UserID       int       `json:"userId"`
	Username     string    `json:"username"`
	AccessToken  string    `json:"accessToken"`
	RefreshToken string    `json:"refreshToken"`
	ExpiresAt    time.Time `json:"expiresAt"`
}
