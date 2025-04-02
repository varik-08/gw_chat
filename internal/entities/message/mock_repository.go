package message

import "github.com/stretchr/testify/mock"

type MockMessageRepository struct {
	mock.Mock
}

func (m *MockMessageRepository) GetMessagesByChatID(chatID int) ([]*Message, error) {
	args := m.Called(chatID)
	return args.Get(0).([]*Message), args.Error(1)
}

func (m *MockMessageRepository) CreateChatMessage(message *Message) (int, error) {
	args := m.Called(message)
	return args.Int(0), args.Error(1)
}
