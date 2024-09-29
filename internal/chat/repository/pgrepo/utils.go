package pgrepo

import (
	"github.com/KozlovNikolai/pfp/internal/chat/domain"
	"github.com/KozlovNikolai/pfp/internal/chat/repository/models"
)

func domainToUser(user domain.User) models.User {
	return models.User{
		ID:       user.ID(),
		Login:    user.Login(),
		Password: user.Password(),
		Role:     user.Role(),
		Token:    user.Token(),
	}
}

func userToDomain(user models.User) domain.User {
	return domain.NewUser(domain.NewUserData{
		ID:       user.ID,
		Login:    user.Login,
		Password: user.Password,
		Role:     user.Role,
		Token:    user.Token,
	})
}
