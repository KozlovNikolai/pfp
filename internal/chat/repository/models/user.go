// Package models ...
package models

type User struct {
	ID        int
	UserExtID int
	Login     string
	Password  string
	Profile   string
	Name      string
	Surname   string
	Email     string
	UserType  string
	CreatedAt int64
	UpdatedAt int64
}

type Contact struct {
	Id        int
	AccountID int
	UserID    int
	Name      string
	Surname   string
	Phone     string
	Email     string
}
