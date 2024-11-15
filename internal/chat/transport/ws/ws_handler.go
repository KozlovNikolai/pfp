package ws

import (
	"fmt"
	"net/http"
	"strconv"

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

// CreateRoomReq ...
type CreateRoomReq struct {
	ID   string `json:"id" example:"1"`
	Name string `json:"name" example:"Room1"`
}

// CreateRoom is ...
// CreateRoomTags		godoc
// @Summary				Создать комнату.
// @Description			Create new room in the system.
// @Param				CreateRoomReq body CreateRoomReq true "Create room."
// @Produce				application/json
// @Tags				Room
// @Security			BearerAuth
// @Success				201 {object} CreateRoomReq
// @failure				400 {string} err.Error()
// @failure				500 {string} string "error-to-create-room"
// @Router				/auth/ws/createRoom [post]
func (h *Handler) CreateRoom(c *gin.Context) {
	var req CreateRoomReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	h.hub.Rooms[req.ID] = &Room{
		ID:      req.ID,
		Name:    req.Name,
		Clients: make(map[string]*Client),
	}

	c.JSON(http.StatusCreated, req)
}

var upgrade = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(_ *http.Request) bool {
		return true
	},
}

// JoinRoom ...
func (h *Handler) JoinRoom(c *gin.Context) {
	user, err := utils.GetUserFromContext(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	conn, err := upgrade.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	roomID := c.Param("roomID")
	clientID := strconv.Itoa(user.ID())
	username := user.Login()
	fmt.Printf("\nroomID: %s, clientID: %s, username: %s\n\n", roomID, clientID, username)

	client := &Client{
		Conn:     conn,
		Message:  make(chan *Message, 10),
		ID:       clientID,
		RoomID:   roomID,
		Username: username,
	}

	msg := &Message{
		Content:  "a new user joined the room",
		RoomID:   roomID,
		Username: username,
	}

	h.hub.Register <- client
	h.hub.Broadcast <- msg

	go client.writeMessage()
	client.readMessage(h.hub)
}

// RoomRes ...
type RoomRes struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// GetRooms is ...
// GetRoomsTags 		godoc
// @Summary			Получить список всех комнат.
// @Description		Return rooms list.
// @Tags			Room
// @Security		BearerAuth
// @Produce      	json
// @Success			200 {object} []RoomRes
// @failure			404 {string} err.Error()
// @Router			/auth/ws/getRooms [get]
func (h *Handler) GetRooms(c *gin.Context) {
	rooms := make([]RoomRes, 0)
	for _, room := range h.hub.Rooms {
		rooms = append(rooms, RoomRes{
			ID:   room.ID,
			Name: room.Name,
		})
	}
	fmt.Printf("rooms: %v\n", rooms)
	c.JSON(http.StatusOK, rooms)
}

// ClientRes ...
type ClientRes struct {
	ID       string `json:"id"`
	Username string `json:"username"`
}

// GetClients is ...
// GetClientsTags 		godoc
// @Summary			Получить список всех участников группы.
// @Description		Return room clients list.
// @Tags			Room
// @Security		BearerAuth
// @Param			roomID path int true "Room ID" example(1) default(1)
// @Produce      	json
// @Success			200 {object} []ClientRes
// @failure			404 {string} err.Error()
// @Router			/auth/ws/getClients/{roomID} [get]
func (h *Handler) GetClients(c *gin.Context) {
	clients := make([]ClientRes, 0)
	roomID := c.Param("roomID")

	if _, ok := h.hub.Rooms[roomID]; !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "room not found"})
		return
	}

	for _, client := range h.hub.Rooms[roomID].Clients {
		fmt.Printf("clients: %+v\n", client)
		clients = append(clients, ClientRes{
			ID:       client.ID,
			Username: client.Username,
		})
	}

	c.JSON(http.StatusOK, clients)
}
