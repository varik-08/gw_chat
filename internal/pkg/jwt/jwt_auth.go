package jwt

import (
	"errors"
	"time"

	jwt "github.com/golang-jwt/jwt/v5"
)

const (
	AccessTokenCountMinutes = 180
	RefreshTokenCountDays   = 7
)

type Claims struct {
	UserID   int    `json:"userId"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

func GenerateAccessToken(secretKey string, userID int, username string) (string, time.Time, error) {
	jwtSecret := []byte(secretKey)
	expirationTime := time.Now().Add(AccessTokenCountMinutes * time.Minute)

	claims := &Claims{
		UserID:   userID,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		return "", time.Time{}, err
	}

	return tokenString, expirationTime, nil
}

func GenerateRefreshToken(secretKey string, userID int, username string) (string, time.Time, error) {
	jwtSecret := []byte(secretKey)
	expirationTime := time.Now().Add(RefreshTokenCountDays * 24 * time.Hour)

	claims := &Claims{
		UserID:   userID,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		return "", time.Time{}, err
	}

	return tokenString, expirationTime, nil
}

func ValidateToken(secretKey string, tokenString string) (*Claims, error) {
	jwtSecret := []byte(secretKey)

	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(_ *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}
