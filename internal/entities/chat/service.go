package chat

import "fmt"

type ChatService struct {
	chatRepository *ChatRepository
}

func NewChatService(chatRepository *ChatRepository) *ChatService {
	return &ChatService{chatRepository: chatRepository}
}

func (s *ChatService) GetChatsByUserID(userID int) ([]Chat, error) {
	return s.chatRepository.GetChatsByUserID(userID)
}

func (s *ChatService) CreateChat(chatDTO *ChatDTO) (int, error) {
	chat := &Chat{
		Name:     chatDTO.Name,
		IsPublic: chatDTO.IsPublic,
		OwnerId:  chatDTO.OwnerId,
	}

	chatID, err := s.chatRepository.CreateChat(chat)
	if err != nil {
		return 0, fmt.Errorf("failed to create chat: %w", err)
	}

	err = s.chatRepository.AddMember(&ChatMember{ChatID: chatID, UserID: chat.OwnerId})
	if err != nil {
		return 0, fmt.Errorf("failed to add owner to chat: %w", err)
	}

	return chatID, nil
}

func (s *ChatService) AddMember(chatMemberDTO *ChatMemberDTO, currentUser int) error {
	chat, err := s.chatRepository.GetChatByID(chatMemberDTO.ChatID)
	if err != nil {
		return fmt.Errorf("failed to get chat: %w", err)
	}

	if chat.OwnerId != currentUser {
		return fmt.Errorf("only chat owner can add members")
	}

	exists, err := s.chatRepository.CheckExistsUserInChat(chatMemberDTO.UserID, chat.ID)
	if err != nil {
		return fmt.Errorf("failed to check user in chat: %w", err)
	}

	if exists {
		return fmt.Errorf("user already exists in chat")
	}

	chatMember := &ChatMember{ChatID: chatMemberDTO.ChatID, UserID: chatMemberDTO.UserID}

	return s.chatRepository.AddMember(chatMember)
}
