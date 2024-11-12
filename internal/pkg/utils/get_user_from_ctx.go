package utils

import (
	"fmt"

	"github.com/KozlovNikolai/pfp/internal/chat/domain"
	"github.com/gin-gonic/gin"
)

type CtxKey string

func (c CtxKey) String() string {
	return string(c)
}
func GetUserFromContext(ctx *gin.Context) (domain.User, error) {
	contextUser, exists := ctx.Get(CtxKey("user").String())
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
