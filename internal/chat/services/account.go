package services

import (
	"context"

	"alfachat/internal/chat/constants"
	"alfachat/internal/chat/domain"
)

// UserService is a User service
type AccountService struct {
	repo IAccountRepository
}

// NewUserService creates a new User service
func NewAccountService(repo IAccountRepository) AccountService {
	return AccountService{
		repo: repo,
	}
}

// // CreateUser ...
func (a AccountService) CreateAccount(ctx context.Context, account domain.Account) (domain.Account, error) {
	// creatingTime := time.Now().Unix()

	// userDb, err := s.GetUserByLogin(ctx, user.Profile(), user.Login())
	// if err != nil {
	// 	if strings.HasSuffix(err.Error(), "no rows in result set") {
	// 		password, err := utils.HashPassword(user.Password())
	// 		if err != nil {
	// 			return domain.User{}, fmt.Errorf("error-hashing-password: %v", err.Error())
	// 		}

	// 		newUser := domain.NewUserData{
	// 			UserExtID: user.UserExtID(),
	// 			Login:     user.Login(),
	// 			Password:  password,
	// 			Profile:   user.Profile(),
	// 			Token:     user.Token(),
	// 			Name:      user.Name(),
	// 			Surname:   user.Surname(),
	// 			Email:     user.Email(),
	// 			UserType:  user.UserType(),
	// 			CreatedAt: creatingTime,
	// 			UpdatedAt: creatingTime,
	// 		}
	// 		creatingUser := domain.NewUser(newUser)
	// 		return s.repo.CreateUser(ctx, creatingUser)
	// 	}
	// }
	// return domain.User{}, fmt.Errorf("user with profile: %s and login: %s already exists", userDb.Profile(), userDb.Login())
	return a.repo.CreateAccount(ctx, account)
}

func (a AccountService) AddUserToAccount(ctx context.Context, userID int, accountID int, inviterID int, role string) error {
	return a.repo.AddUserToAccount(ctx, userID, accountID, inviterID, role)
}

func (a AccountService) NewUserToNewAccount(ctx context.Context, userID int, accountID int) error {
	return a.repo.AddUserToAccount(ctx, userID, accountID, userID, constants.AccountRoleOwner)
}

func (a AccountService) GetAccountByUserID(ctx context.Context, userID int) (int, error) {
	return a.repo.GetAccountByUserID(ctx, userID)
}

func (a AccountService) GetContactsByAccount(ctx context.Context, accID int) ([]int, error) {
	return a.repo.GetContactsByAccount(ctx, accID)
}

// // // GetUserByID ...
// func (s UserService) GetUserByID(ctx context.Context, id int) (domain.User, error) {
// 	return s.repo.GetUserByID(ctx, id)
// }

// // GetUserByID ...
// func (s UserService) GetUserByExtID(ctx context.Context, account, extId string) (domain.User, error) {
// 	return s.repo.GetUserByExtID(ctx, account, extId)
// }

// // GetUserByLogin ...
// func (s UserService) GetUserByLogin(ctx context.Context, profile, login string) (domain.User, error) {
// 	return s.repo.GetUserByLogin(ctx, profile, login)
// }

// // RegisterUser ...
// func (s UserService) RegisterUser(ctx context.Context, user domain.User) (domain.User, bool, error) {
// 	userDb, err := s.GetUserByExtID(ctx, user.Profile(), user.UserExtID())

// 	if err != nil {
// 		if strings.HasSuffix(err.Error(), "no rows in result set") {
// 			userDomain, err := s.repo.CreateUser(ctx, user)

// 			return userDomain, true, err
// 		}
// 		return domain.User{}, false, fmt.Errorf("register failure of user with ext id: %s, error: %s", user.UserExtID(), err.Error())
// 	}
// 	return userDb, false, nil
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
// func (s UserService) GetUsers(ctx context.Context, user domain.User, limit, offset int) ([]domain.User, error) {
// 	account := user.Profile()
// 	// fmt.Println(account, limit, offset)
// 	return s.repo.GetUsers(ctx, account, limit, offset)
// }
