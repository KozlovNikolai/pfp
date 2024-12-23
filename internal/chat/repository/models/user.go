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
