package pgrepo

import (
	"context"
	"fmt"
	"log"

	"github.com/KozlovNikolai/pfp/internal/chat/constants"
	"github.com/KozlovNikolai/pfp/internal/chat/domain"
	"github.com/KozlovNikolai/pfp/internal/chat/repository/models"
	"github.com/KozlovNikolai/pfp/internal/pkg/pg"
)

// UserRepo ...
type MsgRepo struct {
	db *pg.DB
}

// NewUserRepo ...
func NewMsgRepo(db *pg.DB) *MsgRepo {
	return &MsgRepo{
		db: db,
	}
}

// CreateUser implements services.IUserRepository.
func (m *MsgRepo) SaveMsg(ctx context.Context, msg domain.Message) error {
	dbMsg := domainToMessage(msg)
	// log.Printf("  domain Message: %+v", msg)
	var insertedMsg models.Message
	// log.Printf("model db Message: %+v", dbMsg)
	// Начинаем транзакцию
	tx, err := m.db.WR.Begin(ctx)
	if err != nil {
		return fmt.Errorf(constants.FailedToBeginTransaction, err)
	}
	defer func() {
		err := tx.Rollback(ctx)
		if err != nil && err.Error() != "tx is closed" {
			log.Printf("save message rollback error:%v", err)
		}
	}()
	// Сохранение сообщения
	err = tx.QueryRow(
		ctx,
		`
			INSERT INTO messages (sender_id, chat_id, msg_type, text, is_deleted, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7)
			RETURNING id,sender_id, chat_id, msg_type, text, is_deleted, created_at, updated_at`,
		dbMsg.SenderID,
		dbMsg.ChatID,
		dbMsg.MsgType,
		dbMsg.Text,
		dbMsg.IsDeleted,
		dbMsg.CreatedAt,
		dbMsg.UpdatedAt,
	).
		Scan(
			&insertedMsg.Id,
			&insertedMsg.SenderID,
			&insertedMsg.ChatID,
			&insertedMsg.MsgType,
			&insertedMsg.Text,
			&insertedMsg.IsDeleted,
			&insertedMsg.CreatedAt,
			&insertedMsg.UpdatedAt)
	if err != nil {
		return fmt.Errorf("failed to insert Message: %w", err)
	}

	//*******************************************************
	// Обновление номера последнего сообщения в чате
	query := `
UPDATE chats
SET last_message_id = $1
WHERE id = $2
RETURNING last_message_id
`
	var lastMsgId int

	// Выполняем запрос и сканируем обновленный результат в структуру User
	err = tx.QueryRow(ctx, query,
		insertedMsg.Id, insertedMsg.ChatID).
		Scan(&lastMsgId)
	if err != nil {
		return fmt.Errorf("failed to insert Message: %w", err)
	}
	//*******************************************************
	// Фиксация транзакции
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf(constants.FailedToBeginTransaction, err)
	}
	return nil
}

func (m *MsgRepo) GetMessagesByChatID(ctx context.Context, chatID, limit, offset int) ([]domain.Message, error) {
	// Начинаем транзакцию
	tx, err := m.db.RO.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf(constants.FailedToBeginTransaction, err)
	}
	defer func() {
		err := tx.Rollback(ctx)
		if err != nil && err.Error() != "tx is closed" {
			log.Printf("get message by chat id rollback error:%v", err)
		}
	}()

	// SQL-запрос на получение ообщений по userID
	query := `
		SELECT *
		FROM messages
		WHERE (chat_id=$1)
		ORDER BY id
		LIMIT $2 OFFSET $3
	`
	rows, err := m.db.RO.Query(ctx, query, chatID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed get messages by chat_id: %d: %w", chatID, err)
	}
	// log.Printf("SQL-запрос на получение сообщений по chatID: %d", chatID)
	// Фиксация транзакции
	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf(constants.FailedToBeginTransaction, err)
	}
	// log.Println("Заполняем массив")
	// заполняем массив
	var msgs []models.Message
	for rows.Next() {
		var msg models.Message
		if err := rows.Scan(
			&msg.Id,
			&msg.SenderID,
			&msg.ChatID,
			&msg.MsgType,
			&msg.Text,
			&msg.IsDeleted,
			&msg.CreatedAt,
			&msg.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("error scanning user ID: %w", err)
		}
		msgs = append(msgs, msg)
	}
	// log.Println("проверка на ошибки итерации")
	// Проверка на ошибки, возникшие при итерации по строкам
	if rows.Err() != nil {
		return nil, fmt.Errorf("error occurred during row iteration: %w", rows.Err())
	}
	// log.Println("мапим модель в домен")

	// мапим модель в домен
	domainMsgs := make([]domain.Message, len(msgs))
	for i, msg := range msgs {
		domainMsg := messageToDomain(msg)
		domainMsgs[i] = domainMsg
	}
	return domainMsgs, nil
}

func (m *MsgRepo) GetChatMessages(ctx context.Context, chatID, initialMsgID, before, after int) ([]domain.Message, error) {
	// Начинаем транзакцию
	tx, err := m.db.RO.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf(constants.FailedToBeginTransaction, err)
	}
	defer func() {
		err := tx.Rollback(ctx)
		if err != nil && err.Error() != "tx is closed" {
			log.Printf("get chat messages rollback error:%v", err)
		}
	}()

	query := `
		(SELECT * FROM messages
		WHERE chat_id = $1 AND is_deleted = FALSE AND id < $2
		ORDER BY id DESC
		LIMIT $3)		
		UNION ALL 
		(SELECT * FROM messages
		WHERE chat_id = $1 AND is_deleted = FALSE AND id >= $2
		ORDER BY id ASC
		LIMIT $4)
		ORDER BY id ASC
`

	rows, err := m.db.RO.Query(ctx, query, chatID, initialMsgID, before, after)
	if err != nil {
		return nil, fmt.Errorf("failed get messages by chat_id: %d: %w", chatID, err)
	}
	// log.Printf("SQL-запрос на получение сообщений по chatID: %d", chatID)
	// Фиксация транзакции
	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf(constants.FailedToBeginTransaction, err)
	}
	// log.Println("Заполняем массив")
	// заполняем массив
	var msgs []models.Message
	for rows.Next() {
		var msg models.Message
		if err := rows.Scan(
			&msg.Id,
			&msg.SenderID,
			&msg.ChatID,
			&msg.MsgType,
			&msg.Text,
			&msg.IsDeleted,
			&msg.CreatedAt,
			&msg.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("error scanning user ID: %w", err)
		}
		msgs = append(msgs, msg)
	}
	// log.Println("проверка на ошибки итерации")
	// Проверка на ошибки, возникшие при итерации по строкам
	if rows.Err() != nil {
		return nil, fmt.Errorf("error occurred during row iteration: %w", rows.Err())
	}
	// log.Println("мапим модель в домен")

	// мапим модель в домен
	domainMsgs := make([]domain.Message, len(msgs))
	for i, msg := range msgs {
		domainMsg := messageToDomain(msg)
		domainMsgs[i] = domainMsg
	}
	return domainMsgs, nil
}

// // GetUserByExtID implements services.IUserRepository.
// func (m *MsgRepo) GetUserByExtID(ctx context.Context, account, extID string) (domain.User, error) {
// 	if extID == "" {
// 		return domain.User{}, fmt.Errorf("%w: ext ID", domain.ErrRequired)
// 	}
// 	if account == "" {
// 		return domain.User{}, fmt.Errorf("%w: account", domain.ErrRequired)
// 	}
// 	// SQL-запрос на получение данных Пользователя по extID
// 	query := `
// 		SELECT *
// 		FROM users
// 		WHERE (account=$1 AND user_ext_id=$2)
// 	`
// 	var user models.User
// 	// Выполняем запрос и сканируем результат в структуру User
// 	err := u.db.RO.QueryRow(ctx, query, account, extID).Scan(
// 		&user.ID,
// 		&user.UserExtID,
// 		&user.Login,
// 		&user.Password,
// 		&user.Account,
// 		&user.Token,
// 		&user.Name,
// 		&user.Surname,
// 		&user.Email,
// 		&user.UserType,
// 		&user.CreatedAt,
// 		&user.UpdatedAt,
// 	)
// 	if err != nil {
// 		return domain.User{}, fmt.Errorf("failed to get User by extId: %s: %w", extID, err)
// 	}

// 	domainUser := userToDomain(user)

// 	return domainUser, nil
// }

// // DeleteUser implements service.IUserRepository.
// func (u *UserRepo) DeleteUser(ctx context.Context, id int) error {
// 	if id == 0 {
// 		return fmt.Errorf("%w: id", domain.ErrRequired)
// 	}
// 	// Начинаем транзакцию
// 	tx, err := u.db.WR.Begin(ctx)
// 	if err != nil {
// 		return fmt.Errorf(constants.FailedToBeginTransaction, err)
// 	}
// 	defer func() {
// 		err := tx.Rollback(ctx)
// 		if err != nil {
// 			log.Printf("error:%v", err)
// 		}
// 	}()
// 	// Проверяем, что пользователь не связан ни с одним заказом.
// 	var count int
// 	err = tx.QueryRow(ctx, `
// 		SELECT COUNT(*)
// 		FROM orders
// 		WHERE user_id = (SELECT id FROM myusers WHERE id = $1)`, id).Scan(&count)
// 	if err != nil {
// 		return fmt.Errorf("failed to request the orders users: %w", err)
// 	}
// 	if count > 0 {
// 		return fmt.Errorf("error, there are order-related users.: %w", err)
// 	}
// 	// Удаляем пользователя
// 	_, err = tx.Exec(ctx, `
// 		DELETE FROM users
// 		WHERE id = $1`, id)
// 	if err != nil {
// 		return fmt.Errorf("failed to delete User with id %d: %w", id, err)
// 	}
// 	// Фиксация транзакции
// 	if err := tx.Commit(ctx); err != nil {
// 		return fmt.Errorf("failed to commit transaction: %w", err)
// 	}
// 	return nil
// }

// // GetUsers implements service.IUserRepository.
// func (u *UserRepo) GetUsers(ctx context.Context, account string, limit, offset int) ([]domain.User, error) {
// 	query := `
// 		SELECT *
// 		FROM users
// 		WHERE account=$1
// 		ORDER BY id
// 		LIMIT $2 OFFSET $3
// 	`
// 	// Запрос
// 	rows, err := u.db.RO.Query(ctx, query, account, limit, offset)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to execute query: %w", err)
// 	}
// 	defer rows.Close()
// 	// Заполняем массив пользователей
// 	var users []models.User
// 	for rows.Next() {
// 		var user models.User
// 		err := rows.Scan(
// 			&user.ID,
// 			&user.UserExtID,
// 			&user.Login,
// 			&user.Password,
// 			&user.Account,
// 			&user.Token,
// 			&user.Name,
// 			&user.Surname,
// 			&user.Email,
// 			&user.UserType,
// 			&user.CreatedAt,
// 			&user.UpdatedAt)
// 		if err != nil {
// 			return nil, fmt.Errorf("failed to scan row: %w", err)
// 		}
// 		users = append(users, user)
// 	}

// 	// Проверка на ошибки, возникшие при итерации по строкам
// 	if rows.Err() != nil {
// 		return nil, fmt.Errorf("error occurred during row iteration: %w", rows.Err())
// 	}
// 	// мапим модель в домен
// 	domainUsers := make([]domain.User, len(users))
// 	for i, user := range users {
// 		domainUser := userToDomain(user)
// 		domainUsers[i] = domainUser
// 	}
// 	return domainUsers, nil
// }

// // GetUserByID implements service.IUserRepository.
// func (u *UserRepo) GetUserByID(ctx context.Context, id int) (domain.User, error) {
// 	// SQL-запрос на получение данных Пользователя по ID
// 	query := `
// 		SELECT *
// 		FROM users
// 		WHERE id = $1
// 	`
// 	var user models.User
// 	// Выполняем запрос и сканируем результат в структуру User
// 	err := u.db.RO.QueryRow(ctx, query, id).Scan(
// 		&user.ID,
// 		&user.UserExtID,
// 		&user.Login,
// 		&user.Password,
// 		&user.Account,
// 		&user.Token,
// 		&user.Name,
// 		&user.Surname,
// 		&user.Email,
// 		&user.UserType,
// 		&user.CreatedAt,
// 		&user.UpdatedAt,
// 	)
// 	if err != nil {
// 		return domain.User{}, fmt.Errorf("failed to get User by id: %w", err)
// 	}

// 	domainUser := userToDomain(user)
// 	return domainUser, nil
// }

// // GetUserByLogin implements service.IUserRepository.
// func (u *UserRepo) GetUserByLogin(ctx context.Context, account, login string) (domain.User, error) {
// 	if login == "" {
// 		return domain.User{}, fmt.Errorf("%w: login", domain.ErrRequired)
// 	}

// 	// SQL-запрос на получение данных Пользователя по логину
// 	query := `
// 		SELECT *
// 		FROM users
// 		WHERE (account = $1 AND login = $2)
// 	`
// 	var user models.User
// 	// Выполняем запрос и сканируем результат в структуру User
// 	err := u.db.RO.QueryRow(ctx, query, account, login).Scan(
// 		&user.ID,
// 		&user.UserExtID,
// 		&user.Login,
// 		&user.Password,
// 		&user.Account,
// 		&user.Token,
// 		&user.Name,
// 		&user.Surname,
// 		&user.Email,
// 		&user.UserType,
// 		&user.CreatedAt,
// 		&user.UpdatedAt,
// 	)
// 	if err != nil {
// 		return domain.User{}, fmt.Errorf("failed to get User by account: %s, login: %s, error: %w", account, login, err)
// 	}

// 	domainUser := userToDomain(user)

// 	return domainUser, nil
// }

// // UpdateUser implements service.IUserRepository.
// func (u *UserRepo) UpdateUser(ctx context.Context, user domain.User) (domain.User, error) {
// 	dbUser := domainToUser(user)
// 	dbUser.UpdatedAt = time.Now().Unix()
// 	// Начинаем транзакцию
// 	tx, err := u.db.WR.Begin(ctx)
// 	if err != nil {
// 		return domain.User{}, fmt.Errorf("failed to begin transaction: %w", err)
// 	}
// 	defer func() {
// 		err := tx.Rollback(ctx)
// 		if err != nil {
// 			log.Printf("error:%v", err)
// 		}
// 	}()
// 	// SQL-запрос на обновление данных Поставщика
// 	query := `
// 		UPDATE users
// 		SET user_ext_id = $1, login = $2, account = $3, token = $4, name = $5, surname = $6, email = $7, user_type = $8, updated_at = $9
// 		WHERE id = $10
// 		RETURNING id, login, password, role, token
// 	`
// 	var updatedUser models.User

// 	// Выполняем запрос и сканируем обновленный результат в структуру User
// 	err = tx.QueryRow(ctx, query,
// 		dbUser.UserExtID, dbUser.Login, dbUser.Account, dbUser.Token, dbUser.Name, dbUser.Surname, dbUser.Email, dbUser.UserType, dbUser.UpdatedAt).
// 		Scan(&updatedUser.ID,
// 			&updatedUser.UserExtID,
// 			&updatedUser.Login,
// 			&updatedUser.Password,
// 			&updatedUser.Account,
// 			&updatedUser.Token,
// 			&updatedUser.Name,
// 			&updatedUser.Surname,
// 			&updatedUser.Email,
// 			&updatedUser.UserType,
// 			&updatedUser.CreatedAt,
// 			&updatedUser.UpdatedAt,
// 		)
// 	if err != nil {
// 		return domain.User{}, fmt.Errorf("failed to update User: %w", err)
// 	}
// 	// Фиксация транзакции
// 	if err := tx.Commit(ctx); err != nil {
// 		return domain.User{}, fmt.Errorf("failed to commit transaction: %w", err)
// 	}
// 	domainUser := userToDomain(updatedUser)

// 	return domainUser, nil
// }
