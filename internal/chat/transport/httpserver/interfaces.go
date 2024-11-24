package httpserver

import (
	"context"

	"github.com/KozlovNikolai/pfp/internal/chat/domain"
	"github.com/google/uuid"
)

// ITokenService is a token service
type ITokenService interface {
	GenerateTokenForRegisteredUsers(ctx context.Context, user domain.User) (string, error)
	GenerateToken(ctx context.Context, account, login, password string) (domain.User, string, error)
	GetUser(ctx context.Context, token string) (domain.User, error)
	GetPubsubToken(ctx context.Context, user domain.User) (uuid.UUID, error)
}

type IStateService interface {
}

type IUserService interface {
	CreateUser(context.Context, domain.User) (domain.User, error)
	RegisterUser(context.Context, domain.User) (domain.User, bool, error)
	GetUsers(context.Context, domain.User, int, int) ([]domain.User, error)
	GetUserByID(context.Context, int) (domain.User, error)
	GetUserByExtID(context.Context, string, string) (domain.User, error)
	GetUserByLogin(context.Context, string, string) (domain.User, error)
	UpdateUser(context.Context, domain.User) (domain.User, error)
	DeleteUser(context.Context, int) error
}

type IChatService interface {
	CreateChat(context.Context, domain.Chat) (domain.Chat, error)
	AddUserToChat(context.Context, int, int) ([]domain.Chat, error)
	GetChatByNameAndType(context.Context, string, string) (domain.Chat, error)
	GetChatsByUser(ctx context.Context, userID int) ([]domain.Chat, error)
}
