package models

import (
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type State struct {
	Connects []Connect `json:"connects" db:"connects"`
	// Contacts []Contact `json:"contacts" db:"contacts"`
	// Chats    []Chat    `json:"chats" db:"chats"`
}

type Connect struct {
	Conn      *websocket.Conn
	Pubsub    uuid.UUID `json:"pubsub" db:"pubsub"`
	CreatedAt int64     `json:"created_at" db:"created_at"`
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
