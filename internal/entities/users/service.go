package users

import (
	"fmt"

	"github.com/varik-08/gw_chat/internal/pkg"
)

type UserService struct {
	UserRepository UserRepositoryInterface
}

func NewUserService(userRepository UserRepositoryInterface) *UserService {
	return &UserService{UserRepository: userRepository}
}

func (s *UserService) CreateUser(credentials Credentials) (int, error) {
	hashedPassword, err := pkg.HashPassword(credentials.Password)
	if err != nil {
		return 0, fmt.Errorf("failed to hash password: %w", err)
	}

	user := User{
		Username:     credentials.Username,
		PasswordHash: hashedPassword,
	}

	userID, err := s.UserRepository.CreateUser(&user)

	return userID, err
}

func (s *UserService) UpdatePassword(userID int, newPassword string, oldPassword string) error {
	user, err := s.UserRepository.GetUserByID(userID)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}

	if !pkg.CheckPasswordHash(oldPassword, user.PasswordHash) {
		return fmt.Errorf("invalid old password")
	}

	newHashedPassword, err := pkg.HashPassword(newPassword)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	user.PasswordHash = newHashedPassword

	err = s.UserRepository.UpdatePassword(user)
	if err != nil {
		return fmt.Errorf("failed to update password: %w", err)
	}

	return nil
}

func (s *UserService) GetUsers() ([]*User, error) {
	return s.UserRepository.GetUsers()
}
