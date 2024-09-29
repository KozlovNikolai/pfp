package httpserver

import (
	"github.com/go-playground/validator/v10"
)

// #######################################################################################3
// UserRequest is ...
type UserRequest struct {
	Login    string `json:"login" db:"login" example:"cmd@cmd.ru" validate:"required,email"`
	Password string `json:"password" db:"password" example:"123456" validate:"required,min=6,max=32"`
}

func (u *UserRequest) Validate() error {
	validate := validator.New(validator.WithRequiredStructEnabled())
	return validate.Struct(u)
}

type UserResponse struct {
	ID       int    `json:"id" db:"id"`
	Login    string `json:"login" db:"login"`
	Password string `json:"password" db:"password"`
	Role     string `json:"role" db:"role"`
	Token    string `json:"token" db:"token"`
}
