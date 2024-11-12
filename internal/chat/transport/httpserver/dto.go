package httpserver

import (
	"github.com/KozlovNikolai/pfp/internal/chat/domain"
)

// #########################################################
func toResponseUser(user domain.User) UserResponse {
	return UserResponse{
		ID:       user.ID(),
		Login:    user.Login(),
		Password: user.Password(),
		Role:     user.Role(),
		Token:    user.Token(),
	}
}

func toDomainUser(user UserRequest) domain.User {
	return domain.NewUser(domain.NewUserData{
		Login:    user.Login,
		Password: user.Password,
	})
}

// #########################################################
