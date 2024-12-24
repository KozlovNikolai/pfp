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
	"github.com/jackc/pgx/v4"
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
		if err != nil && err.Error() != "tx is closed" {
			log.Printf("create chat rollback error:%v", err)
		}
	}()
	// Вставка данных о чате и получение ID
	err = tx.QueryRow(
		ctx,
		`
			INSERT INTO chats (account_id, name,  chat_type, last_message_id, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5, $6)
			RETURNING id,account_id, name, chat_type, last_message_id, created_at, updated_at`,
		dbChat.AccountID,
		dbChat.Name,
		dbChat.ChatType,
		dbChat.LastChatMsgID,
		dbChat.CreatedAt,
		dbChat.UpdatedAt,
	).
		Scan(
			&insertedChat.Id,
			&insertedChat.AccountID,
			&insertedChat.Name,
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

func (c *ChatRepo) AddUserToChat(ctx context.Context, userID int, chatID int, role string) error {
	createdAt := time.Now().Unix()
	lastReadMsgID := 0
	notifications := true
	var insertedRecordID int
	// Начинаем транзакцию
	tx, err := c.db.WR.Begin(ctx)
	if err != nil {
		return fmt.Errorf(constants.FailedToBeginTransaction, err)
	}
	defer func() {
		err := tx.Rollback(ctx)
		if err != nil && err.Error() != "tx is closed" {
			log.Printf("add user to chat rollback error:%v", err)
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
		return fmt.Errorf("failed add to Chat new User: %w", err)
	}
	// Фиксация транзакции
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf(constants.FailedToBeginTransaction, err)
	}

	return nil
}

func (c *ChatRepo) GetChatsByUser(ctx context.Context, userID int) ([]domain.Chat, error) {
	// Начинаем транзакцию
	tx, err := c.db.RO.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf(constants.FailedToBeginTransaction, err)
	}
	defer func() {
		err := tx.Rollback(ctx)
		if err != nil && err.Error() != "tx is closed" {
			log.Printf("get chats by user rollback error:%v", err)
		}
	}()

	// SQL-запрос на получение чатов по userID
	query := `
		SELECT c.id,  c.name, c.account_id, c.chat_type, c.last_message_id, c.created_at, c.updated_at
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
			&chat.AccountID,
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

// func (c *ChatRepo) GetChatByNameAndType(ctx context.Context, name string, chatType string) (domain.Chat, error) {
// 	// Начинаем транзакцию
// 	tx, err := c.db.RO.Begin(ctx)
// 	if err != nil {
// 		return domain.Chat{}, fmt.Errorf(constants.FailedToBeginTransaction, err)
// 	}
// 	defer func() {
// 		err := tx.Rollback(ctx)
// 		if err != nil && err.Error() != "tx is closed" {
// 			log.Printf("get chats by name and type rollback error:%v", err)
// 		}
// 	}()

//		// SQL-запрос на получение чата по имени и типу
//		query := `
//			SELECT *
//			FROM chats
//			WHERE (name=$1 AND chat_type=$2)
//		`
//		// Выполняем запрос и сканируем результат в структуру
//		var chat models.Chat
//		// Выполняем запрос и сканируем результат в структуру User
//		err = c.db.RO.QueryRow(ctx, query, name, chatType).Scan(
//			&chat.Id,
//			&chat.AccountID,
//			&chat.Name,
//			&chat.ChatType,
//			&chat.LastChatMsgID,
//			&chat.CreatedAt,
//			&chat.UpdatedAt,
//		)
//		if err != nil {
//			return domain.Chat{}, fmt.Errorf("failed get chat by name: %s and type: %s: %w", name, chatType, err)
//		}
//		// Фиксация транзакции
//		if err := tx.Commit(ctx); err != nil {
//			return domain.Chat{}, fmt.Errorf(constants.FailedToBeginTransaction, err)
//		}
//		// fmt.Printf("insertedChat: %v\n", insertedChat)
//		// domainChat := chatToDomain(insertedChat)
//		// fmt.Printf("domainChat: %v\n", domainChat)
//		return chatToDomain(chat), nil
//	}
func (c *ChatRepo) GetUserIDsByChatID(ctx context.Context, chatID int) ([]int, error) {
	// Начинаем транзакцию
	tx, err := c.db.RO.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf(constants.FailedToBeginTransaction, err)
	}
	defer func() {
		err := tx.Rollback(ctx)
		if err != nil && err.Error() != "tx is closed" {
			log.Printf("get user ids by chat id rollback error:%v", err)
		}
	}()

	// SQL-запрос на получение чатов по userID
	query := `
	SELECT user_id
	FROM chat_members
	WHERE chat_id=$1
`
	rows, err := c.db.RO.Query(ctx, query, chatID)
	if err != nil {
		return nil, fmt.Errorf("failed get users by chat_id: %d: %w", chatID, err)
	}
	// log.Printf("SQL-запрос на получение чатов по userID: %d", userID)
	// Фиксация транзакции
	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf(constants.FailedToBeginTransaction, err)
	}
	// log.Println("Заполняем массив")
	// заполняем массив
	var usersID []int
	for rows.Next() {
		var userID int
		if err := rows.Scan(
			&userID,
		); err != nil {
			return nil, fmt.Errorf("error scanning user ID: %w", err)
		}
		usersID = append(usersID, userID)
	}
	// log.Println("проверка на ошибки итерации")
	// Проверка на ошибки, возникшие при итерации по строкам
	if rows.Err() != nil {
		return nil, fmt.Errorf("error occurred during row iteration: %w", rows.Err())
	}
	// log.Println("мапим модель в домен")

	// fmt.Printf("insertedChat: %v\n", insertedChat)
	// domainChat := chatToDomain(insertedChat)
	// fmt.Printf("domainChat: %v\n", domainChat)
	return usersID, nil
}

func (c *ChatRepo) GetChatMember(ctx context.Context, userID int, chatID int) (domain.ChatMember, bool) {
	// Начинаем транзакцию
	tx, err := c.db.RO.Begin(ctx)
	if err != nil {
		log.Print(constants.FailedToBeginTransaction, err.Error())
		return domain.ChatMember{}, false
	}
	defer func() {
		err := tx.Rollback(ctx)
		if err != nil && err.Error() != "tx is closed" {
			log.Printf("is chat member rollback error:%s", err.Error())
		}
	}()

	// SQL-запрос на получение чата по имени и типу
	query := `
		SELECT *
		FROM chat_members
		WHERE (user_id=$1 AND chat_id=$2)
	`
	// Выполняем запрос и сканируем результат в структуру
	var chatMember models.ChatMember
	// Выполняем запрос и сканируем результат в структуру User
	err = c.db.RO.QueryRow(ctx, query, userID, chatID).Scan(
		&chatMember.Id,
		&chatMember.ChatID,
		&chatMember.UserID,
		&chatMember.Role,
		&chatMember.LastReadMsgID,
		&chatMember.Notifications,
		&chatMember.CreatedAt,
		&chatMember.UpdatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			log.Printf("not found user:%d in chat:%d", userID, chatID)
			return domain.ChatMember{}, false
		}
		log.Printf("the search for a user:%d in the chat:%d failed: %s", userID, chatID, err.Error())
		return domain.ChatMember{}, false
	}

	// Фиксация транзакции
	if err := tx.Commit(ctx); err != nil {
		log.Print(constants.FailedToBeginTransaction, err.Error())
		return domain.ChatMember{}, false
	}

	return chatMemberToDomain(chatMember), true
}

func (c *ChatRepo) GetUsersByChatID(ctx context.Context, chatID int) ([]domain.User, error) {
	// Начинаем транзакцию
	tx, err := c.db.RO.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf(constants.FailedToBeginTransaction, err)
	}
	defer func() {
		err := tx.Rollback(ctx)
		if err != nil && err.Error() != "tx is closed" {
			log.Printf("get users by chat id rollback error:%v", err)
		}
	}()

	// SQL-запрос на получение чатов по userID
	query := `
		SELECT u.id,  u.user_ext_id, u.login, u.password, u.profile, u.name, u.surname, u.email, u.user_type, u.created_at, u.updated_at
		FROM chat_members cm
		JOIN users u ON cm.user_id = u.id
		WHERE (cm.chat_id=$1)
	`
	rows, err := c.db.RO.Query(ctx, query, chatID)
	if err != nil {
		return nil, fmt.Errorf("failed get users by chat_id: %d: %w", chatID, err)
	}
	// Фиксация транзакции
	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf(constants.FailedToBeginTransaction, err)
	}
	defer rows.Close()
	// Заполняем массив пользователей
	var users []models.User
	for rows.Next() {
		var user models.User
		err := rows.Scan(
			&user.ID,
			&user.UserExtID,
			&user.Login,
			&user.Password,
			&user.Profile,
			&user.Name,
			&user.Surname,
			&user.Email,
			&user.UserType,
			&user.CreatedAt,
			&user.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		users = append(users, user)
	}

	// Проверка на ошибки, возникшие при итерации по строкам
	if rows.Err() != nil {
		return nil, fmt.Errorf("error occurred during row iteration: %w", rows.Err())
	}
	// мапим модель в домен
	domainUsers := make([]domain.User, len(users))
	for i, user := range users {
		domainUser := userToDomain(user)
		domainUsers[i] = domainUser
	}
	return domainUsers, nil
}

func (c *ChatRepo) GetChatByNameAndType(ctx context.Context, name, chatType string) (domain.Chat, error) {
	var query string
	if name == "" && chatType != "" {
		name = "%"
		query = `
		SELECT *
		FROM chats
		WHERE (name LIKE $1 AND chat_type = $2)
	`
	} else if name != "" && chatType == "" {
		chatType = "%"
		query = `
		SELECT *
		FROM chats
		WHERE (name = $1 AND chat_type LIKE $2)
	`
	} else if name != "" && chatType != "" {
		query = `
		SELECT *
		FROM chats
		WHERE (name = $1 AND chat_type = $2)
	`
	} else {
		return domain.Chat{}, fmt.Errorf("%w: chat name and/or chat type", domain.ErrRequired)
	}

	var chat models.Chat
	// Выполняем запрос и сканируем результат в структуру Chat
	err := c.db.RO.QueryRow(ctx, query, name, chatType).Scan(
		&chat.Id,
		&chat.AccountID,
		&chat.Name,
		&chat.ChatType,
		&chat.LastChatMsgID,
		&chat.CreatedAt,
		&chat.UpdatedAt,
	)
	if err != nil {
		return domain.Chat{}, fmt.Errorf("failed to get Chat by name: %s, type: %s, error: %w", name, chatType, err)
	}

	domainChat := chatToDomain(chat)

	return domainChat, nil
}
