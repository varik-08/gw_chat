package auth

import (
	"crypto/rsa"
	"fmt"
	"time"

	gjwt "github.com/golang-jwt/jwt/v5"
	"github.com/varik-08/gw_chat/internal/entities/users"
	"github.com/varik-08/gw_chat/internal/pkg"
	jwtHS "github.com/varik-08/gw_chat/internal/pkg/jwt"
	jwtRS "github.com/varik-08/gw_chat/internal/security/jwt"
)

type Service struct {
	UserRepository users.UserRepositoryInterface
	secretKey      string
	rsaPriv        *rsa.PrivateKey
	kid            string
}

func NewAuthService(secretKey string, userRepository users.UserRepositoryInterface) *Service {
	return &Service{
		secretKey:      secretKey,
		UserRepository: userRepository,
	}
}

func NewAuthServiceRS256(priv *rsa.PrivateKey, kid string, userRepository users.UserRepositoryInterface) *Service {
	return &Service{
		rsaPriv:        priv,
		kid:            kid,
		UserRepository: userRepository,
	}
}

func (s *Service) Login(credentials users.Credentials) (*Token, error) {
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

	var accessToken, refreshToken string
	var expiresAt time.Time
	if s.rsaPriv != nil {
		accessToken, expiresAt, err = jwtRS.SignTokenRS256(s.rsaPriv, s.kid, int64(user.ID), user.Username, 15*time.Minute)
		if err != nil { return nil, fmt.Errorf("ошибка генерации access token: %w", err) }
		refreshToken, _, err = jwtRS.SignTokenRS256(s.rsaPriv, s.kid, int64(user.ID), user.Username, 7*24*time.Hour)
		if err != nil { return nil, fmt.Errorf("ошибка генерации refresh token: %w", err) }
	} else {
		accessToken, expiresAt, err = jwtHS.GenerateAccessToken(s.secretKey, user.ID, user.Username)
		if err != nil { return nil, fmt.Errorf("ошибка генерации access token: %w", err) }
		refreshToken, _, err = jwtHS.GenerateRefreshToken(s.secretKey, user.ID, user.Username)
		if err != nil { return nil, fmt.Errorf("ошибка генерации refresh token: %w", err) }
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

func (s *Service) RefreshToken(oldToken *Token) (*Token, error) {
	if s.rsaPriv != nil {
		t, err := gjwt.ParseWithClaims(oldToken.RefreshToken, &jwtRS.Claims{}, func(_ *gjwt.Token) (interface{}, error) {
			return &s.rsaPriv.PublicKey, nil
		})
		if err != nil || !t.Valid { return nil, fmt.Errorf("ошибка валидации refresh token: %w", err) }
		claims := t.Claims.(*jwtRS.Claims)
		newAcc, _, err := jwtRS.SignTokenRS256(s.rsaPriv, s.kid, claims.UserID, claims.Username, 15*time.Minute)
		if err != nil { return nil, fmt.Errorf("ошибка генерации access token: %w", err) }
		newRef, _, err := jwtRS.SignTokenRS256(s.rsaPriv, s.kid, claims.UserID, claims.Username, 7*24*time.Hour)
		if err != nil { return nil, fmt.Errorf("ошибка генерации refresh token: %w", err) }
		return &Token{AccessToken: newAcc, RefreshToken: newRef}, nil
	}

	claims, err := jwtHS.ValidateToken(s.secretKey, oldToken.RefreshToken)
	if err != nil {
		return nil, fmt.Errorf("ошибка валидации refresh token: %w", err)
	}

	newAccessToken, _, err := jwtHS.GenerateAccessToken(s.secretKey, claims.UserID, claims.Username)
	if err != nil {
		return nil, fmt.Errorf("ошибка генерации access token: %w", err)
	}

	newRefreshToken, _, err := jwtHS.GenerateRefreshToken(s.secretKey, claims.UserID, claims.Username)
	if err != nil {
		return nil, fmt.Errorf("ошибка генерации refresh token: %w", err)
	}

	token := Token{
		AccessToken:  newAccessToken,
		RefreshToken: newRefreshToken,
	}

	return &token, nil
}
