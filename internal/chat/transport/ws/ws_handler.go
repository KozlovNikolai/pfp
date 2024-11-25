package ws

import (
	"fmt"
	"log"
	"net/http"

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

func (h *Handler) Broadcast(content string, chatID int, senderID int, chatMembers []int) {

	msg := &Message{
		Content:     content,
		ChatID:      chatID,
		Sender:      senderID,
		ChatMembers: chatMembers,
	}

	h.hub.Broadcast <- msg
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
	// chatID, err := strconv.Atoi(c.Param("chatID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	// chatID = chatid
	name := user.Name()
	surname := user.Surname()
	account := user.Account()
	fmt.Printf("\nclientID: %d, login: %s, pubsub: %s\n", clientID, login, pubsub)
	fmt.Printf("name:%s, surname %s, account: %s\nIs connected!\n", name, surname, account)

	username := name + " " + surname
	subscriber := &Subscriber{
		Conn:     conn,
		Message:  make(chan *MessageOne, 10),
		ID:       clientID,
		Pubsub:   pubsub,
		Username: username,
	}
	h.hub.stateService.SetState(c, user.ID(), pubsub, conn)
	state, ok := h.hub.stateService.GetState(c, user.ID())
	log.Printf("state: %+v\n", state)
	log.Printf("ok: %+v\n", ok)

	h.hub.Register <- subscriber

	go subscriber.writeMessage()
	subscriber.readMessage(h.hub)
}
