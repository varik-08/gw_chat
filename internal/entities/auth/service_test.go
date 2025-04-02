package auth

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/varik-08/gw_chat/internal/entities/users"
	"github.com/varik-08/gw_chat/internal/pkg"
	"github.com/varik-08/gw_chat/internal/pkg/jwt"
)

const SecretKey = "secret"

const UserID = 1

const Username = "testuser"

const CorrectPassword = "correctpassword"

const WrongPassword = "wrongpassword"

const WrongRefreshToken = "wrongrefreshtoken"

func TestService_LoginSuccess(t *testing.T) {
	mockRepo := new(users.MockUserRepository)

	hash, err := pkg.HashPassword(CorrectPassword)
	if err != nil {
		t.Fatal(err)
	}

	testUser := &users.User{
		ID:           UserID,
		Username:     Username,
		PasswordHash: hash,
	}

	mockRepo.On("GetUserByUsername", Username).Return(testUser, nil)

	service := NewAuthService(SecretKey, mockRepo)

	credentials := users.Credentials{
		Username: Username,
		Password: CorrectPassword,
	}

	token, err := service.Login(credentials)

	assert.NoError(t, err)
	assert.NotNil(t, token)
	assert.Equal(t, testUser.ID, token.UserID)
	assert.Equal(t, testUser.Username, token.Username)

	mockRepo.AssertExpectations(t)
}

func TestService_LoginFailure_UserNotFound(t *testing.T) {
	mockRepo := new(users.MockUserRepository)

	mockRepo.On("GetUserByUsername", Username).Return(nil, nil)

	service := NewAuthService(SecretKey, mockRepo)

	credentials := users.Credentials{
		Username: Username,
		Password: CorrectPassword,
	}

	_, err := service.Login(credentials)

	assert.Error(t, err)

	mockRepo.AssertExpectations(t)
}

func TestService_LoginFailure_WrongPassword(t *testing.T) {
	mockRepo := new(users.MockUserRepository)

	hash, err := pkg.HashPassword(CorrectPassword)
	if err != nil {
		t.Fatal(err)
	}

	testUser := &users.User{
		ID:           UserID,
		Username:     Username,
		PasswordHash: hash,
	}

	mockRepo.On("GetUserByUsername", Username).Return(testUser, nil)

	service := NewAuthService(SecretKey, mockRepo)

	credentials := users.Credentials{
		Username: Username,
		Password: WrongPassword,
	}

	_, err = service.Login(credentials)

	assert.Error(t, err)

	mockRepo.AssertExpectations(t)
}

func TestService_RefreshTokenSuccess(t *testing.T) {
	mockRepo := new(users.MockUserRepository)

	service := NewAuthService(SecretKey, mockRepo)

	newRefreshToken, _, err := jwt.GenerateRefreshToken(SecretKey, UserID, Username)
	if err != nil {
		t.Fatal(err)
	}

	oldToken := &Token{
		RefreshToken: newRefreshToken,
	}

	token, err := service.RefreshToken(oldToken)

	assert.NoError(t, err)
	assert.NotNil(t, token)
}

func TestService_RefreshTokenFailure_InvalidToken(t *testing.T) {
	mockRepo := new(users.MockUserRepository)

	service := NewAuthService(SecretKey, mockRepo)

	oldToken := &Token{
		RefreshToken: WrongRefreshToken,
	}

	token, err := service.RefreshToken(oldToken)

	assert.Error(t, err)
	assert.Nil(t, token)
}
