package httpserver

import (
	"strconv"

	"github.com/KozlovNikolai/pfp/internal/chat/constants"
	"github.com/KozlovNikolai/pfp/internal/chat/domain"
	"github.com/KozlovNikolai/pfp/internal/chat/transport/httpserver/middlewares"
)

func toDomainUser(user UserRequest) domain.User {
	return domain.NewUser(domain.NewUserData{
		UserExtID: "0",
		Login:     user.Login,
		Password:  user.Password,
		Account:   user.Account,
		Token:     "",
		Name:      user.Name,
		Surname:   user.Surname,
		Email:     user.Login,
		UserType:  "regular",
		CreatedAt: 0,
		UpdatedAt: 0,
	})
}

// #########################################################

func toResponseUser(user domain.User) UserResponse {
	return UserResponse{
		ID:        user.ID(),
		UserExtID: user.UserExtID(),
		Login:     user.Login(),
		// Password:  user.Password(),
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

func toDomainUserFromUserSputnik(user middlewares.ReceiveUserSputnik) domain.User {
	return domain.NewUser(domain.NewUserData{
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

func toDomainChat(chat ChatCreateRequest) domain.Chat {
	return domain.NewChat(domain.NewChatData{
		OwnerID:       chat.OwnerID,
		Name:          chat.Name,
		ChatType:      chat.ChatType,
		LastChatMsgID: 0,
	})
}

func toResponseChat(chat domain.Chat) ChatResponse {
	return ChatResponse{
		Id:            chat.ID(),
		Name:          chat.Name(),
		OwnerID:       chat.OwnerID(),
		ChatType:      chat.ChatType(),
		LastChatMsgID: chat.LastMsgID(),
		Contacts:      chat.Contacts(),
		CreatedAt:     chat.CreatedAt(),
		UpdatedAt:     chat.UpdatedAt(),
	}
}
