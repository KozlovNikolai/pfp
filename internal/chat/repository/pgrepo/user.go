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
type UserRepo struct {
	db *pg.DB
}

// NewUserRepo ...
func NewUserRepo(db *pg.DB) *UserRepo {
	return &UserRepo{
		db: db,
	}
}

// CreateUser implements services.IUserRepository.
func (u *UserRepo) CreateUser(ctx context.Context, user domain.User) (domain.User, error) {
	dbUser := domainToUser(user)

	dbUser.Login = dbUser.Email
	dbUser.CreatedAt = time.Now().Unix()
	dbUser.UpdatedAt = dbUser.CreatedAt
	dbUser.UserType = constants.User_type

	var insertedUser models.User

	// Начинаем транзакцию
	tx, err := u.db.WR.Begin(ctx)
	if err != nil {
		return domain.User{}, fmt.Errorf(constants.FailedToBeginTransaction, err)
	}
	defer func() {
		err := tx.Rollback(ctx)
		if err != nil {
			log.Printf("error:%v", err)
		}
	}()
	// Вставка данных о пользователе и получение ID
	err = tx.QueryRow(
		ctx,
		`
			INSERT INTO users (user_ext_id, login, password, account, token, name, surname, email, user_type, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
			RETURNING id,user_ext_id, login, password, account, token, name, surname, email, user_type, created_at, updated_at`,
		dbUser.UserExtID,
		dbUser.Login,
		dbUser.Password,
		dbUser.Account,
		dbUser.Token,
		dbUser.Name,
		dbUser.Surname,
		dbUser.Email,
		dbUser.UserType,
		dbUser.CreatedAt,
		dbUser.UpdatedAt,
	).
		Scan(
			&insertedUser.ID,
			&insertedUser.UserExtID,
			&insertedUser.Login,
			&insertedUser.Password,
			&insertedUser.Account,
			&insertedUser.Token,
			&insertedUser.Name,
			&insertedUser.Surname,
			&insertedUser.Email,
			&insertedUser.UserType,
			&insertedUser.CreatedAt,
			&insertedUser.UpdatedAt)
	if err != nil {
		return domain.User{}, fmt.Errorf("failed to insert User: %w", err)
	}
	// Фиксация транзакции
	if err := tx.Commit(ctx); err != nil {
		return domain.User{}, fmt.Errorf(constants.FailedToBeginTransaction, err)
	}
	// fmt.Printf("insertedUser: %v\n", insertedUser)
	domainUser := userToDomain(insertedUser)
	// fmt.Printf("domainUser: %v\n", domainUser)
	return domainUser, nil
}

// GetUserByExtID implements services.IUserRepository.
func (u *UserRepo) GetUserByExtID(ctx context.Context, account, extID string) (domain.User, error) {
	if extID == "" {
		return domain.User{}, fmt.Errorf("%w: ext ID", domain.ErrRequired)
	}
	if account == "" {
		return domain.User{}, fmt.Errorf("%w: account", domain.ErrRequired)
	}
	// SQL-запрос на получение данных Пользователя по extID
	query := `
		SELECT *
		FROM users
		WHERE (account=$1 AND user_ext_id=$2)
	`
	var user models.User
	// Выполняем запрос и сканируем результат в структуру User
	err := u.db.RO.QueryRow(ctx, query, account, extID).Scan(
		&user.ID,
		&user.UserExtID,
		&user.Login,
		&user.Password,
		&user.Account,
		&user.Token,
		&user.Name,
		&user.Surname,
		&user.Email,
		&user.UserType,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return domain.User{}, fmt.Errorf("failed to get User by extId: %s: %w", extID, err)
	}

	domainUser := userToDomain(user)

	return domainUser, nil
}

// DeleteUser implements service.IUserRepository.
func (u *UserRepo) DeleteUser(ctx context.Context, id int) error {
	if id == 0 {
		return fmt.Errorf("%w: id", domain.ErrRequired)
	}
	// Начинаем транзакцию
	tx, err := u.db.WR.Begin(ctx)
	if err != nil {
		return fmt.Errorf(constants.FailedToBeginTransaction, err)
	}
	defer func() {
		err := tx.Rollback(ctx)
		if err != nil {
			log.Printf("error:%v", err)
		}
	}()
	// Проверяем, что пользователь не связан ни с одним заказом.
	var count int
	err = tx.QueryRow(ctx, `
		SELECT COUNT(*)
		FROM orders
		WHERE user_id = (SELECT id FROM myusers WHERE id = $1)`, id).Scan(&count)
	if err != nil {
		return fmt.Errorf("failed to request the orders users: %w", err)
	}
	if count > 0 {
		return fmt.Errorf("error, there are order-related users.: %w", err)
	}
	// Удаляем пользователя
	_, err = tx.Exec(ctx, `
		DELETE FROM users
		WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("failed to delete User with id %d: %w", id, err)
	}
	// Фиксация транзакции
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}
	return nil
}

// GetUsers implements service.IUserRepository.
func (u *UserRepo) GetUsers(ctx context.Context, account string, limit, offset int) ([]domain.User, error) {
	query := `
		SELECT *
		FROM users
		WHERE account=$1
		ORDER BY id
		LIMIT $2 OFFSET $3
	`
	// Запрос
	rows, err := u.db.RO.Query(ctx, query, account, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
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
			&user.Account,
			&user.Token,
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

// GetUserByID implements service.IUserRepository.
func (u *UserRepo) GetUserByID(ctx context.Context, id int) (domain.User, error) {
	// SQL-запрос на получение данных Пользователя по ID
	query := `
		SELECT *
		FROM users
		WHERE id = $1
	`
	var user models.User
	// Выполняем запрос и сканируем результат в структуру User
	err := u.db.RO.QueryRow(ctx, query, id).Scan(
		&user.ID,
		&user.UserExtID,
		&user.Login,
		&user.Password,
		&user.Account,
		&user.Token,
		&user.Name,
		&user.Surname,
		&user.Email,
		&user.UserType,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return domain.User{}, fmt.Errorf("failed to get User by id: %w", err)
	}

	domainUser := userToDomain(user)
	return domainUser, nil
}

// GetUserByLogin implements service.IUserRepository.
func (u *UserRepo) GetUserByLogin(ctx context.Context, account, login string) (domain.User, error) {
	if login == "" {
		return domain.User{}, fmt.Errorf("%w: login", domain.ErrRequired)
	}

	// SQL-запрос на получение данных Пользователя по логину
	query := `
		SELECT *
		FROM users
		WHERE (account = $1 AND login = $2)
	`
	var user models.User
	// Выполняем запрос и сканируем результат в структуру User
	err := u.db.RO.QueryRow(ctx, query, account, login).Scan(
		&user.ID,
		&user.UserExtID,
		&user.Login,
		&user.Password,
		&user.Account,
		&user.Token,
		&user.Name,
		&user.Surname,
		&user.Email,
		&user.UserType,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return domain.User{}, fmt.Errorf("failed to get User by account: %s, login: %s, error: %w", account, login, err)
	}

	domainUser := userToDomain(user)

	return domainUser, nil
}

// UpdateUser implements service.IUserRepository.
func (u *UserRepo) UpdateUser(ctx context.Context, user domain.User) (domain.User, error) {
	dbUser := domainToUser(user)
	dbUser.UpdatedAt = time.Now().Unix()
	// Начинаем транзакцию
	tx, err := u.db.WR.Begin(ctx)
	if err != nil {
		return domain.User{}, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		err := tx.Rollback(ctx)
		if err != nil {
			log.Printf("error:%v", err)
		}
	}()
	// SQL-запрос на обновление данных Поставщика
	query := `
		UPDATE users
		SET user_ext_id = $1, login = $2, account = $3, token = $4, name = $5, surname = $6, email = $7, user_type = $8, updated_at = $9
		WHERE id = $10
		RETURNING id, login, password, role, token
	`
	var updatedUser models.User

	// Выполняем запрос и сканируем обновленный результат в структуру User
	err = tx.QueryRow(ctx, query,
		dbUser.UserExtID, dbUser.Login, dbUser.Account, dbUser.Token, dbUser.Name, dbUser.Surname, dbUser.Email, dbUser.UserType, dbUser.UpdatedAt).
		Scan(&updatedUser.ID,
			&updatedUser.UserExtID,
			&updatedUser.Login,
			&updatedUser.Password,
			&updatedUser.Account,
			&updatedUser.Token,
			&updatedUser.Name,
			&updatedUser.Surname,
			&updatedUser.Email,
			&updatedUser.UserType,
			&updatedUser.CreatedAt,
			&updatedUser.UpdatedAt,
		)
	if err != nil {
		return domain.User{}, fmt.Errorf("failed to update User: %w", err)
	}
	// Фиксация транзакции
	if err := tx.Commit(ctx); err != nil {
		return domain.User{}, fmt.Errorf("failed to commit transaction: %w", err)
	}
	domainUser := userToDomain(updatedUser)

	return domainUser, nil
}
