// Package middlewares ...
package middlewares

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"

	"github.com/KozlovNikolai/pfp/internal/pkg/config"
)

// AuthorizationHeader ...
// BearerPrefix ...
const (
	AuthorizationHeader = "Authorization"
	BearerPrefix        = "Bearer "
)

// Claims is ...
type Claims struct {
	AuthID    int    `json:"auth_id"`
	AuthLogin string `json:"auth_login"`
	AuthRole  string `json:"auth_role"`
	jwt.RegisteredClaims
}

// AuthMiddleware ...
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader(AuthorizationHeader)
		if authHeader == "" {
			c.AbortWithStatusJSON(
				http.StatusUnauthorized,
				gin.H{"error": "Authorization header is required"},
			)
			return
		}

		tokenString := authHeader[len(BearerPrefix):]
		claims := &Claims{}
		tkn, err := jwt.ParseWithClaims(
			tokenString,
			claims,
			func(token *jwt.Token) (interface{}, error) {
				_ = token
				return config.JwtKey, nil
			},
		)
		if err != nil {
			if err == jwt.ErrSignatureInvalid {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
				return
			}
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Bad request"})
			return
		}

		if !tkn.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}
		fmt.Printf("auth_id = %v\n", claims.AuthID)
		fmt.Printf("auth_login = %v\n", claims.AuthLogin)
		fmt.Printf("auth_role = %v\n", claims.AuthRole)

		// Сохранение имени пользователя в контексте запроса
		c.Set("auth_id", claims.AuthID)
		c.Set("auth_login", claims.AuthLogin)
		c.Set("auth_role", claims.AuthRole)
		c.Next()
	}
}
