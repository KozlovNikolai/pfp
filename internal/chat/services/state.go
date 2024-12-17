package services

import (
	"context"
	"log"

	"github.com/KozlovNikolai/pfp/internal/chat/domain"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

// UserService is a User service
type StateService struct {
	repoState IStateRepository
	repoUser  IUserRepository
}

// DeleteConnFromState implements ws.IStateService.
func (s StateService) DeleteConnFromState(ctx context.Context, userID int, pubsub uuid.UUID) (domain.State, bool) {
	return s.repoState.DeleteConnFromState(ctx, userID, pubsub)
}

// GetState implements ws.IStateService.
func (s StateService) GetState(ctx context.Context, userID int) (domain.State, bool) {
	return s.repoState.GetState(ctx, userID)
}

// user,state, index of connect, ifExists
func (s StateService) GetStateByPubsub(ctx context.Context, pubsub uuid.UUID) (domain.User, domain.State, int, bool) {
	userID, state, indexConn, ok := s.repoState.GetStateByPubsub(ctx, pubsub)
	user, err := s.repoUser.GetUserByID(ctx, userID)
	if err != nil {
		log.Printf("service GetStateByPubsub GetUserByID - failure")
		return domain.User{}, domain.State{}, 0, false
	}
	return user, state, indexConn, ok
}

// SetState implements ws.IStateService.
func (s StateService) SetState(ctx context.Context, userID int, pubsub uuid.UUID, conn *websocket.Conn) domain.State {
	return s.repoState.SetState(ctx, userID, pubsub, conn)
}

func (s StateService) SetConnIntoState(ctx context.Context, userID int, pubsub uuid.UUID, conn *websocket.Conn, indexConn int) bool {
	return s.repoState.SetConnIntoState(ctx, userID, pubsub, conn, indexConn)
}

// NewUserService creates a new User service
func NewStateService(repoState IStateRepository, repoUser IUserRepository) StateService {
	return StateService{
		repoState: repoState,
		repoUser:  repoUser,
	}
}

func (s StateService) GetAllStates(ctx context.Context) []domain.State {
	return s.repoState.GetAllStates(ctx)
}

func (s StateService) SetCurrentChat(ctx context.Context, userID int, pubsub uuid.UUID, chatID int) bool {
	return s.repoState.SetCurrentChat(ctx, userID, pubsub, chatID)
}

// // GetUserByID ...
// func (s UserService) GetUserByID(ctx context.Context, id int) (domain.User, error) {
// 	return s.repo.GetUserByID(ctx, id)
// }

// // GetUserByLogin ...
// func (s UserService) GetUserByLogin(ctx context.Context, login string) (domain.User, error) {
// 	return s.repo.GetUserByLogin(ctx, login)
// }

// CreateUser ...
// func (s StateService) CreateUser(ctx context.Context, user domain.User) (domain.User, error) {
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
