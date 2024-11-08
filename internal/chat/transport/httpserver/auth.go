// Package httpserver ...
package httpserver

import (
	"net/http"

	"github.com/gin-gonic/gin"
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

// SignOut ...
func (h HTTPServer) SignOut(c *gin.Context) {
	_ = c
}
