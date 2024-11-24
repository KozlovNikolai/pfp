package staterepo

import (
	"context"
	"sync"
	"time"

	"github.com/KozlovNikolai/pfp/internal/chat/domain"
	"github.com/KozlovNikolai/pfp/internal/chat/repository/models"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type StateRepo struct {
	db    *stateStore
	mutex sync.RWMutex
}

func NewStateRepo(db *stateStore) *StateRepo {
	return &StateRepo{
		db: db,
	}
}

// SetState implements services.IStateRepository.
func (s *StateRepo) SetState(ctx context.Context, userID int, pubsub uuid.UUID, conn *websocket.Conn) domain.State {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	state, ok := s.db.states[userID]
	if ok {
		state.Connects = append(state.Connects, models.Connect{
			Conn:      conn,
			Pubsub:    pubsub,
			CreatedAt: time.Now().Unix(),
		})
		s.db.states[userID] = state
		return stateToDomain(state)
	} else {
		state := models.State{
			Connects: []models.Connect{
				{
					Conn:      conn,
					Pubsub:    pubsub,
					CreatedAt: time.Now().Unix(),
				},
			},
		}
		s.db.states[userID] = state

		return stateToDomain(state)
	}
}

// GetState implements services.IStateRepository.
func (s *StateRepo) GetState(ctx context.Context, userID int) (domain.State, bool) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	state, ok := s.db.states[userID]
	if !ok {
		return domain.State{}, false
	}
	return stateToDomain(state), true
}

func (s *StateRepo) DeleteConnFromState(ctx context.Context, userID int, pubsub uuid.UUID) (domain.State, bool) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	state, ok := s.db.states[userID]
	if !ok {
		return domain.State{}, false
	}
	for i, connect := range state.Connects {
		if connect.Pubsub == pubsub {
			state.Connects[i] = state.Connects[len(state.Connects)-1]
			state.Connects = state.Connects[:len(state.Connects)-1]
			if len(state.Connects) == 0 {
				delete(s.db.states, userID)
			}
			return domain.State{}, true
		}
	}
	return stateToDomain(state), false
}
