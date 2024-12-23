package httpserver

import (
	"github.com/KozlovNikolai/pfp/internal/chat/constants"
	"github.com/KozlovNikolai/pfp/internal/chat/domain"
	"github.com/KozlovNikolai/pfp/internal/chat/transport/httpserver/middlewares"
)

func toDomainUser(user UserRequest) domain.User {
	return domain.NewUser(domain.NewUserData{
		UserExtID: 0,
		Login:     user.Login,
		Password:  user.Password,
		Profile:   user.Profile,
		Name:      user.Name,
		Surname:   user.Surname,
		Email:     user.Login,
		UserType:  "regular",
		CreatedAt: 0,
		UpdatedAt: 0,
	})
}

func toDomainAccount(account AccountRequest) domain.Account {
	return domain.NewAccount(domain.NewAccountData{
		Name:      account.Name,
		CreatedAt: 0,
		UpdatedAt: 0,
	})
}

// #########################################################

func toResponseUser(user domain.User, status string) UserResponse {
	return UserResponse{
		ID:        user.ID(),
		UserExtID: user.UserExtID(),
		Login:     user.Login(),
		// Password:  user.Password(),
		Profile:   user.Profile(),
		Name:      user.Name(),
		Surname:   user.Surname(),
		Email:     user.Email(),
		UserType:  user.UserType(),
		CreatedAt: user.CreatedAt(),
		UpdatedAt: user.UpdatedAt(),
		Status:    status,
	}
}

func toResponseAccount(account domain.Account) AccountResponse {
	return AccountResponse{
		ID:        account.ID(),
		Name:      account.Name(),
		CreatedAt: account.CreatedAt(),
		UpdatedAt: account.UpdatedAt(),
	}
}

func toDomainUserFromUserSputnik(user middlewares.ReceiveUserSputnik) domain.User {
	return domain.NewUser(domain.NewUserData{
		UserExtID: user.Payload.UserID,
		Login:     user.Payload.Email,
		Password:  "",
		Profile:   constants.Account_name_magnum,
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
		AccountID:     chat.OwnerID,
		Name:          chat.Name,
		ChatType:      chat.ChatType,
		LastChatMsgID: 0,
	})
}

func toResponseChat(chat domain.Chat) ChatResponse {
	return ChatResponse{
		Id:            chat.ID(),
		Name:          chat.Name(),
		AccountID:     chat.AccountID(),
		ChatType:      chat.ChatType(),
		LastChatMsgID: chat.LastMsgID(),
		// Contacts:      chat.Contacts(),
		CreatedAt: chat.CreatedAt(),
		UpdatedAt: chat.UpdatedAt(),
	}
}

func toDomainMessage(msg SendMessageRequest) domain.Message {
	return domain.NewMessage(domain.NewMessageData{
		SenderID:  msg.SenderID,
		ChatID:    msg.ChatID,
		MsgType:   msg.MsgType,
		Text:      msg.Text,
		CreatedAt: msg.CreatedAt,
		UpdatedAt: msg.UpdatedAt,
	})
}

func toResponseMessage(msg domain.Message) MessageResponse {
	return MessageResponse{
		Id:        msg.ID(),
		SenderID:  msg.SenderID(),
		ChatID:    msg.ChatID(),
		MsgType:   msg.MsgType(),
		Text:      msg.Text(),
		IsDeleted: msg.IsDeleted(),
		CreatedAt: msg.CreatedAt(),
		UpdatedAt: msg.UpdatedAt(),
	}
}

func toResponseState(state domain.State) StateResponse {
	var resp StateResponse
	resp.UserID = state.UserID
	resp.Connects = make([]Connect, len(state.Connects))
	for index, connect := range state.Connects {
		resp.Connects[index].Pubsub = connect.Pubsub
		resp.Connects[index].CreatedAt = connect.CreatedAt
		resp.Connects[index].CurrentChat = connect.CurrentChat
		if connect.Conn == nil {
			resp.Connects[index].Conn = false
		} else {
			resp.Connects[index].Conn = true
		}

	}
	return resp
}
