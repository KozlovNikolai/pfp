package services

import (
	"context"

	"alfachat/internal/chat/domain"

	"github.com/google/uuid"
)

func (s TokenService) GetPubsubToken(ctx context.Context, user domain.User) (uuid.UUID, error) {
	_ = ctx
	pubsub := uuid.New()
	return pubsub, nil
}
