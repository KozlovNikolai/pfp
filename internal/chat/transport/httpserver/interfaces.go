package httpserver

import (
	"context"

	"github.com/KozlovNikolai/pfp/internal/chat/domain"
	"github.com/google/uuid"
)

// ITokenService is a token service
type ITokenService interface {
	GenerateTokenForRegisteredUsers(ctx context.Context, user domain.UserChat) (string, error)
	GenerateToken(ctx context.Context, account, login, password string) (domain.UserChat, string, error)
	GetUser(ctx context.Context, token string) (domain.UserChat, error)
	GetPubsubToken(ctx context.Context, user domain.UserChat) (uuid.UUID, error)
}

type IStateService interface {
}

type IUserChatService interface {
	CreateUser(context.Context, domain.UserChat) (domain.UserChat, error)
	RegisterUser(context.Context, domain.UserChat) (domain.UserChat, error)
	GetUsers(context.Context, domain.UserChat, int, int) ([]domain.UserChat, error)
	GetUserByID(context.Context, int) (domain.UserChat, error)
	GetUserByExtID(context.Context, string, string) (domain.UserChat, error)
	GetUserByLogin(context.Context, string, string) (domain.UserChat, error)
	UpdateUser(context.Context, domain.UserChat) (domain.UserChat, error)
	DeleteUser(context.Context, int) error
}
