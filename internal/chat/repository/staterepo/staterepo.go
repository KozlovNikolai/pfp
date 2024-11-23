package staterepo

import "github.com/KozlovNikolai/pfp/internal/chat/repository/models"

type UserID uint64

type stateStore struct {
	states map[UserID]models.State
}

func NewStateDB() *stateStore {
	return &stateStore{
		states: make(map[UserID]models.State),
	}
}
