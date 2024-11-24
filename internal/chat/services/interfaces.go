// Package services ...
package services

import (
	"context"

	"github.com/KozlovNikolai/pfp/internal/chat/domain"
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
	GetState(context.Context, domain.User) (domain.State, error)
	SetState(context.Context, domain.User) (domain.State, error)
}

type IChatRepository interface {
	CreateChat(context.Context, domain.Chat) (domain.Chat, error)
	AddUserToChat(context.Context, int, int) ([]domain.Chat, error)
	GetChatByNameAndType(context.Context, string, string) (domain.Chat, error)
	GetChatsByUser(ctx context.Context, userID int) ([]domain.Chat, error)
}
