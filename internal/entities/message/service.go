package message

import (
	"fmt"
	"github.com/varik-08/gw_chat/internal/entities/chat"
)

type MessageService struct {
	messageRepository *MessageRepository
	chatRepository    *chat.ChatRepository
}

func NewChatService(messageRepository *MessageRepository, chatRepository *chat.ChatRepository) *MessageService {
	return &MessageService{
		messageRepository: messageRepository,
		chatRepository:    chatRepository,
	}
}

func (s *MessageService) GetMessagesByChatID(chatID int, currentUser int) ([]*Message, error) {
	exists, err := s.chatRepository.CheckExistsUserInChat(currentUser, chatID)
	if err != nil {
		return nil, fmt.Errorf("failed to check user in chat: %w", err)
	}

	if !exists {
		return nil, fmt.Errorf("user is not in chat: %w", err)
	}

	messages, err := s.messageRepository.GetMessagesByChatID(chatID)
	if err != nil {
		return nil, fmt.Errorf("failed to get messages: %w", err)
	}

	return messages, nil
}

func (s *MessageService) CreateChatMessage(messageDTO *MessageDTO) (int, error) {
	message := &Message{
		ChatID:  messageDTO.ChatID,
		UserID:  messageDTO.UserID,
		Content: messageDTO.Content,
	}

	messageID, err := s.messageRepository.CreateChatMessage(message)
	if err != nil {
		return 0, fmt.Errorf("failed to create message: %w", err)
	}

	return messageID, nil
}
