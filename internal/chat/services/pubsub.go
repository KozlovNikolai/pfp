package services

import (
	"context"

	"github.com/KozlovNikolai/pfp/internal/chat/domain"
	"github.com/google/uuid"
)

func (s TokenService) GetPubsubToken(ctx context.Context, user domain.UserChat) (uuid.UUID, error) {
	_ = ctx
	pubsub := uuid.New()
	return pubsub, nil
}
