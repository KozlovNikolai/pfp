package services

import (
	"context"
	"fmt"
	"strings"
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

// // CreateUser ...
func (s UserChatService) CreateUser(ctx context.Context, user domain.UserChat) (domain.UserChat, error) {
	creatingTime := time.Now().Unix()

	password, err := utils.HashPassword(user.Password())
	if err != nil {
		return domain.UserChat{}, fmt.Errorf("error-hashing-password: %v", err.Error())
	}

	newUser := domain.NewUserChatData{
		UserExtID: user.UserExtID(),
		Login:     user.Login(),
		Password:  password,
		Account:   user.Account(),
		Token:     user.Token(),
		Name:      user.Name(),
		Surname:   user.Surname(),
		Email:     user.Email(),
		UserType:  user.UserType(),
		CreatedAt: creatingTime,
		UpdatedAt: creatingTime,
	}
	creatingUser := domain.NewUserChat(newUser)
	return s.repo.CreateUserChat(ctx, creatingUser)
}

// // GetUserByID ...
func (s UserChatService) GetUserByID(ctx context.Context, id int) (domain.UserChat, error) {
	return s.repo.GetUserByID(ctx, id)
}

// GetUserByID ...
func (s UserChatService) GetUserByExtID(ctx context.Context, account, extId string) (domain.UserChat, error) {
	return s.repo.GetUserByExtID(ctx, account, extId)
}

// GetUserByLogin ...
func (s UserChatService) GetUserByLogin(ctx context.Context, account, login string) (domain.UserChat, error) {
	return s.repo.GetUserByLogin(ctx, account, login)
}

// RegisterUser ...
func (s UserChatService) RegisterUser(ctx context.Context, user domain.UserChat) (domain.UserChat, error) {
	userChat, err := s.GetUserByExtID(ctx, user.Account(), user.UserExtID())

	if err != nil {
		if strings.HasSuffix(err.Error(), "no rows in result set") {
			return s.repo.CreateUserChat(ctx, user)
		}
		return domain.UserChat{}, fmt.Errorf("register failure of user with ext id: %s, error: %s", user.UserExtID(), err.Error())
	}
	return userChat, nil
}

// UpdateUser ...
func (s UserChatService) UpdateUser(ctx context.Context, user domain.UserChat) (domain.UserChat, error) {
	return s.repo.UpdateUser(ctx, user)
}

// DeleteUser ...
func (s UserChatService) DeleteUser(ctx context.Context, id int) error {
	return s.repo.DeleteUser(ctx, id)
}

// GetUsers ...
func (s UserChatService) GetUsers(ctx context.Context, limit, offset int) ([]domain.UserChat, error) {
	return s.repo.GetUsers(ctx, limit, offset)
}
