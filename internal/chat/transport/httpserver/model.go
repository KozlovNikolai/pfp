package httpserver

import (
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

// UserRequest is ...
type UserRequest struct {
	Profile  string `json:"profile" db:"profile" example:"email" validate:"required"`
	Login    string `json:"login"    db:"login"    example:"cmd@cmd.ru" validate:"required,email"`
	Password string `json:"password" db:"password" example:"123456"     validate:"required,min=6,max=32"`
	Name     string `json:"name" db:"name" example:"Ivan"`
	Surname  string `json:"surname" db:"surname" example:"Ivanov"`
}

// Validate ...
func (u *UserRequest) Validate() error {
	validate := validator.New(validator.WithRequiredStructEnabled())
	return validate.Struct(u)
}

type UserResponse struct {
	ID        int    `json:"id" db:"id"`
	UserExtID int    `json:"user_ext_id" db:"user_ext_id"`
	Login     string `json:"login" db:"login"`
	// Password  string `json:"password" db:"password"`
	Profile   string `json:"profile" db:"profile"`
	Name      string `json:"name" db:"name"`
	Surname   string `json:"surname" db:"surname"`
	Email     string `json:"email" db:"email"`
	UserType  string `json:"type" db:"type"`
	CreatedAt int64  `json:"created_at" db:"created_at"`
	UpdatedAt int64  `json:"updated_at" db:"updated_at"`
	Status    string `json:"status"`
}

// UserRequest is ...
type AccountRequest struct {
	ID   int    `json:"id" db:"id" example:"1"`
	Name string `json:"name" db:"name" example:"MyAccount"`
}

// Validate ...
func (a *AccountRequest) Validate() error {
	validate := validator.New(validator.WithRequiredStructEnabled())
	return validate.Struct(a)
}

type AccountResponse struct {
	ID        int    `json:"id" db:"id"`
	Name      string `json:"name" db:"name"`
	CreatedAt int64  `json:"created_at" db:"created_at"`
	UpdatedAt int64  `json:"updated_at" db:"updated_at"`
}

// // UserRequest is ...
// type AccountCreateRequest struct {
// 	Name string `json:"name" db:"name" example:"MyAccount" validate:"required"`
// }

// // Validate ...
// func (a *AccountCreateRequest) Validate() error {
// 	validate := validator.New(validator.WithRequiredStructEnabled())
// 	return validate.Struct(a)
// }

// type AccountCreateResponse struct {
// 	ID   int    `json:"id" db:"id"`
// 	Name string `json:"name" db:"name"`
// }

type ChatCreateRequest struct {
	OwnerID  int
	Name     string `json:"name" db:"name" validate:"required"`
	ChatType string `json:"chat_type" db:"chat_type" validate:"required"`
}

func (c *ChatCreateRequest) Validate() error {
	validate := validator.New(validator.WithRequiredStructEnabled())
	return validate.Struct(c)
}

type ChatResponse struct {
	Id            int    `json:"id"`
	Name          string `json:"name"`
	AccountID     int    `json:"account_id"`
	ChatType      string `json:"chat_type"`
	LastChatMsgID uint64 `json:"last_message_id"`
	// Contacts      []domain.Contact `json:"contacts"`
	CreatedAt int64 `json:"created_at"`
	UpdatedAt int64 `json:"updated_at"`
}

type SendMessageRequest struct {
	SenderID  int    `json:"-"`
	ChatID    int    `json:"chat_id" db:"chat_id" validate:"required"`
	MsgType   string `json:"msg_type" db:"msg_type" validate:"required"`
	Text      string `json:"text" db:"text" validate:"required"`
	CreatedAt int64  `json:"-"`
	UpdatedAt int64  `json:"-"`
}

func (c *SendMessageRequest) Validate() error {
	validate := validator.New(validator.WithRequiredStructEnabled())
	return validate.Struct(c)
}

type MessageResponse struct {
	Id        int    `json:"id"`
	SenderID  int    `json:"sender_id"`
	ChatID    int    `json:"chat_id"`
	MsgType   string `json:"msg_type"`
	Text      string `json:"text"`
	IsDeleted bool   `json:"is_deleted"`
	CreatedAt int64  `json:"created_at"`
	UpdatedAt int64  `json:"updated_at"`
}

type GetMessagesRequest struct {
	UserID  int
	ChatID  int    `json:"chat_id" db:"chat_id" validate:"required"`
	MsgType string `json:"msg_type" db:"msg_type" validate:"required"`
	Limit   int    `json:"limit" db:"limit" validate:"gt=0,max=20,required"`
	Offset  int    `json:"offset" db:"offset" validate:"gte=0"`
}

func (c *GetMessagesRequest) Validate() error {
	validate := validator.New(validator.WithRequiredStructEnabled())
	return validate.Struct(c)
}

type GetChatMessagesRequest struct {
	ChatID       int `json:"chat_id" validate:"required"`
	InitialMsgID int `json:"initial_msg_id" validate:"gte=0"`
	Before       int `json:"before" validate:"gte=0"`
	After        int `json:"after" validate:"gte=0"`
}

func (c *GetChatMessagesRequest) Validate() error {
	validate := validator.New(validator.WithRequiredStructEnabled())
	return validate.Struct(c)
}

type AddToChatRequest struct {
	ChatID int    `json:"chat_id" db:"chat_id" validate:"gt=0,required"`
	UserID int    `json:"user_id" db:"user_id" validate:"gt=0,required"`
	Role   string `json:"role" db:"role" validate:"required"`
}

func (a *AddToChatRequest) Validate() error {
	validate := validator.New(validator.WithRequiredStructEnabled())
	return validate.Struct(a)
}

type AddToChatResponse struct {
	Status string `json:"status"`
}

type StateResponse struct {
	UserID   int       `json:"user_id"`
	Connects []Connect `json:"connects"`
}

type Connect struct {
	Conn        bool      `json:"conn"`
	Pubsub      uuid.UUID `json:"pubsub"`
	CreatedAt   int64     `json:"created_at"`
	CurrentChat int       `json:"current_chat"`
}
