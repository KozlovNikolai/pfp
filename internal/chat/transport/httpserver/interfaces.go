package httpserver

import (
	"context"

	"github.com/KozlovNikolai/pfp/internal/chat/domain"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

// ITokenService is a token service
type ITokenService interface {
	GenerateTokenForRegisteredUsers(ctx context.Context, user domain.User) (string, error)
	GenerateToken(ctx context.Context, account, login, password string) (domain.User, string, error)
	GetUser(ctx context.Context, token string) (domain.User, error)
	GetPubsubToken(ctx context.Context, user domain.User) (uuid.UUID, error)
}

type IStateService interface {
	SetState(ctx context.Context, userID int, pubsub uuid.UUID, conn *websocket.Conn) domain.State
	GetState(ctx context.Context, userID int) (domain.State, bool)
	GetStateByPubsub(ctx context.Context, pubsub uuid.UUID) (domain.User, domain.State, int, bool) //user,state, index of connect, ifExists
	DeleteConnFromState(ctx context.Context, userID int, pubsub uuid.UUID) (domain.State, bool)
	GetAllStates(ctx context.Context) []domain.State
	SetCurrentChat(ctx context.Context, userID int, pubsub uuid.UUID, chatID int) bool
}

type IUserService interface {
	CreateUser(context.Context, domain.User) (domain.User, error)
	RegisterUser(context.Context, domain.User) (domain.User, bool, error)
	GetUsers(context.Context, domain.User, int, int) ([]domain.User, error)
	GetUserByID(context.Context, int) (domain.User, error)
	GetUserByExtID(context.Context, string, int) (domain.User, error)
	GetUserByLogin(context.Context, string, string) (domain.User, error)
	UpdateUser(context.Context, domain.User) (domain.User, error)
	DeleteUser(context.Context, int) error
	AddContact(ctx context.Context, user domain.User, userID int) error
	FindUsers(ctx context.Context, search string, start int, stop int) ([]domain.User, error)
}

type IChatService interface {
	CreateChat(context.Context, domain.Chat) (domain.Chat, error)
	AddUserToChat(ctx context.Context, userID int, chatID int, role string) error
	GetChatByNameAndType(ctx context.Context, name, chatType string) (domain.Chat, error)
	GetChatsByUser(ctx context.Context, userID int) ([]domain.Chat, error)
	GetUserIDsByChatID(ctx context.Context, chatID int) ([]int, error)
	GetChatMember(ctx context.Context, userID int, chatID int) (domain.ChatMember, bool)
	GetUsersByChatID(ctx context.Context, chatID int) ([]domain.User, error)
}

type IMessageService interface {
	SaveMsg(ctx context.Context, msg domain.Message) error
	GetMessagesByChatID(ctx context.Context, chatID, limit, offset int) ([]domain.Message, error)
	GetChatMessages(ctx context.Context, chatID, initialMsgID, before, after int) ([]domain.Message, error)
}

type IAccountService interface {
	CreateAccount(context.Context, domain.Account) (domain.Account, error)
	NewUserToNewAccount(ctx context.Context, userID int, accountID int) error
	AddUserToAccount(ctx context.Context, userID int, accountID int, inviterID int, role string) error
	GetAccountByUserID(ctx context.Context, userID int) (int, error)
	GetContactsByAccount(ctx context.Context, accID int) ([]int, error)
}
