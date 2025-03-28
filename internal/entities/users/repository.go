package users

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) CreateUser(user *User) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var userID int

	err := r.db.QueryRow(ctx, "INSERT INTO users (username, password_hash) VALUES ($1, $2) RETURNING id",
		user.Username, user.PasswordHash).Scan(&userID)

	return userID, err
}

func (r *UserRepository) GetUserById(id int) (*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var user User
	err := r.db.QueryRow(ctx, "SELECT id, username, password_hash FROM users WHERE id=$1", id).
		Scan(&user.ID, &user.Username, &user.PasswordHash)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) GetUserByUsername(username string) (*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var user User
	err := r.db.QueryRow(ctx, "SELECT id, username, password_hash FROM users WHERE username=$1", username).
		Scan(&user.ID, &user.Username, &user.PasswordHash)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) UpdatePassword(user *User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := r.db.Exec(ctx, "UPDATE users SET password_hash=$1 WHERE id=$2", user.PasswordHash, user.ID)

	return err
}
