package httpserver

import (
	"github.com/go-playground/validator/v10"
)

// UserRequest is ...
type UserRequest struct {
	Login    string `json:"login"    db:"login"    example:"cmd@cmd.ru" validate:"required,email"`
	Password string `json:"password" db:"password" example:"123456"     validate:"required,min=6,max=32"`
}

// Validate ...
func (u *UserRequest) Validate() error {
	validate := validator.New(validator.WithRequiredStructEnabled())
	return validate.Struct(u)
}

// UserResponse ...
type UserResponse struct {
	ID       int    `json:"id"       db:"id"`
	Login    string `json:"login"    db:"login"`
	Password string `json:"password" db:"password"`
	Role     string `json:"role"     db:"role"`
	Token    string `json:"token"    db:"token"`
}

type UserChatResponse struct {
	ID        int    `json:"id" db:"id"`
	UserExtID int    `json:"user_ext_id" db:"user_ext_id"`
	Name      string `json:"name" db:"name"`
	Surname   string `json:"surname" db:"surname"`
	Email     string `json:"email" db:"email"`
	UserType  string `json:"type" db:"type"`
	CreatedAt uint64 `json:"created_at" db:"created_at"`
	UpdatedAt uint64 `json:"updated_at" db:"updated_at"`
}

// type ReceiveUser struct {
// 	Payload Payload `json:"payload"`
// }

// type Payload struct {
// 	Email         string        `json:"email" db:"email"`
// 	Name          string        `json:"name" db:"name"`
// 	Surname       string        `json:"surname" db:"surname"`
// 	UserID        int           `json:"id" db:"id"`
// 	TokensSputnik TokensSputnik `json:"tokens"`
// 	Lang          string        `json:"lang" db:"lang"`
// 	AccountId     int           `json:"account_id"`
// }

// type TokensSputnik struct {
// 	AccessToken  string `json:"access_token"`
// 	RefreshToken string `json:"refresh_token"`
// }

// func (ru *ReceiveUser) Validate() error {
// 	validate := validator.New(validator.WithRequiredStructEnabled())
// 	return validate.Struct(ru)
// }
