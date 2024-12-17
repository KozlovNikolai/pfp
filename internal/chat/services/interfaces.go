// Package services ...
package services

import (
	"context"

	"github.com/KozlovNikolai/pfp/internal/chat/domain"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

// IUserRepository is ...
type IUserRepository interface {
	CreateUser(context.Context, domain.User) (domain.User, error)
	GetUsers(context.Context, string, int, int) ([]domain.User, error)
	GetUserByID(context.Context, int) (domain.User, error)
	GetUserByExtID(context.Context, string, string) (domain.User, error)
	GetUserByLogin(context.Context, string, string) (domain.User, error)
	UpdateUser(context.Context, domain.User) (domain.User, error)
	DeleteUser(context.Context, int) error
}

type IStateRepository interface {
	SetState(ctx context.Context, userID int, pubsub uuid.UUID, conn *websocket.Conn) domain.State
	GetState(ctx context.Context, userID int) (domain.State, bool)
	DeleteConnFromState(ctx context.Context, userID int, pubsub uuid.UUID) (domain.State, bool)
	GetStateByPubsub(ctx context.Context, pubsub uuid.UUID) (int, domain.State, int, bool) //userID,state, index of connect, ifExists
	SetConnIntoState(ctx context.Context, userID int, pubsub uuid.UUID, conn *websocket.Conn, indexConn int) bool
	GetAllStates(ctx context.Context) []domain.State
	SetCurrentChat(ctx context.Context, userID int, pubsub uuid.UUID, chatID int) bool
}

type IChatRepository interface {
	CreateChat(context.Context, domain.Chat) (domain.Chat, error)
	AddUserToChat(ctx context.Context, userID int, chatID int) error
	GetChatByNameAndType(context.Context, string, string) (domain.Chat, error)
	GetChatsByUser(ctx context.Context, userID int) ([]domain.Chat, error)
	GetUserIDsByChatID(ctx context.Context, chatID int) ([]int, error)
	IsChatMember(ctx context.Context, userID int, chatID int) bool
	GetUsersByChatID(ctx context.Context, chatID int) ([]domain.User, error)
}

type IMessageRepository interface {
	SaveMsg(ctx context.Context, msg domain.Message) error
	GetMessagesByChatID(ctx context.Context, chatID, limit, offset int) ([]domain.Message, error)
}
