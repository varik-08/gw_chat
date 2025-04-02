package chat

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

const ChatID = 1

const UserID = 2

const CurrentUserID = 1

func TestService_CreateChatSuccess(t *testing.T) {
	mockRepo := new(MockChatRepository)

	mockRepo.On("CreateChat", &Chat{}).Return(ChatID, nil)
	mockRepo.On("AddMember", &Member{ChatID: ChatID}).Return(nil)

	service := NewChatService(mockRepo)

	newChatID, err := service.CreateChat(&DTO{})

	assert.NoError(t, err)
	assert.NotNil(t, newChatID)
	assert.Equal(t, newChatID, ChatID)

	mockRepo.AssertExpectations(t)
}

func TestService_CreateChatFailure(t *testing.T) {
	mockRepo := new(MockChatRepository)

	mockRepo.On("CreateChat", &Chat{}).Return(0, errors.New("error"))

	service := NewChatService(mockRepo)

	_, err := service.CreateChat(&DTO{})

	assert.Error(t, err)

	mockRepo.AssertExpectations(t)
}

func TestService_AddMemberSuccess(t *testing.T) {
	mockRepo := new(MockChatRepository)

	mockRepo.On("GetChatByID", ChatID).Return(&Chat{ID: ChatID, OwnerID: 1}, nil)
	mockRepo.On("CheckExistsUserInChat", UserID, ChatID).Return(false, nil)
	mockRepo.On("AddMember", &Member{ChatID: ChatID, UserID: UserID}).Return(nil)

	service := NewChatService(mockRepo)

	err := service.AddMember(&MemberDTO{ChatID: ChatID, UserID: UserID}, CurrentUserID)

	assert.NoError(t, err)

	mockRepo.AssertExpectations(t)
}

func TestService_AddMemberExists(t *testing.T) {
	mockRepo := new(MockChatRepository)

	mockRepo.On("GetChatByID", ChatID).Return(&Chat{ID: ChatID, OwnerID: 1}, nil)
	mockRepo.On("CheckExistsUserInChat", UserID, ChatID).Return(true, nil)

	service := NewChatService(mockRepo)

	err := service.AddMember(&MemberDTO{ChatID: ChatID, UserID: UserID}, CurrentUserID)

	assert.Error(t, err)

	mockRepo.AssertExpectations(t)
}
