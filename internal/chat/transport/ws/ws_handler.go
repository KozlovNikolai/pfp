package ws

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/KozlovNikolai/pfp/internal/chat/constants"
	"github.com/KozlovNikolai/pfp/internal/chat/domain"
	"github.com/KozlovNikolai/pfp/internal/pkg/utils"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

// Handler ...
type Handler struct {
	hub *Hub
}

// NewHandler ...
func NewHandler(hub *Hub) *Handler {
	return &Handler{
		hub: hub,
	}
}

// CreateChatRequest ...
type CreateChatRequest struct {
	ID   string `json:"id" example:"1"`
	Name string `json:"name" example:"Chat1"`
}

// CreateChat is ...
// CreateChatTags		godoc
// @Summary				Создать комнату.
// @Description			Create new chat in the system.
// @Param				CreateChatReq body CreateChatReq true "Create chat."
// @Produce				application/json
// @Tags				Chat
// @Security			BearerAuth
// @Success				201 {object} CreateChatReq
// @failure				400 {string} err.Error()
// @failure				500 {string} string "error-to-create-chat"
// @Router				/auth/ws/createChat [post]
// func (h *Handler) CreateChat(c *gin.Context) {
// 	var req CreateChatRequest
// 	if err := c.ShouldBindJSON(&req); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}

// 	h.hub.Chats[req.ID] = &Chat{
// 		ID:      req.ID,
// 		Name:    req.Name,
// 		Clients: make(map[string]*Subscriber),
// 	}

// 	c.JSON(http.StatusCreated, req)
// }

var upgrade = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(_ *http.Request) bool {
		return true
	},
}

// JoinChat ...
func (h *Handler) JoinChat(c *gin.Context) {
	user, err := utils.GetDataFromContext[domain.User](c, "user")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	conn, err := upgrade.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	chatID := c.Param("chatID")
	clientID := strconv.Itoa(user.ID())
	username := user.Login()
	fmt.Printf("\nchatID: %s, clientID: %s, username: %s\n\n", chatID, clientID, username)

	client := &Subscriber{
		Conn:     conn,
		Message:  make(chan *Message, 10),
		ID:       clientID,
		ChatID:   chatID,
		Username: username,
	}

	msg := &Message{
		Content:  "a new user joined the chat",
		ChatID:   chatID,
		Username: username,
	}

	h.hub.Register <- client
	h.hub.Broadcast <- msg

	go client.writeMessage()
	client.readMessage(h.hub)
}

// JoinChat ...
func (h *Handler) Subscribe(c *gin.Context) {
	user, err := utils.GetDataFromContext[domain.User](c, "user")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	conn, err := upgrade.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	pubsub := c.Param("pubsub")
	clientID := strconv.Itoa(user.ID())
	login := user.Login()
	name := user.Name()
	surname := user.Surname()
	account := user.Account()
	fmt.Printf("\nclientID: %s, login: %s, pubsub: %s\n", clientID, login, pubsub)
	fmt.Printf("%s %s, %s\nIs connected!\n", name, surname, account)

	username := name + " " + surname
	subscriber := &Subscriber{
		Conn:     conn,
		Message:  make(chan *Message, 10),
		ID:       clientID,
		ChatID:   constants.SystemChat,
		Username: username,
	}

	msg := &Message{
		Content:  "a new user joined the chat",
		ChatID:   constants.SystemChat,
		Username: username,
	}

	h.hub.Register <- subscriber
	h.hub.Broadcast <- msg

	go subscriber.writeMessage()
	subscriber.readMessage(h.hub)
}

// ChatResponse ...
type ChatResponse struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// GetChats is ...
// GetChatsTags 		godoc
// @Summary			Получить список всех комнат.
// @Description		Return chats list.
// @Tags			Chat
// @Security		BearerAuth
// @Produce      	json
// @Success			200 {object} []ChatRes
// @failure			404 {string} err.Error()
// @Router			/auth/ws/getChats [get]
func (h *Handler) GetChats(c *gin.Context) {
	chats := make([]ChatResponse, 0)
	for _, chat := range h.hub.Chats {
		chats = append(chats, ChatResponse{
			ID:   chat.ID,
			Name: chat.Name,
		})
	}
	fmt.Printf("chats: %v\n", chats)
	c.JSON(http.StatusOK, chats)
}

// ClientResponse ...
type ClientResponse struct {
	ID       string `json:"id"`
	Username string `json:"username"`
}

// GetClients is ...
// GetClientsTags 		godoc
// @Summary			Получить список всех участников группы.
// @Description		Return chat clients list.
// @Tags			Chat
// @Security		BearerAuth
// @Param			chatID path int true "Chat ID" example(1) default(1)
// @Produce      	json
// @Success			200 {object} []ClientRes
// @failure			404 {string} err.Error()
// @Router			/auth/ws/getClients/{chatID} [get]
func (h *Handler) GetClients(c *gin.Context) {
	clients := make([]ClientResponse, 0)
	chatID := c.Param("chatID")

	if _, ok := h.hub.Chats[chatID]; !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "chat not found"})
		return
	}

	for _, client := range h.hub.Chats[chatID].Clients {
		fmt.Printf("clients: %+v\n", client)
		clients = append(clients, ClientResponse{
			ID:       client.ID,
			Username: client.Username,
		})
	}

	c.JSON(http.StatusOK, clients)
}
