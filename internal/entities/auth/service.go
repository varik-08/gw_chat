package auth

import (
	"fmt"
	"github.com/varik-08/gw_chat/internal/entities/users"
	"github.com/varik-08/gw_chat/pkg"
	"github.com/varik-08/gw_chat/pkg/jwt"
)

type AuthService struct {
	UserRepository *users.UserRepository
	secretKey      string
}

func NewAuthService(secretKey string, userRepository *users.UserRepository) *AuthService {
	return &AuthService{
		secretKey:      secretKey,
		UserRepository: userRepository,
	}
}

func (s *AuthService) Login(credentials users.Credentials) (*Token, error) {
	user, err := s.UserRepository.GetUserByUsername(credentials.Username)
	if err != nil {
		return nil, fmt.Errorf("ошибка получения пользователя")
	}
	if user == nil {
		return nil, fmt.Errorf("пользователь не найден")
	}

	if !pkg.CheckPasswordHash(credentials.Password, user.PasswordHash) {
		return nil, fmt.Errorf("пользователь не найден")
	}

	accessToken, expiresAt, err := jwt.GenerateAccessToken(s.secretKey, user.ID, user.Username)
	if err != nil {
		return nil, fmt.Errorf("ошибка генерации access token: %w", err)
	}

	refreshToken, _, err := jwt.GenerateRefreshToken(s.secretKey, user.ID, user.Username)
	if err != nil {
		return nil, fmt.Errorf("ошибка генерации refresh token: %w", err)
	}

	token := Token{
		UserID:       user.ID,
		Username:     user.Username,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresAt:    expiresAt,
	}

	return &token, nil
}

func (s *AuthService) RefreshToken(oldToken *Token) (*Token, error) {
	// Валидация refresh токена
	claims, err := jwt.ValidateToken(s.secretKey, oldToken.RefreshToken)
	if err != nil {
		return nil, fmt.Errorf("ошибка валидации refresh token: %w", err)
	}

	// Генерация нового access token
	newAccessToken, _, err := jwt.GenerateAccessToken(s.secretKey, claims.UserID, claims.Username)
	if err != nil {
		return nil, fmt.Errorf("ошибка генерации access token: %w", err)
	}

	// Генерация нового refresh token
	newRefreshToken, _, err := jwt.GenerateRefreshToken(s.secretKey, claims.UserID, claims.Username)
	if err != nil {
		return nil, fmt.Errorf("ошибка генерации refresh token: %w", err)
	}

	token := Token{
		AccessToken:  newAccessToken,
		RefreshToken: newRefreshToken,
	}

	return &token, nil
}
