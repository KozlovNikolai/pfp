package domain

import "github.com/google/uuid"

type State struct {
	Connects []Connect `json:"connects" db:"connects"`
	Contacts []Contact `json:"contacts" db:"contacts"`
	Chats    []Chat    `json:"chats" db:"chats"`
}

type Connect struct {
	Pubsub    uuid.UUID `json:"pubsub" db:"pubsub"`
	CreatedAt uint64    `json:"created_at" db:"created_at"`
}

type Contact struct {
	UserID  uint64 `json:"user_id" db:"user_id"`
	Status  string `json:"status" db:"status"`
	Event   string `json:"event" db:"event"`
	Name    string `json:"name" db:"name"`
	Surname string `json:"surname" db:"surname"`
	Email   string `json:"email" db:"email"`
}

func NewState() *State {
	return &State{
		Connects: []Connect{},
		Contacts: []Contact{},
		Chats:    []Chat{},
	}
}
