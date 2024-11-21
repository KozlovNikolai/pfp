package services

// UserService is a User service
type StateService struct {
	repo IStateRepository
}

// NewUserService creates a new User service
func NewStateService(repo IStateRepository) StateService {
	return StateService{
		repo: repo,
	}
}

// // GetUserByID ...
// func (s UserService) GetUserByID(ctx context.Context, id int) (domain.User, error) {
// 	return s.repo.GetUserByID(ctx, id)
// }

// // GetUserByLogin ...
// func (s UserService) GetUserByLogin(ctx context.Context, login string) (domain.User, error) {
// 	return s.repo.GetUserByLogin(ctx, login)
// }

// // CreateUser ...
// func (s UserService) CreateUser(ctx context.Context, user domain.User) (domain.User, error) {
// 	creatingTime := time.Now()

// 	password, err := utils.HashPassword(user.Password())
// 	if err != nil {
// 		return domain.User{}, fmt.Errorf("error-hashing-password: %v", err.Error())
// 	}

// 	newUser := domain.NewUserData{
// 		Login:     user.Login(),
// 		Password:  password,
// 		Role:      "regular",
// 		Token:     "",
// 		CreatedAt: creatingTime,
// 		UpdatedAt: creatingTime,
// 	}
// 	creatingUser := domain.NewUser(newUser)
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
