package domain

import "time"

type Chat struct {
	id            int
	name          string
	accountID     int
	chatType      string
	lastChatMsgID uint64
	// contacts      []Contact
	createdAt int64
	updatedAt int64
}

type NewChatData struct {
	ID            int
	Name          string
	AccountID     int
	ChatType      string
	LastChatMsgID uint64
	// Contacts      []Contact
}

func NewChat(data NewChatData) Chat {
	createdTime := time.Now().Unix()
	return Chat{
		id:            data.ID,
		name:          data.Name,
		accountID:     data.AccountID,
		chatType:      data.ChatType,
		lastChatMsgID: data.LastChatMsgID,
		// contacts:      data.Contacts,
		createdAt: createdTime,
		updatedAt: createdTime,
	}
}

func (c Chat) ID() int {
	return c.id
}

func (c Chat) Name() string {
	return c.name
}

func (c Chat) AccountID() int {
	return c.accountID
}
func (c Chat) ChatType() string {
	return c.chatType
}
func (c Chat) LastMsgID() uint64 {
	return c.lastChatMsgID
}
func (c Chat) CreatedAt() int64 {
	return c.createdAt
}
func (c Chat) UpdatedAt() int64 {
	return c.updatedAt
}

type ChatMember struct {
	id            int
	chatID        int
	userID        int
	role          string
	lastReadMsgID int
	notifications bool
	createdAt     int64
	updatedAt     int64
}

type NewChatMemberData struct {
	Id            int
	ChatID        int
	UserID        int
	Role          string
	LastReadMsgID int
	Notifications bool
	CreatedAt     int64
	UpdatedAt     int64
}

func NewChatMember(data NewChatMemberData) ChatMember {
	return ChatMember{
		id:            data.Id,
		chatID:        data.ChatID,
		userID:        data.UserID,
		role:          data.Role,
		lastReadMsgID: data.LastReadMsgID,
		notifications: data.Notifications,
		createdAt:     data.CreatedAt,
		updatedAt:     data.UpdatedAt,
	}
}

func (cm ChatMember) ID() int {
	return cm.id
}
func (cm ChatMember) ChatID() int {
	return cm.chatID
}
func (cm ChatMember) UserID() int {
	return cm.userID
}
func (cm ChatMember) Role() string {
	return cm.role
}
func (cm ChatMember) LastReadMsgID() int {
	return cm.lastReadMsgID
}
func (cm ChatMember) Notifications() bool {
	return cm.notifications
}
func (cm ChatMember) CreatedAt() int64 {
	return cm.createdAt
}
func (cm ChatMember) UpdatedAt() int64 {
	return cm.updatedAt
}
