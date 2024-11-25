// Package httpserver ...
package httpserver

import (
	"log"
	"net/http"

	"github.com/KozlovNikolai/pfp/internal/chat/domain"
	"github.com/KozlovNikolai/pfp/internal/chat/transport/httpserver/middlewares"
	"github.com/KozlovNikolai/pfp/internal/pkg/utils"
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
	var chats []ChatResponse
	if !(createdUser.Login() == "root@admin.ru" && createdUser.Account() == "system") {

		chat, err := h.chatService.GetChatByNameAndType(c, "common chat", "system")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"failure get common system chat id": err.Error()})
			return
		}
		// log.Printf("chat response: %+v", chat)
		err = h.chatService.AddUserToChat(c, createdUser.ID(), chat.ID())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"failure add user to system chat": err.Error()})
			return
		}
		chatsDomain, err := h.chatService.GetChatsByUser(c, createdUser.ID())

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		chatsResponse := make([]ChatResponse, len(chatsDomain))
		for i, chatDomain := range chatsDomain {
			chatsResponse[i] = toResponseChat(chatDomain)
		}
	}

	response := toResponseUser(createdUser)
	// c.JSON(http.StatusCreated, response, chats)
	c.JSON(http.StatusCreated, gin.H{"user": response, "chats": chats})
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

	user, token, err := h.tokenService.GenerateToken(c, userRequest.Account, userRequest.Login, userRequest.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error-generated-token": err.Error()})
		return
	}

	pubsub, err := h.tokenService.GetPubsubToken(c, user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error-generated-pubsub": err.Error()})
		return
	}
	var chats []ChatResponse
	chatsDomain, err := h.chatService.GetChatsByUser(c, user.ID())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"failure get chats by user id": err.Error()})
		return
	}
	for _, chatDomain := range chatsDomain {
		chats = append(chats, toResponseChat(chatDomain))
	}

	c.JSON(http.StatusOK, gin.H{"pubsub": pubsub, "token": token, "user": toResponseUser(user), "chats": chats})
}

// LoginUserByTokenSputnik is ...
// LoginUserByTokenTags		godoc
// @Summary				Авторизоваться по токену.
// @Description			Logging in as an existing user from remote Auth service.
// @Produce				application/json
// @Tags				Auth
// @Success				200 {string} pubsub_token
// @failure				400 {string} err.Error()
// @failure				500 {string} err.Error()
// @Router				/auth/login [get]
func (h HTTPServer) LoginUserByTokenSputnik(c *gin.Context) {

	userSputnik, err := utils.GetDataFromContext[middlewares.ReceiveUserSputnik](c, "user_sputnik")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": domain.ErrNoUserInContext.Error()})
		return
	}
	domainUser := toDomainUserFromUserSputnik(userSputnik)

	registeredUser, isNew, err := h.userService.RegisterUser(c, domainUser)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error service User": err.Error()})
		return
	}
	// log.Printf("isNew: %v", isNew)
	var chats []ChatResponse
	log.Printf("is New: %v", isNew)
	if isNew {
		chat, err := h.chatService.GetChatByNameAndType(c, "common chat", "system")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"failure get common system chat id": err.Error()})
			return
		}
		// log.Printf("chat response: %+v", chat)
		err = h.chatService.AddUserToChat(c, registeredUser.ID(), chat.ID())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"failure add user to system chat": err.Error()})
			return
		}
		chatsDomain, err := h.chatService.GetChatsByUser(c, registeredUser.ID())

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		for _, chatDomain := range chatsDomain {
			chats = append(chats, toResponseChat(chatDomain))
		}
	} else {
		chatsDomain, err := h.chatService.GetChatsByUser(c, registeredUser.ID())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"failure get chats by user id": err.Error()})
			return
		}
		for _, chatDomain := range chatsDomain {
			chats = append(chats, toResponseChat(chatDomain))
		}
	}
	token, err := h.tokenService.GenerateTokenForRegisteredUsers(c, registeredUser)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error-generated-token": err.Error()})
		return
	}

	pubsub, err := h.tokenService.GetPubsubToken(c, registeredUser)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error-generated-pubsub": err.Error()})
		return
	}

	response := toResponseUser(registeredUser)
	c.JSON(http.StatusOK, gin.H{"pubsub": pubsub, "token": token, "user": response, "chats": chats})
}

// SignOut ...
func (h HTTPServer) SignOut(c *gin.Context) {
	_ = c
}
