package staterepo

import (
	"context"
	"fmt"
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
			UserID: userID,
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

func (s *StateRepo) GetStateByPubsub(ctx context.Context, pubsub uuid.UUID) (int, domain.State, int, bool) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	for userID, state := range s.db.states {
		for index, connect := range state.Connects {
			if connect.Pubsub == pubsub {
				return userID, stateToDomain(state), index, true
			}
		}
	}
	return 0, domain.State{}, 0, false
}

func (s *StateRepo) SetConnIntoState(ctx context.Context, userID int, pubsub uuid.UUID, conn *websocket.Conn, indexConn int) bool {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	state, ok := s.db.states[userID]
	if !ok {
		return false
	}
	st := state
	st.Connects[indexConn].Conn = conn
	s.db.states[userID] = st
	return true
}

func (s *StateRepo) DeleteConnFromState(ctx context.Context, userID int, pubsub uuid.UUID) (domain.State, bool) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	state, ok := s.db.states[userID]
	if !ok {
		return domain.State{}, false
	}
	for i, connect := range state.Connects {
		fmt.Printf("i: %d, connect.Pubsub: %v, pubsub: %v\n", i, connect.Pubsub, pubsub)
		if connect.Pubsub == pubsub {
			fmt.Printf("len(state.Connects): %d\n", len(state.Connects))
			state.Connects[i], state.Connects[len(state.Connects)-1] = state.Connects[len(state.Connects)-1], state.Connects[i]

			state.Connects = state.Connects[:len(state.Connects)-1]
			s.db.states[userID] = state
			if len(state.Connects) == 0 {
				delete(s.db.states, userID)
				fmt.Printf("Length = 0\n")
				return domain.State{}, true
			}
			fmt.Printf("Length = %d\n", len(state.Connects))
			return stateToDomain(state), true
		}
	}
	return stateToDomain(state), false
}

func (s *StateRepo) GetAllStates(ctx context.Context) []domain.State {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	var states []domain.State
	for _, state := range s.db.states {
		states = append(states, stateToDomain(state))
	}
	return states
}
