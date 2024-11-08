package httpserver

import (
	"context"

	"github.com/KozlovNikolai/pfp/internal/chat/domain"
)

// IUserService is ...
type IUserService interface {
	CreateUser(context.Context, domain.User) (domain.User, error)
	GetUsers(context.Context, int, int) ([]domain.User, error)
	GetUserByID(context.Context, int) (domain.User, error)
	GetUserByLogin(context.Context, string) (domain.User, error)
	UpdateUser(context.Context, domain.User) (domain.User, error)
	DeleteUser(context.Context, int) error
}

// ITokenService is a token service
type ITokenService interface {
	GenerateToken(ctx context.Context, login string, password string) (string, error)
	GetUser(token string) (domain.User, error)
}
