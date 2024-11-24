package ws

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/KozlovNikolai/pfp/internal/chat/domain"
	"github.com/KozlovNikolai/pfp/internal/pkg/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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

var upgrade = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(_ *http.Request) bool {
		return true
	},
}

func (h *Handler) JoinChat(c *gin.Context) {
	user, err := utils.GetDataFromContext[domain.User](c, "user")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	conn, err := upgrade.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	chatID, err := strconv.Atoi(c.Param("chatID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	clientID := user.ID()
	username := user.Login()
	fmt.Printf("\nchatID: %d, clientID: %d, username: %s\n\n", chatID, clientID, username)

	client := &Subscriber{
		Conn:    conn,
		Message: make(chan *Message, 10),
		ID:      clientID,
		// ChatID:   chatID,
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
	log.Println("subscribe is runing")
	user, err := utils.GetDataFromContext[domain.User](c, "user")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	log.Println("getting data from context - good")
	conn, err := upgrade.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	log.Println("connecting to websocket")

	pubsub, err := uuid.Parse(c.Param("pubsub"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	clientID := user.ID()
	login := user.Login()
	chatID, err := strconv.Atoi(c.Param("chatID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	name := user.Name()
	surname := user.Surname()
	account := user.Account()
	fmt.Printf("\nclientID: %d, login: %s, pubsub: %s\n", clientID, login, pubsub)
	fmt.Printf("name:%s, surname %s, account: %s\nIs connected!\n", name, surname, account)

	username := name + " " + surname
	subscriber := &Subscriber{
		Conn:    conn,
		Message: make(chan *Message, 10),
		ID:      clientID,
		// ChatID:   chatID,
		Username: username,
	}
	h.hub.stateService.SetState(c, user.ID(), pubsub, conn)
	state, ok := h.hub.stateService.GetState(c, user.ID())
	log.Printf("state: %+v\n", state)
	log.Printf("ok: %+v\n", ok)

	msg := &Message{
		Content:  "a new user joined the chat",
		ChatID:   chatID,
		Username: username,
	}

	h.hub.Register <- subscriber
	h.hub.Broadcast <- msg

	go subscriber.writeMessage()
	subscriber.readMessage(h.hub)
}

// ChatResponse ...
type ChatResponse struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type ClientResponse struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
}

func (h *Handler) GetClients(c *gin.Context) {
	clients := make([]ClientResponse, 0)

	for _, client := range h.hub.Node {
		fmt.Printf("clients: %+v\n", client)
		clients = append(clients, ClientResponse{
			ID:       client.ID,
			Username: client.Username,
		})
	}

	c.JSON(http.StatusOK, clients)
}
