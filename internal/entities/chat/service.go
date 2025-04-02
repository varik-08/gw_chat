package chat

import "fmt"

type Service struct {
	chatRepository RepositoryInterface
}

func NewChatService(chatRepository RepositoryInterface) *Service {
	return &Service{chatRepository: chatRepository}
}

func (s *Service) GetChatsByUserID(userID int) ([]Chat, error) {
	return s.chatRepository.GetChatsByUserID(userID)
}

func (s *Service) CreateChat(chatDTO *DTO) (int, error) {
	chat := &Chat{
		Name:     chatDTO.Name,
		IsPublic: chatDTO.IsPublic,
		OwnerID:  chatDTO.OwnerID,
	}

	chatID, err := s.chatRepository.CreateChat(chat)
	if err != nil {
		return 0, fmt.Errorf("failed to create chat: %w", err)
	}

	err = s.chatRepository.AddMember(&Member{ChatID: chatID, UserID: chat.OwnerID})
	if err != nil {
		return 0, fmt.Errorf("failed to add owner to chat: %w", err)
	}

	return chatID, nil
}

func (s *Service) AddMember(chatMemberDTO *MemberDTO, currentUser int) error {
	chat, err := s.chatRepository.GetChatByID(chatMemberDTO.ChatID)
	if err != nil {
		return fmt.Errorf("failed to get chat: %w", err)
	}

	if chat.OwnerID != currentUser {
		return fmt.Errorf("only chat owner can add members")
	}

	exists, err := s.chatRepository.CheckExistsUserInChat(chatMemberDTO.UserID, chat.ID)
	if err != nil {
		return fmt.Errorf("failed to check user in chat: %w", err)
	}

	if exists {
		return fmt.Errorf("user already exists in chat")
	}

	chatMember := &Member{ChatID: chatMemberDTO.ChatID, UserID: chatMemberDTO.UserID}

	return s.chatRepository.AddMember(chatMember)
}
