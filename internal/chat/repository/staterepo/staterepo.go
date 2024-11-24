package staterepo

import "github.com/KozlovNikolai/pfp/internal/chat/repository/models"

type stateStore struct {
	states map[int]models.State
}

func NewStateDB() *stateStore {
	return &stateStore{
		states: make(map[int]models.State),
	}
}
