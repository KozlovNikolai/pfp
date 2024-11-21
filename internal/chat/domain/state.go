package domain

import "github.com/google/uuid"

type State struct {
	Connects []Connect `json:"connects" db:"connects"`
	Contacts []Contact `json:"contacts" db:"contacts"`
	Chats    []Chat    `json:"chats" db:"chats"`
}

type Connect struct {
	WStoken   uuid.UUID `json:"ws_token" db:"ws_token"`
	CreatedAt uint64    `json:"created_at" db:"created_at"`
}

type Contact struct {
}

type Chat struct {
}
