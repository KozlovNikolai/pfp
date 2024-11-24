package services

import (
	"context"

	"github.com/KozlovNikolai/pfp/internal/chat/domain"
)

type MessageService struct {
	repo IMessageRepository
}

// NewUserService creates a new User service
func NewMessageService(repo IMessageRepository) MessageService {
	return MessageService{
		repo: repo,
	}
}

func (m MessageService) SaveMsg(ctx context.Context, msg domain.Message) error {
	return m.repo.SaveMsg(ctx, msg)
}

func (m MessageService) GetMessagesByChatID(ctx context.Context, chatID, limit, offset int) ([]domain.Message, error) {
	return m.repo.GetMessagesByChatID(ctx, chatID, limit, offset)
}
