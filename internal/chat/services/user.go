package services

// import (
// 	"context"
// 	"fmt"
// 	"time"

// 	"github.com/KozlovNikolai/pfp/internal/chat/domain"
// 	"github.com/KozlovNikolai/pfp/internal/pkg/utils"
// )

// // UserService is a User service
// type UserService struct {
// 	repo IUserRepository
// }

// // NewUserService creates a new User service
// func NewUserService(repo IUserRepository) UserService {
// 	return UserService{
// 		repo: repo,
// 	}
// }

// // GetUserByID ...
// func (s UserService) GetUserByID(ctx context.Context, id int) (domain.UserChat, error) {
// 	return s.repo.GetUserByID(ctx, id)
// }

// // GetUserByLogin ...
// func (s UserService) GetUserByLogin(ctx context.Context, account, login string) (domain.UserChat, error) {
// 	return s.repo.GetUserByLogin(ctx, account, login)
// }

// // CreateUser ...
// func (s UserService) CreateUser(ctx context.Context, user domain.UserChat) (domain.UserChat, error) {
// 	creatingTime := time.Now().Unix()

// 	password, err := utils.HashPassword(user.Password())
// 	if err != nil {
// 		return domain.UserChat{}, fmt.Errorf("error-hashing-password: %v", err.Error())
// 	}

// 	newUser := domain.NewUserChatData{
// 		UserExtID: user.UserExtID(),
// 		Login:     user.Login(),
// 		Password:  password,
// 		Account:   user.Account(),
// 		Token:     user.Token(),
// 		Name:      user.Name(),
// 		Surname:   user.Surname(),
// 		Email:     user.Email(),
// 		UserType:  user.UserType(),
// 		CreatedAt: creatingTime,
// 		UpdatedAt: creatingTime,
// 	}
// 	creatingUser := domain.NewUserChat(newUser)
// 	return s.repo.CreateUser(ctx, creatingUser)
// }

// // UpdateUser ...
// func (s UserService) UpdateUser(ctx context.Context, user domain.User) (domain.User, error) {
// 	return s.repo.UpdateUser(ctx, user)
// }

// // DeleteUser ...
// func (s UserService) DeleteUser(ctx context.Context, id int) error {
// 	return s.repo.DeleteUser(ctx, id)
// }

// // GetUsers ...
// func (s UserService) GetUsers(ctx context.Context, limit, offset int) ([]domain.User, error) {
// 	return s.repo.GetUsers(ctx, limit, offset)
// }
