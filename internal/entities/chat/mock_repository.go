package chat

import "github.com/stretchr/testify/mock"

type MockChatRepository struct {
	mock.Mock
}

func (m *MockChatRepository) CreateChat(chat *Chat) (int, error) {
	args := m.Called(chat)
	return args.Int(0), args.Error(1)
}

func (m *MockChatRepository) AddMember(member *Member) error {
	args := m.Called(member)
	return args.Error(0)
}

func (m *MockChatRepository) GetChatByID(chatID int) (*Chat, error) {
	args := m.Called(chatID)
	return args.Get(0).(*Chat), args.Error(1)
}

func (m *MockChatRepository) CheckExistsUserInChat(userID int, chatID int) (bool, error) {
	args := m.Called(userID, chatID)
	return args.Bool(0), args.Error(1)
}

func (m *MockChatRepository) GetChatsByUserID(userID int) ([]Chat, error) {
	args := m.Called(userID)
	if chats, ok := args.Get(0).([]Chat); ok {
		return chats, args.Error(1)
	}
	return nil, args.Error(1)
}
