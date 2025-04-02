package users

import "github.com/stretchr/testify/mock"

type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) GetUserByUsername(username string) (*User, error) {
	args := m.Called(username)
	if args.Get(0) != nil {
		return args.Get(0).(*User), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockUserRepository) CreateUser(user *User) (int, error) {
	args := m.Called(user)
	return args.Int(0), args.Error(1)
}

func (m *MockUserRepository) GetUserByID(id int) (*User, error) {
	args := m.Called(id)
	if args.Get(0) != nil {
		return args.Get(0).(*User), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockUserRepository) UpdatePassword(user *User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserRepository) GetUsers() ([]*User, error) {
	args := m.Called()
	if args.Get(0) != nil {
		return args.Get(0).([]*User), args.Error(1)
	}
	return nil, args.Error(1)
}
