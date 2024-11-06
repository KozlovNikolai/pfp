package ws

import (
	"fmt"
	"net/http"

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
	ID   string `json:"id"`
	Name string `json:"name"`
}

// CreateRoom ...
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
	conn, err := upgrade.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	roomID := c.Param("roomID")
	clientID := c.Query("userID")
	username := c.Query("username")
	fmt.Printf("roomID: %s, clientID: %s, username: %s\n", roomID, clientID, username)

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

// GetRooms ...
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

// GetClients ...
func (h *Handler) GetClients(c *gin.Context) {
	clients := make([]ClientRes, 0)
	roomID := c.Param("roomID")

	if _, ok := h.hub.Rooms[roomID]; !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "room not found"})
		return
	}
	for _, client := range h.hub.Rooms[roomID].Clients {
		clients = append(clients, ClientRes{
			ID:       client.ID,
			Username: client.Username,
		})
	}

	c.JSON(http.StatusOK, clients)
}
