package services

import (
	"context"

	"github.com/KozlovNikolai/pfp/internal/chat/domain"
)

// UserService is a User service
type ChatService struct {
	repo IChatRepository
}

// NewUserService creates a new User service
func NewChatService(repo IChatRepository) ChatService {
	return ChatService{
		repo: repo,
	}
}

// CreateChat implements services.IChatRepository.
func (c ChatService) CreateChat(ctx context.Context, chat domain.Chat) (domain.Chat, error) {
	return c.repo.CreateChat(ctx, chat)
}
func (c ChatService) AddUserToChat(ctx context.Context, userID int, chatID int) error {
	return c.repo.AddUserToChat(ctx, userID, chatID)
}
func (c ChatService) GetChatByNameAndType(ctx context.Context, name string, chatType string) (domain.Chat, error) {
	return c.repo.GetChatByNameAndType(ctx, name, chatType)
}
func (c ChatService) GetChatsByUser(ctx context.Context, userID int) ([]domain.Chat, error) {
	return c.repo.GetChatsByUser(ctx, userID)
}
func (c ChatService) GetUsersByChatID(ctx context.Context, chatID int) ([]int, error) {
	return c.repo.GetUsersByChatID(ctx, chatID)
}
