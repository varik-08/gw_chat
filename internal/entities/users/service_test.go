package users

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

const UserID = 1

const UserPassword = "password"

func TestService_CreateUserSuccess(t *testing.T) {
	mockUserRepo := new(MockUserRepository)

	mockUserRepo.On("CreateUser", mock.Anything).Return(UserID, nil)

	service := NewUserService(mockUserRepo)

	newUserID, err := service.CreateUser(Credentials{Password: UserPassword})

	assert.NoError(t, err)
	assert.Equal(t, newUserID, UserID)

	mockUserRepo.AssertExpectations(t)
}

func TestService_CreateUserFailure(t *testing.T) {
	mockUserRepo := new(MockUserRepository)

	mockUserRepo.On("CreateUser", mock.Anything).Return(0, errors.New("error"))

	service := NewUserService(mockUserRepo)

	_, err := service.CreateUser(Credentials{Password: UserPassword})

	assert.Error(t, err)

	mockUserRepo.AssertExpectations(t)
}
