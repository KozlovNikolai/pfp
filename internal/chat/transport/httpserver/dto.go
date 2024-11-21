package httpserver

import (
	"github.com/KozlovNikolai/pfp/internal/chat/domain"
	"github.com/KozlovNikolai/pfp/internal/chat/transport/httpserver/middlewares"
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

func toResponseUserChat(user domain.UserChat) UserChatResponse {
	return UserChatResponse{
		ID:        user.GetID(),
		UserExtID: user.GetUserExtID(),
		Name:      user.GetName(),
		Surname:   user.GetSurname(),
		Email:     user.GetEmail(),
		UserType:  user.GetUserType(),
		CreatedAt: user.GetCreatedAt(),
		UpdatedAt: user.GetUpdatedAt(),
	}
}

func toDomainUserChat(user middlewares.ReceiveUserSputnik) domain.UserChat {
	return domain.NewUserChat(domain.NewUserChatData{
		UserExtID: user.Payload.UserID,
		Name:      user.Payload.Name,
		Surname:   user.Payload.Surname,
		Email:     user.Payload.Email,
		UserType:  user.Payload.Lang,
		CreatedAt: 0,
		UpdatedAt: 0,
	})
}
