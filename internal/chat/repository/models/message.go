package models

type Message struct {
	Id        int
	SenderID  int
	ChatID    int
	MsgType   string
	Text      string
	IsDeleted bool
	CreatedAt int64
	UpdatedAt int64
}
