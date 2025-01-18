// Package httpserver ...
package httpserver

import (
	"fmt"
	"log"
	"net/http"

	"alfachat/internal/chat/constants"
	"alfachat/internal/chat/domain"
	"alfachat/internal/chat/transport/httpserver/middlewares"
	"alfachat/internal/pkg/utils"

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
	// var accountRequest AccountRequest
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
	domainAccount := toDomainAccount(AccountRequest{
		Name: domainUser.Name() + " " + domainUser.Surname(),
	})
	createdAccount, err := h.accountService.CreateAccount(c, domainAccount)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	log.Printf("Account created, id: %d, name: %s", createdAccount.ID(), createdAccount.Name())

	createdUser, err := h.userService.CreateUser(c, domainUser)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error service User": err.Error()})
		return
	}

	err = h.accountService.NewUserToNewAccount(c, createdUser.ID(), createdAccount.ID())
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	log.Printf("User id: %d binded to Account id: %d", createdUser.ID(), createdAccount.ID())

	var chats []ChatResponse
	if !(createdUser.Login() == constants.AdminEmaleLogin && createdUser.Profile() == constants.SystemProfile) {

		chat, err := h.chatService.GetChatByNameAndType(c, constants.SystemChatName, constants.SystemChatType)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"failure get common system chat id": err.Error()})
			return
		}
		// log.Printf("chat response: %+v", chat)
		err = h.chatService.AddUserToChat(c, createdUser.ID(), chat.ID(), constants.ChatRoleRegular)
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

	_, ok := h.stateService.GetState(c, createdUser.ID())
	status := "offline"
	if ok {
		status = "online"
	}

	response := toResponseUser(createdUser, status)
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

	user, token, err := h.tokenService.GenerateToken(c, userRequest.Profile, userRequest.Login, userRequest.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error-generated-token": err.Error()})
		return
	}

	pubsub, err := h.tokenService.GetPubsubToken(c, user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error-generated-pubsub": err.Error()})
		return
	}

	h.stateService.SetState(c, user.ID(), pubsub, nil)
	var currentChatID int
	if !(user.Login() == constants.AdminEmaleLogin && user.Profile() == constants.SystemProfile) {
		chat, err := h.chatService.GetChatByNameAndType(c, constants.SystemChatName, constants.SystemChatType)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"failure get common system chat id": err.Error()})
			return
		}
		ok := h.stateService.SetCurrentChat(c, user.ID(), pubsub, chat.ID())
		if !ok {
			log.Printf("Не удалось установить текущий чат")
		}
		currentChatID = chat.ID()
	}

	state, ok := h.stateService.GetState(c, user.ID())
	log.Printf("user ID: %d\n", user.ID())
	for _, v := range state.Connects {
		fmt.Printf("pubsub: %v, chat: %d\n", v.Pubsub, v.CurrentChat)
	}
	log.Printf("ok: %+v\n", ok)

	var chats []ChatResponse
	chatsDomain, err := h.chatService.GetChatsByUser(c, user.ID())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"failure get chats by user id": err.Error()})
		return
	}
	for _, chatDomain := range chatsDomain {
		chats = append(chats, toResponseChat(chatDomain))
	}

	_, ok = h.stateService.GetState(c, user.ID())
	status := "offline"
	if ok {
		status = "online"
	}

	c.JSON(http.StatusOK, gin.H{
		"pubsub":          pubsub,
		"token":           token,
		"user":            toResponseUser(user, status),
		"chats":           chats,
		"current_chat_id": currentChatID,
	})
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

		domainAccount := toDomainAccount(AccountRequest{
			Name: domainUser.Name() + " " + domainUser.Surname(),
		})
		createdAccount, err := h.accountService.CreateAccount(c, domainAccount)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		log.Printf("Account created, id: %d, name: %s", createdAccount.ID(), createdAccount.Name())

		// createdUser, err := h.userService.CreateUser(c, domainUser)
		// if err != nil {
		// 	c.JSON(http.StatusInternalServerError, gin.H{"error service User": err.Error()})
		// 	return
		// }

		err = h.accountService.NewUserToNewAccount(c, registeredUser.ID(), createdAccount.ID())
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		log.Printf("User id: %d binded to Account id: %d", registeredUser.ID(), createdAccount.ID())

		chat, err := h.chatService.GetChatByNameAndType(c, constants.SystemChatName, constants.SystemChatType)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"failure get common system chat id": err.Error()})
			return
		}
		// log.Printf("chat response: %+v", chat)
		err = h.chatService.AddUserToChat(c, registeredUser.ID(), chat.ID(), constants.ChatRoleRegular)
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

	h.stateService.SetState(c, registeredUser.ID(), pubsub, nil)
	sysChat, err := h.chatService.GetChatByNameAndType(c, constants.SystemChatName, constants.SystemChatType)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ok := h.stateService.SetCurrentChat(c, registeredUser.ID(), pubsub, sysChat.ID())
	if !ok {
		log.Printf("Не удалось установить текущий чат")
	}
	state, ok := h.stateService.GetState(c, registeredUser.ID())
	log.Printf("user ID: %d, state: %+v\n", registeredUser.ID(), state)
	log.Printf("ok: %+v\n", ok)

	_, ok = h.stateService.GetState(c, registeredUser.ID())
	status := "offline"
	if ok {
		status = "online"
	}

	response := toResponseUser(registeredUser, status)
	c.JSON(http.StatusOK, gin.H{"pubsub": pubsub, "token": token, "user": response, "chats": chats})
}

// SignOut ...
func (h HTTPServer) SignOut(c *gin.Context) {
	pubsub, err := uuid.Parse(c.Param("pubsub"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	userCtx, err := utils.GetDataFromContext[domain.User](c, "user")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": domain.ErrNoUserInContext.Error()})
		return
	}

	isError := false
	state, ok := h.stateService.DeleteConnFromState(c, userCtx.ID(), pubsub)
	if !ok {
		isError = true
		log.Printf("DeleteConnFromState - failure")
	}
	ok = h.wsHandler.Unsubscribe(pubsub)
	if !ok {
		isError = true
		log.Printf("Unsubscribe - failure")
	}
	if isError {
		c.JSON(http.StatusOK, gin.H{"status": "connection not found", "state": toResponseState(state)})
	}
	c.JSON(http.StatusOK, gin.H{"status": "disconnected"})
}
