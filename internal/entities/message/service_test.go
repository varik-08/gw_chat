package message

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/varik-08/gw_chat/internal/entities/chat"
)

const ChatID = 1

const MessageID = 1

const CurrentUserID = 1

func TestService_GetMessagesByChatIDSuccess(t *testing.T) {
	mockMessageRepo := new(MockMessageRepository)
	mockChatRepo := new(chat.MockChatRepository)

	mockChatRepo.On("CheckExistsUserInChat", CurrentUserID, ChatID).Return(true, nil)
	mockMessageRepo.On("GetMessagesByChatID", ChatID).Return([]*Message{}, nil)

	service := NewMessageService(mockMessageRepo, mockChatRepo)

	messages, err := service.GetMessagesByChatID(ChatID, CurrentUserID)

	assert.NoError(t, err)
	assert.Equal(t, []*Message{}, messages)

	mockMessageRepo.AssertExpectations(t)
	mockChatRepo.AssertExpectations(t)
}

func TestService_GetMessagesByChatIDFailure(t *testing.T) {
	mockMessageRepo := new(MockMessageRepository)
	mockChatRepo := new(chat.MockChatRepository)

	mockChatRepo.On("CheckExistsUserInChat", CurrentUserID, ChatID).Return(false, nil)

	service := NewMessageService(mockMessageRepo, mockChatRepo)

	_, err := service.GetMessagesByChatID(ChatID, CurrentUserID)

	assert.Error(t, err)

	mockMessageRepo.AssertExpectations(t)
	mockChatRepo.AssertExpectations(t)
}

func TestService_SendMessageSuccess(t *testing.T) {
	mockMessageRepo := new(MockMessageRepository)
	mockChatRepo := new(chat.MockChatRepository)

	mockMessageRepo.On("CreateChatMessage", &Message{}).Return(MessageID, nil)

	service := NewMessageService(mockMessageRepo, mockChatRepo)

	newMessageID, err := service.CreateChatMessage(&DTO{})

	assert.NoError(t, err)
	assert.Equal(t, newMessageID, MessageID)

	mockMessageRepo.AssertExpectations(t)
	mockChatRepo.AssertExpectations(t)
}

func TestService_SendMessageFailure(t *testing.T) {
	mockMessageRepo := new(MockMessageRepository)
	mockChatRepo := new(chat.MockChatRepository)

	mockMessageRepo.On("CreateChatMessage", &Message{}).Return(0, errors.New("error"))

	service := NewMessageService(mockMessageRepo, mockChatRepo)

	_, err := service.CreateChatMessage(&DTO{})

	assert.Error(t, err)

	mockMessageRepo.AssertExpectations(t)
	mockChatRepo.AssertExpectations(t)
}
