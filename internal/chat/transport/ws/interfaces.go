package ws

import (
	"context"

	"github.com/KozlovNikolai/pfp/internal/chat/domain"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type IStateService interface {
	SetState(ctx context.Context, userID int, pubsub uuid.UUID, conn *websocket.Conn) domain.State
	GetState(ctx context.Context, userID int) (domain.State, bool)
	DeleteConnFromState(ctx context.Context, userID int, pubsub uuid.UUID) (domain.State, bool)
}
