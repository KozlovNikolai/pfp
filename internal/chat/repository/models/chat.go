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

type ChatMember struct {
	Id            int    `json:"id" db:"id"`
	ChatID        int    `json:"chat_id" db:"chat_id"`
	UserID        int    `json:"user_id" db:"user_id"`
	Role          string `json:"role" db:"role"`
	LastReadMsgID int    `json:"last_read_msg_id" db:"last_read_msg_id"`
	Notifications bool   `json:"notifications" db:"notifications"`
	CreatedAt     int64  `json:"created_at" db:"created_at"`
	UpdatedAt     int64  `json:"updated_at" db:"updated_at"`
}
