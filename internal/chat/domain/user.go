package domain

import "time"

// User is a domain User
type User struct {
	id        int
	login     string
	password  string
	role      string
	token     string
	createdAt time.Time
	updatedAt time.Time
}

// NewUserData is a domain User
type NewUserData struct {
	ID        int
	Login     string
	Password  string
	Role      string
	Token     string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewUser(data NewUserData) User {
	return User{
		id:        data.ID,
		login:     data.Login,
		password:  data.Password,
		role:      data.Role,
		token:     data.Token,
		createdAt: data.CreatedAt,
		updatedAt: data.UpdatedAt,
	}
}

func (u User) ID() int {
	return u.id
}
func (u User) Login() string {
	return u.login
}
func (u User) Password() string {
	return u.password
}
func (u User) Role() string {
	return u.role
}
func (u User) Token() string {
	return u.token
}
func (u User) CreratedAt() time.Time {
	return u.createdAt
}
func (u User) UpdatedAt() time.Time {
	return u.updatedAt
}
