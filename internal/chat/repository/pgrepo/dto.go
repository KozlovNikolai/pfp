package pgrepo

import (
	"github.com/KozlovNikolai/pfp/internal/chat/domain"
	"github.com/KozlovNikolai/pfp/internal/chat/repository/models"
)

func domainToUserChat(user domain.UserChat) models.UserChat {
	return models.UserChat{
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

func userChatToDomain(user models.UserChat) domain.UserChat {
	return domain.NewUserChat(domain.NewUserChatData{
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