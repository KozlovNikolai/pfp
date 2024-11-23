// Package models ...
package models

type User struct {
	ID        int
	UserExtID string
	Login     string
	Password  string
	Account   string
	Token     string
	Name      string
	Surname   string
	Email     string
	UserType  string
	CreatedAt int64
	UpdatedAt int64
}
