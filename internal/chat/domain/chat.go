package domain

import "time"

type Chat struct {
	id            int
	name          string
	ownerID       int
	chatType      string
	lastChatMsgID uint64
	// contacts      []Contact
	createdAt int64
	updatedAt int64
}

type NewChatData struct {
	ID            int
	Name          string
	OwnerID       int
	ChatType      string
	LastChatMsgID uint64
	// Contacts      []Contact
}

func NewChat(data NewChatData) Chat {
	createdTime := time.Now().Unix()
	return Chat{
		id:            data.ID,
		name:          data.Name,
		ownerID:       data.OwnerID,
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

func (c Chat) OwnerID() int {
	return c.ownerID
}
func (c Chat) ChatType() string {
	return c.chatType
}
func (c Chat) LastMsgID() uint64 {
	return c.lastChatMsgID
}

//	func (c Chat) Contacts() []Contact {
//		return c.contacts
//	}
func (c Chat) CreatedAt() int64 {
	return c.createdAt
}
func (c Chat) UpdatedAt() int64 {
	return c.updatedAt
}
