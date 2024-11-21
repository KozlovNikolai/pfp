// Package httpserver ...
package httpserver

import (
	"fmt"
	"net/http"

	"github.com/KozlovNikolai/pfp/internal/chat/domain"
	"github.com/KozlovNikolai/pfp/internal/chat/transport/httpserver/middlewares"
	"github.com/KozlovNikolai/pfp/internal/pkg/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const (
	invaldRequest = "invalid-request"
)

// SignUp is ...
// SignUpTags		godoc
// @Summary				Загеристрироваться.
// @Description			Sign up a new user in the system.
// @Param				UserRequest body httpserver.UserRequest true "Create user. Логин указывается в формате электронной почты. Пароль не меньше 6 символов. Роль: super или regular"
// @Produce				application/json
// @Tags				Auth
// @Success				201 {object} httpserver.UserResponse
// @failure				400 {string} err.Error()
// @failure				500 {string} string "error-to-create-domain-user"
// @Router				/signup [post]
func (h HTTPServer) SignUp(c *gin.Context) {
	var userRequest UserRequest
	var err error
	if err = c.ShouldBindJSON(&userRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"invalid-json": err.Error()})

		return
	}

	if err = userRequest.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{invaldRequest: err.Error()})
		return
	}

	// userRequest.Password, err = hashPassword(userRequest.Password)
	// if err != nil {
	// 	c.JSON(http.StatusBadRequest, gin.H{"error-hashing-password": err.Error()})
	// 	return
	// }

	domainUser := toDomainUser(userRequest)

	createdUser, err := h.userService.CreateUser(c, domainUser)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error service User": err.Error()})
		return
	}
	response := toResponseUser(createdUser)
	c.JSON(http.StatusCreated, response)
}

// SignIn is ...
// SignInTags		godoc
// @Summary				Авторизоваться.
// @Description			Sign in as an existing user.
// @Param				UserRequest body httpserver.UserRequest true "SignIn user. Логин указывается в формате электронной почты. Пароль не меньше 6 символов. Роль: super или regular"
// @Produce				application/json
// @Tags				Auth
// @Success				200 {string} token
// @failure				400 {string} err.Error()
// @failure				500 {string} err.Error()
// @Router				/signin [post]
func (h HTTPServer) SignIn(c *gin.Context) {
	var userRequest UserRequest
	if err := c.ShouldBindJSON(&userRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"invalid-json": err.Error()})
		return
	}

	if err := userRequest.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{invaldRequest: err.Error()})
		return
	}

	token, err := h.tokenService.GenerateToken(c, userRequest.Login, userRequest.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error-generated-token": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"token": token})

}

// LoginUserByToken is ...
// LoginUserByTokenTags		godoc
// @Summary				Авторизоваться по токену.
// @Description			Logging in as an existing user from remote Auth service.
// @Produce				application/json
// @Tags				Auth
// @Success				200 {string} pubsub_token
// @failure				400 {string} err.Error()
// @failure				500 {string} err.Error()
// @Router				/auth/login [get]
func (h HTTPServer) LoginUserByToken(c *gin.Context) {

	userSputnik, err := utils.GetDataFromContext[middlewares.ReceiveUserSputnik](c, "user_sputnik")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": domain.ErrNoUserInContext.Error()})
		return
	}

	tokenWS := uuid.New()

	fmt.Println()
	fmt.Printf("Result: %+v\n", userSputnik)
	fmt.Println()

	c.JSON(http.StatusOK, gin.H{"websocket token": tokenWS})
}

// SignOut ...
func (h HTTPServer) SignOut(c *gin.Context) {
	_ = c
}
