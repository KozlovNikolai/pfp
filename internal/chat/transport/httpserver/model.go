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
