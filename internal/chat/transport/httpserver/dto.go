package httpserver

import (
	"fmt"

	"github.com/KozlovNikolai/pfp/internal/chat/domain"
	"github.com/gin-gonic/gin"
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
func getUserFromContext(ctx *gin.Context) (domain.User, error) {
	contextUser, exists := ctx.Get(ctxKey("user").String())
	if !exists {
		return domain.User{}, domain.ErrNoUserInContext
	}
	// contextUser := ctx.Value("user")

	fmt.Printf("contextUser %v\n", contextUser)
	if contextUser == nil {
		return domain.User{}, domain.ErrNoUserInContext
	}
	user, ok := contextUser.(domain.User)
	if !ok {
		return domain.User{}, domain.ErrNoUserInContext
	}
	return user, nil
}
