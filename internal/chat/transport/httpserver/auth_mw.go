package httpserver

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/KozlovNikolai/pfp/internal/pkg/config"
	"github.com/KozlovNikolai/pfp/internal/pkg/utils"
)

// AuthorizationHeader ...
// BearerPrefix ...
const (
	AuthorizationHeader        = "Authorization"
	BearerPrefix               = "Bearer "
	UserIdToken                = "user-id"
	HeaderApplication          = "Application"
	sputnikUrl          string = "https://api.sputnik-monitor.ru/api/v1/auth/login"
)

// CheckAdmin ...
func (h HTTPServer) CheckAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		fmt.Println("CheckAdmin----------------------->>>")

		token := c.GetHeader(AuthorizationHeader)
		token = strings.TrimPrefix(token, BearerPrefix)
		user, err := h.tokenService.GetUser(token)
		if err != nil {
			c.AbortWithStatusJSON(
				http.StatusInternalServerError,
				gin.H{"check-admin-validate-token": err.Error()},
			)
			return
		}
		if user.Login() == "" {
			c.AbortWithStatusJSON(
				http.StatusInternalServerError,
				gin.H{"check-admin-invalid-token": ""},
			)
			return
		}
		if user.Role() != config.AdminRole {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"not-admin": ""})
			return
		}
		c.Set(utils.CtxKey("user").String(), user)
		// v, e := c.Get(ctxKey("user").String())
		// if !e {
		// 	fmt.Println("значение в конетексте не существует")
		// }
		// fmt.Printf("user in context = %v\n", v)
		c.Next()
	}
}

// CheckAuthorizedUser ...
func (h HTTPServer) CheckAuthorizedUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		fmt.Println("CheckAuthorizedUser---------------------->>>")

		token := c.GetHeader(AuthorizationHeader)
		fmt.Printf("token = %s\n", token)
		token = strings.TrimPrefix(token, BearerPrefix)
		user, err := h.tokenService.GetUser(token)
		if err != nil {
			c.AbortWithStatusJSON(
				http.StatusInternalServerError,
				gin.H{"check-auth-validate-token": err.Error()},
			)
			return
		}
		fmt.Printf("user = %+v\n", user)
		if user.Login() == "" {
			c.AbortWithStatusJSON(
				http.StatusInternalServerError,
				gin.H{"check-auth-invalid-token": ""},
			)
			return
		}

		c.Set(utils.CtxKey("user").String(), user)

		c.Next()
	}
}
