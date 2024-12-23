package models

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
// 	UserID  uint64
// 	Status  string
// 	Event   string
// 	Name    string
// 	Surname string
// 	Email   string
// }

// type Chat struct {
// 	ChatID        uint64
// 	ChatType      string
// 	LastChatMsgID uint64
// 	Contacts      []Contact
// }
