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
type UserService struct {
	repo IUserRepository
}

// NewUserService creates a new User service
func NewUserService(repo IUserRepository) UserService {
	return UserService{
		repo: repo,
	}
}

// // CreateUser ...
func (s UserService) CreateUser(ctx context.Context, user domain.User) (domain.User, error) {
	creatingTime := time.Now().Unix()

	userDb, err := s.GetUserByLogin(ctx, user.Account(), user.Login())
	if err != nil {
		if strings.HasSuffix(err.Error(), "no rows in result set") {
			password, err := utils.HashPassword(user.Password())
			if err != nil {
				return domain.User{}, fmt.Errorf("error-hashing-password: %v", err.Error())
			}

			newUser := domain.NewUserData{
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
			creatingUser := domain.NewUser(newUser)
			return s.repo.CreateUser(ctx, creatingUser)
		}
	}
	return domain.User{}, fmt.Errorf("user with account: %s and login: %s already exists", userDb.Account(), userDb.Login())
}

// // GetUserByID ...
func (s UserService) GetUserByID(ctx context.Context, id int) (domain.User, error) {
	return s.repo.GetUserByID(ctx, id)
}

// GetUserByID ...
func (s UserService) GetUserByExtID(ctx context.Context, account, extId string) (domain.User, error) {
	return s.repo.GetUserByExtID(ctx, account, extId)
}

// GetUserByLogin ...
func (s UserService) GetUserByLogin(ctx context.Context, account, login string) (domain.User, error) {
	return s.repo.GetUserByLogin(ctx, account, login)
}

// RegisterUser ...
func (s UserService) RegisterUser(ctx context.Context, user domain.User) (domain.User, bool, error) {
	userDb, err := s.GetUserByExtID(ctx, user.Account(), user.UserExtID())

	if err != nil {
		if strings.HasSuffix(err.Error(), "no rows in result set") {
			userDomain, err := s.repo.CreateUser(ctx, user)

			return userDomain, true, err
		}
		return domain.User{}, false, fmt.Errorf("register failure of user with ext id: %s, error: %s", user.UserExtID(), err.Error())
	}
	return userDb, false, nil
}

// UpdateUser ...
func (s UserService) UpdateUser(ctx context.Context, user domain.User) (domain.User, error) {
	return s.repo.UpdateUser(ctx, user)
}

// DeleteUser ...
func (s UserService) DeleteUser(ctx context.Context, id int) error {
	return s.repo.DeleteUser(ctx, id)
}

// GetUsers ...
func (s UserService) GetUsers(ctx context.Context, user domain.User, limit, offset int) ([]domain.User, error) {
	account := user.Account()
	// fmt.Println(account, limit, offset)
	return s.repo.GetUsers(ctx, account, limit, offset)
}
