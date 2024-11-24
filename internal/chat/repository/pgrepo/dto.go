package pgrepo

import (
	"github.com/KozlovNikolai/pfp/internal/chat/domain"
	"github.com/KozlovNikolai/pfp/internal/chat/repository/models"
)

func domainToUser(user domain.User) models.User {
	return models.User{
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

func userToDomain(user models.User) domain.User {
	return domain.NewUser(domain.NewUserData{
		ID:        user.ID,
		UserExtID: user.UserExtID,
		Login:     user.Login,
		Password:  user.Password,
		Account:   user.Account,
		Token:     user.Token,
		Name:      user.Name,
		Surname:   user.Surname,
		Email:     user.Email,
		UserType:  user.UserType,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	})
}

func domainToChat(chat domain.Chat) models.Chat {
	return models.Chat{
		Id:            chat.ID(),
		Name:          chat.Name(),
		OwnerID:       chat.OwnerID(),
		ChatType:      chat.ChatType(),
		LastChatMsgID: chat.LastMsgID(),
		CreatedAt:     chat.CreatedAt(),
		UpdatedAt:     chat.UpdatedAt(),
	}
}

func chatToDomain(chat models.Chat) domain.Chat {
	return domain.NewChat(domain.NewChatData{
		ID:            chat.Id,
		Name:          chat.Name,
		OwnerID:       chat.OwnerID,
		ChatType:      chat.ChatType,
		LastChatMsgID: chat.LastChatMsgID,
	})
}
