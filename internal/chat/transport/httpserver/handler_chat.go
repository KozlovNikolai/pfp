package httpserver

import (
	"log"
	"net/http"
	"strconv"

	"github.com/KozlovNikolai/pfp/internal/chat/domain"
	"github.com/KozlovNikolai/pfp/internal/pkg/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (h HTTPServer) CreateChat(c *gin.Context) {
	var chatCreateRequest ChatCreateRequest
	var err error
	if err = c.ShouldBindJSON(&chatCreateRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"invalid-json": err.Error()})
		return
	}

	if err = chatCreateRequest.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{invaldRequest: err.Error()})
		return
	}
	userCtx, err := utils.GetDataFromContext[domain.User](c, "user")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": domain.ErrNoUserInContext.Error()})
		return
	}
	// fmt.Printf("\nchatFromCtx: %+v\n\n", userCtx)
	chatCreateRequest.OwnerID = userCtx.ID()
	// fmt.Printf("\nchatRequest: %+v\n\n", chatCreateRequest)
	domainChat := toDomainChat(chatCreateRequest)
	// fmt.Printf("\nchatDomain: %+v\n\n", domainChat)
	createdChat, err := h.chatService.CreateChat(c, domainChat)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error service User": err.Error()})
		return
	}
	response := toResponseChat(createdChat)
	c.JSON(http.StatusCreated, response)
}

func (h HTTPServer) AddToChat(c *gin.Context) {
	userCtx, err := utils.GetDataFromContext[domain.User](c, "user")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": domain.ErrNoUserInContext.Error()})
		return
	}

	var addToChatRequest AddToChatRequest

	if err = c.ShouldBindJSON(&addToChatRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"invalid-json": err.Error()})
		return
	}

	if err = addToChatRequest.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{invaldRequest: err.Error()})
		return
	}

	err = h.chatService.AddUserToChat(c, addToChatRequest.UserID, addToChatRequest.ChatID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error service User": err.Error()})
		return
	}
	chatsDomain, err := h.chatService.GetChatsByUser(c, userCtx.ID())

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	chatsResponse := make([]ChatResponse, len(chatsDomain))
	for i, chatDomain := range chatsDomain {
		chatsResponse[i] = toResponseChat(chatDomain)
	}
	c.JSON(http.StatusCreated, chatsResponse)
}

func (h HTTPServer) GetChatsByUser(c *gin.Context) {
	userCtx, err := utils.GetDataFromContext[domain.User](c, "user")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": domain.ErrNoUserInContext.Error()})
		return
	}

	chatsDomain, err := h.chatService.GetChatsByUser(c, userCtx.ID())

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	chatsResponse := make([]ChatResponse, len(chatsDomain))
	for i, chatDomain := range chatsDomain {
		chatsResponse[i] = toResponseChat(chatDomain)
	}
	c.JSON(http.StatusOK, chatsResponse)
}

func (h HTTPServer) EnterToChat(c *gin.Context) {
	userCtx, err := utils.GetDataFromContext[domain.User](c, "user")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": domain.ErrNoUserInContext.Error()})
		return
	}
	pubsub, err := uuid.Parse(c.Param("pubsub"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	chatID, err := strconv.Atoi(c.Query("chat_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	isChatMember := h.chatService.IsChatMember(c, userCtx.ID(), chatID)
	if !isChatMember {
		ok := h.stateService.SetCurrentChat(c, userCtx.ID(), pubsub, 0)
		if !ok {
			log.Print("SetCurrentChat is fault")
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": "not member chat " + strconv.Itoa(chatID)})
		return
	}

	ok := h.stateService.SetCurrentChat(c, userCtx.ID(), pubsub, chatID)
	if !ok {
		log.Print("SetCurrentChat is fault")
	}

	users, err := h.chatService.GetUsersByChatID(c, chatID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	response := make([]UserResponse, 0, len(users))
	for _, user := range users {
		_, ok := h.stateService.GetState(c, user.ID())
		status := "offline"
		if ok {
			status = "online"
		}
		response = append(response, toResponseUser(user, status))
	}

	c.JSON(http.StatusOK, response)
}
