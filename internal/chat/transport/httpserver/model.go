package httpserver

import (
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

// UserRequest is ...
type UserRequest struct {
	Account  string `json:"account" db:"account" example:"email" validate:"required"`
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
	UserExtID string `json:"user_ext_id" db:"user_ext_id"`
	Login     string `json:"login" db:"login"`
	// Password  string `json:"password" db:"password"`
	Account   string `json:"account" db:"account"`
	Token     string `json:"token" db:"token"`
	Name      string `json:"name" db:"name"`
	Surname   string `json:"surname" db:"surname"`
	Email     string `json:"email" db:"email"`
	UserType  string `json:"type" db:"type"`
	CreatedAt int64  `json:"created_at" db:"created_at"`
	UpdatedAt int64  `json:"updated_at" db:"updated_at"`
}

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
	OwnerID       int    `json:"owner_id"`
	ChatType      string `json:"chat_type"`
	LastChatMsgID uint64 `json:"last_message_id"`
	// Contacts      []domain.Contact `json:"contacts"`
	CreatedAt int64 `json:"created_at"`
	UpdatedAt int64 `json:"updated_at"`
}

type SendMessageRequest struct {
	SenderID  int
	ChatID    int    `json:"chat_id" db:"chat_id" validate:"required"`
	MsgType   string `json:"msg_type" db:"msg_type" validate:"required"`
	Text      string `json:"text" db:"text" validate:"required"`
	CreatedAt int64
	UpdatedAt int64
}

func (c *SendMessageRequest) Validate() error {
	validate := validator.New(validator.WithRequiredStructEnabled())
	return validate.Struct(c)
}

type MessageResponse struct {
	Id        int
	SenderID  int
	ChatID    int
	MsgType   string
	Text      string
	IsDeleted bool
	CreatedAt int64
	UpdatedAt int64
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

type AddToChatRequest struct {
	ChatID int `json:"chat_id" db:"chat_id" validate:"gt=0,required"`
	UserID int `json:"user_id" db:"user_id" validate:"gt=0,required"`
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
	Conn      bool      `json:"conn"`
	Pubsub    uuid.UUID `json:"pubsub"`
	CreatedAt int64     `json:"created_at"`
}
