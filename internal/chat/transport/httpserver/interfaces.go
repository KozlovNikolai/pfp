package httpserver

import (
	"context"

	"github.com/KozlovNikolai/pfp/internal/chat/domain"
)

// // IUserService is ...
// type IUserService interface {
// 	CreateUser(context.Context, domain.UserChat) (domain.UserChat, error)
// 	GetUsers(context.Context, int, int) ([]domain.UserChat, error)
// 	GetUserByID(context.Context, int) (domain.UserChat, error)
// 	GetUserByLogin(context.Context, string) (domain.UserChat, error)
// 	UpdateUser(context.Context, domain.UserChat) (domain.UserChat, error)
// 	DeleteUser(context.Context, int) error
// }

// ITokenService is a token service
type ITokenService interface {
	GenerateToken(ctx context.Context, account, login, password string) (string, error)
	GetUser(token string) (domain.UserChat, error)
}

type IStateService interface {
}

type IUserChatService interface {
	CreateUser(context.Context, domain.UserChat) (domain.UserChat, error)
	RegisterUser(context.Context, domain.UserChat) (domain.UserChat, error)
	GetUsers(context.Context, int, int) ([]domain.UserChat, error)
	GetUserByID(context.Context, int) (domain.UserChat, error)
	GetUserByExtID(context.Context, string, string) (domain.UserChat, error)
	GetUserByLogin(context.Context, string, string) (domain.UserChat, error)
	UpdateUser(context.Context, domain.UserChat) (domain.UserChat, error)
	DeleteUser(context.Context, int) error
}
