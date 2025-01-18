// Package middlewares ...
package middlewares

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

	"alfachat/internal/pkg/utils"
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

type ReceiveUserSputnik struct {
	Payload Payload `json:"payload"`
}

type Payload struct {
	Email         string        `json:"email" db:"email"`
	Name          string        `json:"name" db:"name"`
	Surname       string        `json:"surname" db:"surname"`
	UserID        int           `json:"id" db:"id"`
	TokensSputnik TokensSputnik `json:"tokens"`
	Lang          string        `json:"lang" db:"lang"`
	AccountId     int           `json:"account_id"`
}

type TokensSputnik struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func (ru *ReceiveUserSputnik) Validate() error {
	validate := validator.New(validator.WithRequiredStructEnabled())
	return validate.Struct(ru)
}

// // Claims is ...
// type Claims struct {
// 	AuthID    int    `json:"auth_id"`
// 	AuthLogin string `json:"auth_login"`
// 	AuthRole  string `json:"auth_role"`
// 	jwt.RegisteredClaims
// }

// // AuthMiddleware ...
// func AuthMiddleware() gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		authHeader := c.GetHeader(AuthorizationHeader)
// 		if authHeader == "" {
// 			c.AbortWithStatusJSON(
// 				http.StatusUnauthorized,
// 				gin.H{"error": "Authorization header is required"},
// 			)
// 			return
// 		}

// 		tokenString := authHeader[len(BearerPrefix):]
// 		claims := &Claims{}
// 		tkn, err := jwt.ParseWithClaims(
// 			tokenString,
// 			claims,
// 			func(token *jwt.Token) (interface{}, error) {
// 				_ = token
// 				return config.JwtKey, nil
// 			},
// 		)
// 		if err != nil {
// 			if err == jwt.ErrSignatureInvalid {
// 				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
// 				return
// 			}
// 			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Bad request"})
// 			return
// 		}

// 		if !tkn.Valid {
// 			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
// 			return
// 		}
// 		fmt.Printf("auth_id = %v\n", claims.AuthID)
// 		fmt.Printf("auth_login = %v\n", claims.AuthLogin)
// 		fmt.Printf("auth_role = %v\n", claims.AuthRole)

// 		// Сохранение имени пользователя в контексте запроса
// 		c.Set("auth_id", claims.AuthID)
// 		c.Set("auth_login", claims.AuthLogin)
// 		c.Set("auth_role", claims.AuthRole)
// 		c.Next()
// 	}
// }

func AuthSputnikMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		rawToken := c.GetHeader(AuthorizationHeader)
		userIdToken := c.GetHeader(UserIdToken)
		headerApplication := c.GetHeader(HeaderApplication)
		_ = userIdToken
		_ = headerApplication

		// fmt.Println("Get headers:")
		// fmt.Printf("rawToken: %s, userIdToken: %s, headerApplication: %s\n",
		// rawToken, userIdToken, headerApplication)

		authHeaders := map[string]string{
			AuthorizationHeader: rawToken,
		}

		resp, err := utils.DoRequest[ReceiveUserSputnik]("GET", sputnikUrl, nil, authHeaders)
		if err != nil {
			log.Println("Error DoRequest")
			c.JSON(http.StatusUnauthorized, gin.H{"error": err})
			return
		}

		c.Set(utils.CtxKey("user_sputnik").String(), resp)
		c.Next()
	}
}
