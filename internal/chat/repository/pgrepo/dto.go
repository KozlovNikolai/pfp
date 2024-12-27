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
		Profile:   user.Profile(),
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
		Profile:   user.Profile,
		Name:      user.Name,
		Surname:   user.Surname,
		Email:     user.Email,
		UserType:  user.UserType,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	})
}

func domainToAccount(account domain.Account) models.Account {
	return models.Account{
		ID:        account.ID(),
		Name:      account.Name(),
		CreatedAt: account.CreatedAt(),
		UpdatedAt: account.UpdatedAt(),
	}
}

func accountToDomain(account models.Account) domain.Account {
	return domain.NewAccount(domain.NewAccountData{
		ID:        account.ID,
		Name:      account.Name,
		CreatedAt: account.CreatedAt,
		UpdatedAt: account.UpdatedAt,
	})
}

func domainToChat(chat domain.Chat) models.Chat {
	return models.Chat{
		Id:            chat.ID(),
		Name:          chat.Name(),
		AccountID:     chat.AccountID(),
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
		AccountID:     chat.AccountID,
		ChatType:      chat.ChatType,
		LastChatMsgID: chat.LastChatMsgID,
	})
}

func domainToMessage(msg domain.Message) models.Message {
	return models.Message{
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

func messageToDomain(msg models.Message) domain.Message {
	return domain.NewMessage(domain.NewMessageData{
		Id:        msg.Id,
		SenderID:  msg.SenderID,
		ChatID:    msg.ChatID,
		MsgType:   msg.MsgType,
		Text:      msg.Text,
		IsDeleted: msg.IsDeleted,
		CreatedAt: msg.CreatedAt,
		UpdatedAt: msg.UpdatedAt,
	})
}

func domainToChatMember(chatMember domain.ChatMember) models.ChatMember {
	return models.ChatMember{
		Id:            chatMember.ID(),
		ChatID:        chatMember.ChatID(),
		UserID:        chatMember.UserID(),
		Role:          chatMember.Role(),
		LastReadMsgID: chatMember.LastReadMsgID(),
		Notifications: chatMember.Notifications(),
		CreatedAt:     chatMember.CreatedAt(),
		UpdatedAt:     chatMember.UpdatedAt(),
	}
}

func chatMemberToDomain(chatMember models.ChatMember) domain.ChatMember {
	return domain.NewChatMember(domain.NewChatMemberData{
		Id:            chatMember.Id,
		ChatID:        chatMember.ChatID,
		UserID:        chatMember.UserID,
		Role:          chatMember.Role,
		LastReadMsgID: chatMember.LastReadMsgID,
		Notifications: chatMember.Notifications,
		CreatedAt:     chatMember.CreatedAt,
		UpdatedAt:     chatMember.UpdatedAt,
	})
}

func domainToContact(contact domain.Contact) models.Contact {
	return models.Contact{
		Id:        contact.ID(),
		AccountID: contact.AccountID(),
		UserID:    contact.UserID(),
		Name:      contact.Name(),
		Surname:   contact.Surname(),
		Phone:     contact.Phone(),
		Email:     contact.Email(),
	}
}

func contactToDomain(contact models.Contact) domain.Contact {
	return domain.NewContact(domain.NewContactData{
		Id:        contact.Id,
		AccountID: contact.AccountID,
		UserID:    contact.UserID,
		Name:      contact.Name,
		Surname:   contact.Surname,
		Phone:     contact.Phone,
		Email:     contact.Email,
	})
}
