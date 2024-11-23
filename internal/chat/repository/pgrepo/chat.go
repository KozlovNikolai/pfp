package pgrepo

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/KozlovNikolai/pfp/internal/chat/constants"
	"github.com/KozlovNikolai/pfp/internal/chat/domain"
	"github.com/KozlovNikolai/pfp/internal/chat/repository/models"
	"github.com/KozlovNikolai/pfp/internal/pkg/pg"
)

// UserRepo ...
type ChatRepo struct {
	db *pg.DB
}

// NewUserRepo ...
func NewChatRepo(db *pg.DB) *ChatRepo {
	return &ChatRepo{
		db: db,
	}
}

// CreateUserChat implements services.IUserRepository.
func (c *ChatRepo) CreateChat(ctx context.Context, chat domain.Chat) (domain.Chat, error) {
	dbChat := domainToChat(chat)

	dbChat.CreatedAt = time.Now().Unix()
	dbChat.UpdatedAt = dbChat.CreatedAt
	fmt.Printf("\ndbChat: %+v\n\n", dbChat)
	var insertedChat models.Chat

	// Начинаем транзакцию
	tx, err := c.db.WR.Begin(ctx)
	if err != nil {
		return domain.Chat{}, fmt.Errorf(constants.FailedToBeginTransaction, err)
	}
	defer func() {
		err := tx.Rollback(ctx)
		if err != nil {
			log.Printf("error:%v", err)
		}
	}()
	// Вставка данных о чате и получение ID
	err = tx.QueryRow(
		ctx,
		`
			INSERT INTO chats (name, owner_id, chat_type, last_message_id, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5, $6)
			RETURNING id,name, owner_id, chat_type, last_message_id, created_at, updated_at`,
		dbChat.Name,
		dbChat.OwnerID,
		dbChat.ChatType,
		dbChat.LastChatMsgID,
		dbChat.CreatedAt,
		dbChat.UpdatedAt,
	).
		Scan(
			&insertedChat.Id,
			&insertedChat.Name,
			&insertedChat.OwnerID,
			&insertedChat.ChatType,
			&insertedChat.LastChatMsgID,
			&insertedChat.CreatedAt,
			&insertedChat.UpdatedAt)
	if err != nil {
		return domain.Chat{}, fmt.Errorf("failed to insert new Chat: %w", err)
	}
	// Фиксация транзакции
	if err := tx.Commit(ctx); err != nil {
		return domain.Chat{}, fmt.Errorf(constants.FailedToBeginTransaction, err)
	}
	fmt.Printf("insertedChat: %v\n", insertedChat)
	domainChat := chatToDomain(insertedChat)
	fmt.Printf("domainChat: %v\n", domainChat)
	return domainChat, nil
}
