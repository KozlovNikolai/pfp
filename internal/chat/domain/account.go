package domain

import "time"

// User is a domain User
type Account struct {
	id        int
	name      string
	createdAt int64
	updatedAt int64
}

// NewUserData is a domain User
type NewAccountData struct {
	ID        int    `json:"id" db:"id"`
	Name      string `json:"name" db:"name"`
	CreatedAt int64  `json:"created_at" db:"created_at"`
	UpdatedAt int64  `json:"updated_at" db:"updated_at"`
}

// NewUser ...
func NewAccount(data NewAccountData) Account {
	createdTime := time.Now().Unix()
	return Account{
		id:        data.ID,
		name:      data.Name,
		createdAt: createdTime,
		updatedAt: createdTime,
	}
}

// ID ...
func (a Account) ID() int {
	return a.id
}

// Password ...
func (a Account) Name() string {
	return a.name
}

// CreratedAt ...
func (a Account) CreatedAt() int64 {
	return a.createdAt
}

// UpdatedAt ...
func (a Account) UpdatedAt() int64 {
	return a.updatedAt
}
