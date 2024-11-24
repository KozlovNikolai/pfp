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
	// fmt.Printf("\ndbChat: %+v\n\n", dbChat)
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
	// fmt.Printf("insertedChat: %v\n", insertedChat)
	domainChat := chatToDomain(insertedChat)
	// fmt.Printf("domainChat: %v\n", domainChat)
	return domainChat, nil
}

func (c *ChatRepo) AddUserToChat(ctx context.Context, userID int, chatID int) ([]domain.Chat, error) {
	createdAt := time.Now().Unix()
	role := ""
	lastReadMsgID := 0
	notifications := true
	var insertedRecordID int
	// Начинаем транзакцию
	tx, err := c.db.WR.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf(constants.FailedToBeginTransaction, err)
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
			INSERT INTO chat_members (chat_id, user_id, role, last_read_msg_id, notifications, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7)
			RETURNING id`,
		chatID,
		userID,
		role,
		lastReadMsgID,
		notifications,
		createdAt,
		createdAt,
	).
		Scan(
			&insertedRecordID)
	if err != nil {
		return nil, fmt.Errorf("failed add to Chat new User: %w", err)
	}
	// Фиксация транзакции
	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf(constants.FailedToBeginTransaction, err)
	}
	chats, err := c.GetChatsByUser(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed add user to chat: %w", err)
	}
	// fmt.Printf("insertedChat: %v\n", insertedChat)
	// domainChat := chatToDomain(insertedChat)
	// fmt.Printf("domainChat: %v\n", domainChat)
	return chats, nil
}

func (c *ChatRepo) GetChatsByUser(ctx context.Context, userID int) ([]domain.Chat, error) {
	// Начинаем транзакцию
	tx, err := c.db.WR.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf(constants.FailedToBeginTransaction, err)
	}
	defer func() {
		err := tx.Rollback(ctx)
		if err != nil {
			log.Printf("error:%v", err)
		}
	}()

	// SQL-запрос на получение чатов по userID
	query := `
		SELECT c.id,  c.name, c.owner_id, c.chat_type, c.last_message_id, c.created_at, c.updated_at
		FROM chat_members cm
		JOIN chats c ON cm.chat_id = c.id
		WHERE (cm.user_id=$1)
	`
	rows, err := c.db.RO.Query(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed get chats by user_id: %d: %w", userID, err)
	}
	// log.Printf("SQL-запрос на получение чатов по userID: %d", userID)
	// Фиксация транзакции
	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf(constants.FailedToBeginTransaction, err)
	}
	// log.Println("Заполняем массив")
	// заполняем массив
	var chats []models.Chat
	for rows.Next() {
		var chat models.Chat
		if err := rows.Scan(
			&chat.Id,
			&chat.Name,
			&chat.OwnerID,
			&chat.ChatType,
			&chat.LastChatMsgID,
			&chat.CreatedAt,
			&chat.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("error scanning user ID: %w", err)
		}
		chats = append(chats, chat)
	}
	// log.Println("проверка на ошибки итерации")
	// Проверка на ошибки, возникшие при итерации по строкам
	if rows.Err() != nil {
		return nil, fmt.Errorf("error occurred during row iteration: %w", rows.Err())
	}
	// log.Println("мапим модель в домен")

	// мапим модель в домен
	domainChats := make([]domain.Chat, len(chats))
	for i, chat := range chats {
		domainChat := chatToDomain(chat)
		domainChats[i] = domainChat
	}
	// fmt.Printf("insertedChat: %v\n", insertedChat)
	// domainChat := chatToDomain(insertedChat)
	// fmt.Printf("domainChat: %v\n", domainChat)
	return domainChats, nil
}

func (c *ChatRepo) GetChatByNameAndType(ctx context.Context, name string, chatType string) (domain.Chat, error) {
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

	// SQL-запрос на получение чата по имени и типу
	query := `
		SELECT *
		FROM chats
		WHERE (name=$1 AND chat_type=$2)
	`
	// Выполняем запрос и сканируем результат в структуру
	var chat models.Chat
	// Выполняем запрос и сканируем результат в структуру User
	err = c.db.RO.QueryRow(ctx, query, name, chatType).Scan(
		&chat.Id,
		&chat.OwnerID,
		&chat.Name,
		&chat.ChatType,
		&chat.LastChatMsgID,
		&chat.CreatedAt,
		&chat.UpdatedAt,
	)
	if err != nil {
		return domain.Chat{}, fmt.Errorf("failed get chat by name: %s and type: %s: %w", name, chatType, err)
	}
	// Фиксация транзакции
	if err := tx.Commit(ctx); err != nil {
		return domain.Chat{}, fmt.Errorf(constants.FailedToBeginTransaction, err)
	}
	// fmt.Printf("insertedChat: %v\n", insertedChat)
	// domainChat := chatToDomain(insertedChat)
	// fmt.Printf("domainChat: %v\n", domainChat)
	return chatToDomain(chat), nil
}

// func (c *ChatRepo) GetUsersFromChat(ctx context.Context, chatID int) ([]domain.User, error) {
// 	createdAt := time.Now().Unix()
// 	role := ""
// 	lastReadMsgID := 0
// 	notifications := true
// 	var insertedRecord int
// 	// Начинаем транзакцию
// 	tx, err := c.db.WR.Begin(ctx)
// 	if err != nil {
// 		return nil, fmt.Errorf(constants.FailedToBeginTransaction, err)
// 	}
// 	defer func() {
// 		err := tx.Rollback(ctx)
// 		if err != nil {
// 			log.Printf("error:%v", err)
// 		}
// 	}()
// 	// Вставка данных о чате и получение ID
// 	err = tx.QueryRow(
// 		ctx,
// 		`
// 			INSERT INTO chat_members (chat_id, user_id, role, last_read_msg_id, notifications, created_at, updated_at)
// 			VALUES ($1, $2, $3, $4, $5, $6, $7)
// 			RETURNING id`,
// 		chatID,
// 		userID,
// 		role,
// 		lastReadMsgID,
// 		notifications,
// 		createdAt,
// 		createdAt,
// 	).
// 		Scan(
// 			&insertedRecord)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed add to Chat new User: %w", err)
// 	}
// 	// Фиксация транзакции
// 	if err := tx.Commit(ctx); err != nil {
// 		return nil, fmt.Errorf(constants.FailedToBeginTransaction, err)
// 	}
// 	// fmt.Printf("insertedChat: %v\n", insertedChat)
// 	// domainChat := chatToDomain(insertedChat)
// 	// fmt.Printf("domainChat: %v\n", domainChat)
// 	return nil, nil
// }
