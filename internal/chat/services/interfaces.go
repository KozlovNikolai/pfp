// Package services ...
package services

import (
	"context"

	"github.com/KozlovNikolai/pfp/internal/chat/domain"
)

// IUserRepository is ...
type IUserRepository interface {
	CreateUserChat(context.Context, domain.UserChat) (domain.UserChat, error)
	GetUsers(context.Context, int, int) ([]domain.UserChat, error)
	GetUserByID(context.Context, int) (domain.UserChat, error)
	GetUserByExtID(context.Context, string, string) (domain.UserChat, error)
	GetUserByLogin(context.Context, string, string) (domain.UserChat, error)
	UpdateUser(context.Context, domain.UserChat) (domain.UserChat, error)
	DeleteUser(context.Context, int) error
}

type IStateRepository interface {
	GetState(context.Context, domain.UserChat) (domain.State, error)
}
