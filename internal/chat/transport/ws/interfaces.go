package ws

import (
	"context"

	"alfachat/internal/chat/domain"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type IStateService interface {
	SetState(ctx context.Context, userID int, pubsub uuid.UUID, conn *websocket.Conn) domain.State
	GetState(ctx context.Context, userID int) (domain.State, bool)
	GetStateByPubsub(ctx context.Context, pubsub uuid.UUID) (domain.User, domain.State, int, bool) //user,state, index of connect, ifExists
	DeleteConnFromState(ctx context.Context, userID int, pubsub uuid.UUID) (domain.State, bool)
	SetConnIntoState(ctx context.Context, userID int, pubsub uuid.UUID, conn *websocket.Conn, indexConn int) bool
	GetAllStates(ctx context.Context) []domain.State
}
