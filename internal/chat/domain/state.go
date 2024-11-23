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
	UserID  uint64
	Status  string
	Event   string
	Name    string
	Surname string
	Email   string
}

func NewState() *State {
	return &State{
		Connects: []Connect{},
		Contacts: []Contact{},
		Chats:    []Chat{},
	}
}
