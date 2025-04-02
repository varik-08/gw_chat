package jwt

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const (
	secretKey = "supersecretkey"
	userID    = 123
	username  = "testuser"
)

func TestGenerateAccessToken(t *testing.T) {
	token, expirationTime, err := GenerateAccessToken(secretKey, userID, username)

	assert.NoError(t, err)

	assert.NotEmpty(t, token)

	assert.True(t, expirationTime.After(time.Now()))
}

func TestGenerateRefreshToken(t *testing.T) {
	token, expirationTime, err := GenerateRefreshToken(secretKey, userID, username)

	assert.NoError(t, err)

	assert.NotEmpty(t, token)

	assert.True(t, expirationTime.After(time.Now()))
}

func TestValidateToken_Valid(t *testing.T) {
	token, _, err := GenerateAccessToken(secretKey, userID, username)
	assert.NoError(t, err)

	claims, err := ValidateToken(secretKey, token)
	assert.NoError(t, err)

	assert.Equal(t, userID, claims.UserID)
	assert.Equal(t, username, claims.Username)
}

func TestValidateToken_Invalid(t *testing.T) {
	token, _, err := GenerateAccessToken(secretKey, userID, username)
	assert.NoError(t, err)

	invalidToken := token + "invalid"

	claims, err := ValidateToken(secretKey, invalidToken)

	assert.Error(t, err)
	assert.Nil(t, claims)
}

func TestGenerateAndValidateRefreshToken(t *testing.T) {
	token, _, err := GenerateRefreshToken(secretKey, userID, username)
	assert.NoError(t, err)

	claims, err := ValidateToken(secretKey, token)
	assert.NoError(t, err)

	assert.Equal(t, userID, claims.UserID)
	assert.Equal(t, username, claims.Username)
}
