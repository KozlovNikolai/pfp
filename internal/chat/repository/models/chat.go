package models

type Chat struct {
	Id            int    `json:"id" db:"id"`
	Name          string `json:"name" db:"name"`
	AccountID     int    `json:"account_id" db:"account_id"`
	ChatType      string `json:"chat_type" db:"chat_type"`
	LastChatMsgID uint64 `json:"last_message_id" db:"last_message_id"`
	CreatedAt     int64  `json:"created_at" db:"created_at"`
	UpdatedAt     int64  `json:"updated_at" db:"updated_at"`
}
