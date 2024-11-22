package httpserver

import (
	"strconv"

	"github.com/KozlovNikolai/pfp/internal/chat/constants"
	"github.com/KozlovNikolai/pfp/internal/chat/domain"
	"github.com/KozlovNikolai/pfp/internal/chat/transport/httpserver/middlewares"
)

// #########################################################
// func toResponseUser(user domain.User) UserResponse {
// 	return UserResponse{
// 		ID:       user.ID(),
// 		Login:    user.Login(),
// 		Password: user.Password(),
// 		Role:     user.Role(),
// 		Token:    user.Token(),
// 	}
// }

func toDomainUserChat(user UserRequest) domain.UserChat {
	return domain.NewUserChat(domain.NewUserChatData{
		UserExtID: "0",
		Login:     user.Login,
		Password:  user.Password,
		Account:   user.Account,
		Token:     "",
		Name:      "",
		Surname:   "",
		Email:     user.Login,
		UserType:  "regular",
		CreatedAt: 0,
		UpdatedAt: 0,
	})
}

// #########################################################

func toResponseUserChat(user domain.UserChat) UserChatResponse {
	return UserChatResponse{
		ID:        user.ID(),
		UserExtID: user.UserExtID(),
		Login:     user.Login(),
		Password:  user.Password(),
		Account:   user.Account(),
		Token:     user.Token(),
		Name:      user.Name(),
		Surname:   user.Surname(),
		Email:     user.Email(),
		UserType:  user.UserType(),
		CreatedAt: user.CreatedAt(),
		UpdatedAt: user.UpdatedAt(),
	}
}

func toDomainUserSputnik(user middlewares.ReceiveUserSputnik) domain.UserChat {
	return domain.NewUserChat(domain.NewUserChatData{
		UserExtID: strconv.Itoa(user.Payload.UserID),
		Login:     user.Payload.Email,
		Password:  "",
		Account:   constants.Account_name_magnum,
		Token:     "",
		Name:      user.Payload.Name,
		Surname:   user.Payload.Surname,
		Email:     user.Payload.Email,
		UserType:  "",
		CreatedAt: 0,
		UpdatedAt: 0,
	})
}
