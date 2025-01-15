package services

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/KozlovNikolai/pfp/internal/chat/constants"
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

	userDb, err := s.GetUserByLogin(ctx, user.Profile(), user.Login())
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
				Profile:   user.Profile(),
				Name:      user.Name(),
				Surname:   user.Surname(),
				Email:     user.Email(),
				UserType:  user.UserType(),
				CreatedAt: creatingTime,
				UpdatedAt: creatingTime,
			}
			newUser.Login = newUser.Email

			//	fmt.Printf("Login: %s, Profile: %s, Role: %s\n", newUser.Login, newUser.Profile, newUser.UserType)
			if newUser.Login == constants.AdminEmaleLogin && newUser.Profile == constants.SystemProfile {
				newUser.UserType = constants.UserTypeAdmin
			}
			//	fmt.Printf("Login: %s, Profile: %s, Role: %s\n", newUser.Login, newUser.Profile, newUser.UserType)
			creatingUser := domain.NewUser(newUser)
			return s.repo.CreateUser(ctx, creatingUser)
		}
	}
	return domain.User{}, fmt.Errorf("user with profile: %s and login: %s already exists", userDb.Profile(), userDb.Login())
}

// // GetUserByID ...
func (s UserService) GetUserByID(ctx context.Context, id int) (domain.User, error) {
	return s.repo.GetUserByID(ctx, id)
}

// GetUserByID ...
func (s UserService) GetUserByExtID(ctx context.Context, account string, extId int) (domain.User, error) {
	return s.repo.GetUserByExtID(ctx, account, extId)
}

// GetUserByLogin ...
func (s UserService) GetUserByLogin(ctx context.Context, profile, login string) (domain.User, error) {
	return s.repo.GetUserByLogin(ctx, profile, login)
}

// RegisterUser ...
func (s UserService) RegisterUser(ctx context.Context, user domain.User) (domain.User, bool, error) {
	userDb, err := s.GetUserByExtID(ctx, user.Profile(), user.UserExtID())

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
	account := user.Profile()
	// fmt.Println(account, limit, offset)
	return s.repo.GetUsers(ctx, account, limit, offset)
}

func (u UserService) AddContact(ctx context.Context, user domain.User, userID int) error {
	return u.repo.AddContact(ctx, user, userID)
}

func (u UserService) FindUsers(ctx context.Context, search string, start int, stop int) ([]domain.User, error) {
	return u.repo.FindUsers(ctx, search, start, stop)
}
