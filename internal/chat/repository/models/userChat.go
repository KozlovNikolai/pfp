// Package models ...
package models

type UserChat struct {
	ID        int
	UserExtID string
	Name      string
	Surname   string
	Email     string
	UserType  string
	CreatedAt uint64
	UpdatedAt uint64
}
