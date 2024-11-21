package staterepo

import "github.com/KozlovNikolai/pfp/internal/chat/repository/models"

type UserChatID uint64

type stateStore struct {
	states map[UserChatID]models.State
}

func NewStateDB() *stateStore {
	return &stateStore{
		states: make(map[UserChatID]models.State),
	}
}
