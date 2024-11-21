package services

import (
	"context"
	"fmt"
	"time"

	"github.com/KozlovNikolai/pfp/internal/chat/domain"
	"github.com/KozlovNikolai/pfp/internal/pkg/utils"
)

// UserService is a User service
type UserChatService struct {
	repo IUserRepository
}

// NewUserService creates a new User service
func NewUserChatService(repo IUserRepository) UserChatService {
	return UserChatService{
		repo: repo,
	}
}

// GetUserByID ...
func (s UserChatService) GetUserByID(ctx context.Context, id int) (domain.User, error) {
	return s.repo.GetUserByID(ctx, id)
}

// GetUserByLogin ...
func (s UserChatService) GetUserByLogin(ctx context.Context, login string) (domain.User, error) {
	return s.repo.GetUserByLogin(ctx, login)
}

// CreateUser ...
func (s UserChatService) RegisterUser(ctx context.Context, user domain.UserChat) (domain.UserChat, error) {
	creatingTime := time.Now()

	password, err := utils.HashPassword(user.Password())
	if err != nil {
		return domain.User{}, fmt.Errorf("error-hashing-password: %v", err.Error())
	}

	newUser := domain.NewUserData{
		Login:     user.Login(),
		Password:  password,
		Role:      "regular",
		Token:     "",
		CreatedAt: creatingTime,
		UpdatedAt: creatingTime,
	}
	creatingUser := domain.NewUser(newUser)
	return s.repo.CreateUser(ctx, creatingUser)
}

// UpdateUser ...
func (s UserChatService) UpdateUser(ctx context.Context, user domain.User) (domain.User, error) {
	return s.repo.UpdateUser(ctx, user)
}

// DeleteUser ...
func (s UserChatService) DeleteUser(ctx context.Context, id int) error {
	return s.repo.DeleteUser(ctx, id)
}

// GetUsers ...
func (s UserChatService) GetUsers(ctx context.Context, limit, offset int) ([]domain.User, error) {
	return s.repo.GetUsers(ctx, limit, offset)
}
