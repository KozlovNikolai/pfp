package httpserver

import (
	"github.com/go-playground/validator/v10"
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

type UserChatResponse struct {
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
