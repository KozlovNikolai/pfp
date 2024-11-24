package domain

import "time"

type Message struct {
	id        int
	senderID  int
	chatID    int
	msgType   string
	text      string
	isDeleted bool
	createdAt int64
	updatedAt int64
}

type NewMessageData struct {
	Id        int
	SenderID  int
	ChatID    int
	MsgType   string
	Text      string
	IsDeleted bool
	CreatedAt int64
	UpdatedAt int64
}

func NewMessage(data NewMessageData) Message {
	createdTime := time.Now().Unix()
	return Message{
		id:        data.Id,
		senderID:  data.SenderID,
		chatID:    data.ChatID,
		msgType:   data.MsgType,
		text:      data.Text,
		isDeleted: false,
		createdAt: createdTime,
		updatedAt: createdTime,
	}
}

func (m Message) ID() int {
	return m.id
}
func (m Message) SenderID() int {
	return m.senderID
}
func (m Message) ChatID() int {
	return m.chatID
}
func (m Message) MsgType() string {
	return m.msgType
}
func (m Message) Text() string {
	return m.text
}
func (m Message) IsDeleted() bool {
	return m.isDeleted
}
func (m Message) CreatedAt() int64 {
	return m.createdAt
}
func (m Message) UpdatedAt() int64 {
	return m.updatedAt
}
