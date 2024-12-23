package domain

import (
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type State struct {
	UserID   int
	Connects []Connect `json:"connects"`
	// Contacts []Contact `json:"contacts" db:"contacts"`
	// Chats    []Chat    `json:"chats" db:"chats"`
}

type Connect struct {
	Conn         *websocket.Conn
	Pubsub       uuid.UUID     `json:"pubsub"`
	CreatedAt    int64         `json:"created_at"`
	CurrentChat  int           `json:"current_chat"`
	CanselCannel chan struct{} `json:"-"`
}

// type Contact struct {
// 	UserID  uint64 `json:"user_id" db:"user_id"`
// 	Status  string `json:"status" db:"status"`
// 	Event   string `json:"event" db:"event"`
// 	Name    string `json:"name" db:"name"`
// 	Surname string `json:"surname" db:"surname"`
// 	Email   string `json:"email" db:"email"`
// }

// func NewState() *State {
// 	return &State{
// 		Connects: []Connect{},
// 		// Contacts: []Contact{},
// 		// Chats:    []Chat{},
// 	}
// }
