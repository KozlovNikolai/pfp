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
type AccountRepo struct {
	db *pg.DB
}

// NewUserRepo ...
func NewAccountRepo(db *pg.DB) *AccountRepo {
	return &AccountRepo{
		db: db,
	}
}

// CreateUser implements services.IUserRepository.
func (a *AccountRepo) CreateAccount(ctx context.Context, account domain.Account) (domain.Account, error) {
	dbAccount := domainToAccount(account)

	dbAccount.CreatedAt = time.Now().Unix()
	dbAccount.UpdatedAt = dbAccount.CreatedAt

	var insertedAccount models.Account

	// Начинаем транзакцию
	tx, err := a.db.WR.Begin(ctx)
	if err != nil {
		return domain.Account{}, fmt.Errorf(constants.FailedToBeginTransaction, err)
	}
	defer func() {
		err := tx.Rollback(ctx)
		if err != nil && err.Error() != "tx is closed" {
			log.Printf("create account rollback error:%v", err)
		}
	}()
	// Вставка данных о пользователе и получение ID
	err = tx.QueryRow(
		ctx,
		`
			INSERT INTO accounts ( name, created_at, updated_at)
			VALUES ($1, $2, $3)
			RETURNING id, name, created_at, updated_at`,
		dbAccount.Name,
		dbAccount.CreatedAt,
		dbAccount.UpdatedAt,
	).
		Scan(
			&insertedAccount.ID,
			&insertedAccount.Name,
			&insertedAccount.CreatedAt,
			&insertedAccount.UpdatedAt)
	if err != nil {
		return domain.Account{}, fmt.Errorf("failed to insert Account: %w", err)
	}
	// Фиксация транзакции
	if err := tx.Commit(ctx); err != nil {
		return domain.Account{}, fmt.Errorf(constants.FailedToBeginTransaction, err)
	}
	// fmt.Printf("insertedUser: %v\n", insertedUser)
	domainAccount := accountToDomain(insertedAccount)
	// fmt.Printf("domainUser: %v\n", domainUser)
	return domainAccount, nil
}

func (a *AccountRepo) AddUserToAccount(ctx context.Context, userID int, accountID int, inviterID int, role string) error {
	createdAt := time.Now().Unix()
	updatedAt := createdAt

	// Начинаем транзакцию
	tx, err := a.db.WR.Begin(ctx)
	if err != nil {
		return fmt.Errorf(constants.FailedToBeginTransaction, err)
	}
	defer func() {
		err := tx.Rollback(ctx)
		if err != nil && err.Error() != "tx is closed" {
			log.Printf("bind user to account rollback error:%v", err)
		}
	}()
	var id int
	// Вставка данных и получение ID
	err = tx.QueryRow(
		ctx,
		`
			INSERT INTO account_user (account_id, user_id, role, inviter_id, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5, $6)
			RETURNING id`,
		accountID,
		userID,
		role,
		inviterID,
		createdAt,
		updatedAt,
	).
		Scan(
			&id)
	if err != nil {
		return fmt.Errorf("failed to bind User to Account: %w", err)
	}
	// Фиксация транзакции
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf(constants.FailedToBeginTransaction, err)
	}

	return nil
}

func (a *AccountRepo) GetAccountByUserID(ctx context.Context, userID int) (int, error) {
	// SQL-запрос на получение данных Пользователя по ID
	query := `
		SELECT account_id
		FROM account_user
		WHERE (user_id = $1 AND role = $2)
	`
	var accID int
	// Выполняем запрос и сканируем результат в структуру User
	err := a.db.RO.QueryRow(ctx, query, userID, constants.AccountRoleOwner).Scan(
		&accID,
	)
	if err != nil {
		return 0, fmt.Errorf("failed to get Accaunt ID by User id Where role is OWNER: %w", err)
	}

	return accID, nil
}

func (a *AccountRepo) GetContactsByAccount(ctx context.Context, accID int) ([]int, error) {
	query := `
	SELECT user_id
	FROM account_user
	WHERE (account_id=$1 AND role=$2)
`
	// Запрос
	rows, err := a.db.RO.Query(ctx, query, accID, constants.AccountRoleContact)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}
	defer rows.Close()
	// Заполняем массив идентификаторов пользователей
	var userIDs []int
	for rows.Next() {
		var userID int
		err := rows.Scan(
			&userID)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		userIDs = append(userIDs, userID)
	}

	// Проверка на ошибки, возникшие при итерации по строкам
	if rows.Err() != nil {
		return nil, fmt.Errorf("error occurred during row iteration: %w", rows.Err())
	}
	return userIDs, nil
}

// // GetUserByExtID implements services.IUserRepository.
// func (u *UserRepo) GetUserByExtID(ctx context.Context, account, extID string) (domain.User, error) {
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
// 		&user.Profile,
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
// 		if err != nil && err.Error() != "tx is closed" {
// 			log.Printf("delete user rollback error:%v", err)
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
// 			&user.Profile,
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
// 		&user.Profile,
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
// func (u *UserRepo) GetUserByLogin(ctx context.Context, profile, login string) (domain.User, error) {
// 	if login == "" {
// 		return domain.User{}, fmt.Errorf("%w: login", domain.ErrRequired)
// 	}

// 	// SQL-запрос на получение данных Пользователя по логину
// 	query := `
// 		SELECT *
// 		FROM users
// 		WHERE (profile = $1 AND login = $2)
// 	`
// 	var user models.User
// 	// Выполняем запрос и сканируем результат в структуру User
// 	err := u.db.RO.QueryRow(ctx, query, profile, login).Scan(
// 		&user.ID,
// 		&user.UserExtID,
// 		&user.Login,
// 		&user.Password,
// 		&user.Profile,
// 		&user.Token,
// 		&user.Name,
// 		&user.Surname,
// 		&user.Email,
// 		&user.UserType,
// 		&user.CreatedAt,
// 		&user.UpdatedAt,
// 	)
// 	if err != nil {
// 		return domain.User{}, fmt.Errorf("failed to get User by account: %s, login: %s, error: %w", profile, login, err)
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
// 		if err != nil && err.Error() != "tx is closed" {
// 			log.Printf("update user rollback error:%v", err)
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
// 		dbUser.UserExtID, dbUser.Login, dbUser.Profile, dbUser.Token, dbUser.Name, dbUser.Surname, dbUser.Email, dbUser.UserType, dbUser.UpdatedAt).
// 		Scan(&updatedUser.ID,
// 			&updatedUser.UserExtID,
// 			&updatedUser.Login,
// 			&updatedUser.Password,
// 			&updatedUser.Profile,
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
