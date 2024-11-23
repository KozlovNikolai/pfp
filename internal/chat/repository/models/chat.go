package models

import "github.com/KozlovNikolai/pfp/internal/chat/domain"

type Chat struct {
	Id            int
	Name          string
	OwnerID       int
	ChatType      string
	LastChatMsgID uint64
	Contacts      []domain.Contact
	CreatedAt     int64
	UpdatedAt     int64
}
