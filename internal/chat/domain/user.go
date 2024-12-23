package domain

import "time"

// User is a domain User
type User struct {
	id        int
	userExtID int
	login     string
	password  string
	profile   string
	name      string
	surname   string
	email     string
	userType  string
	createdAt int64
	updatedAt int64
}

// NewUserData is a domain User
type NewUserData struct {
	ID        int    `json:"id" db:"id"`
	UserExtID int    `json:"user_ext_id" db:"user_ext_id"`
	Login     string `json:"login" db:"login"`
	Password  string `json:"password" db:"password"`
	Profile   string `json:"profile" db:"profile"`
	Name      string `json:"name" db:"name"`
	Surname   string `json:"surname" db:"surname"`
	Email     string `json:"email" db:"email"`
	UserType  string `json:"type" db:"type"`
	CreatedAt int64  `json:"created_at" db:"created_at"`
	UpdatedAt int64  `json:"updated_at" db:"updated_at"`
}

// NewUser ...
func NewUser(data NewUserData) User {
	createdTime := time.Now().Unix()
	return User{
		id:        data.ID,
		userExtID: data.UserExtID,
		login:     data.Login,
		password:  data.Password,
		profile:   data.Profile,
		name:      data.Name,
		surname:   data.Surname,
		email:     data.Email,
		userType:  data.UserType,
		createdAt: createdTime,
		updatedAt: createdTime,
	}
}

// ID ...
func (u User) ID() int {
	return u.id
}

// Login ...
func (u User) UserExtID() int {
	return u.userExtID
}

func (u User) Login() string {
	return u.login
}
func (u User) Password() string {
	return u.password
}
func (u User) Profile() string {
	return u.profile
}

// Password ...
func (u User) Name() string {
	return u.name
}

// Role ...
func (u User) Surname() string {
	return u.surname
}

// Token ...
func (u User) Email() string {
	return u.email
}

// Token ...
func (u User) UserType() string {
	return u.userType
}

// CreratedAt ...
func (u User) CreatedAt() int64 {
	return u.createdAt
}

// UpdatedAt ...
func (u User) UpdatedAt() int64 {
	return u.updatedAt
}
