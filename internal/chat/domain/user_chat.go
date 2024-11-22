package domain

// User is a domain User
type UserChat struct {
	id        int
	userExtID string
	login     string
	password  string
	account   string
	token     string
	name      string
	surname   string
	email     string
	userType  string
	createdAt int64
	updatedAt int64
}

// NewUserData is a domain User
type NewUserChatData struct {
	ID        int    `json:"id" db:"id"`
	UserExtID string `json:"user_ext_id" db:"user_ext_id"`
	Login     string `json:"login" db:"login"`
	Password  string `json:"password" db:"password"`
	Account   string `json:"account" db:"account"`
	Token     string `json:"token" db:"token"`
	Name      string `json:"name" db:"name"`
	Surname   string `json:"surname" db:"surname"`
	Email     string `json:"email" db:"email"`
	UserType  string `json:"type" db:"type"`
	CreatedAt int64  `json:"created_at" db:"created_at"`
	UpdatedAt int64  `json:"updated_at" db:"updated_at"`
}

// NewUser ...
func NewUserChat(data NewUserChatData) UserChat {
	return UserChat{
		id:        data.ID,
		userExtID: data.UserExtID,
		login:     data.Login,
		password:  data.Password,
		account:   data.Account,
		token:     data.Token,
		name:      data.Name,
		surname:   data.Surname,
		email:     data.Email,
		userType:  data.UserType,
		createdAt: data.CreatedAt,
		updatedAt: data.UpdatedAt,
	}
}

// ID ...
func (u UserChat) ID() int {
	return u.id
}

// Login ...
func (u UserChat) UserExtID() string {
	return u.userExtID
}

func (u UserChat) Login() string {
	return u.login
}
func (u UserChat) Password() string {
	return u.password
}
func (u UserChat) Account() string {
	return u.account
}
func (u UserChat) Token() string {
	return u.token
}

// Password ...
func (u UserChat) Name() string {
	return u.name
}

// Role ...
func (u UserChat) Surname() string {
	return u.surname
}

// Token ...
func (u UserChat) Email() string {
	return u.email
}

// Token ...
func (u UserChat) UserType() string {
	return u.userType
}

// CreratedAt ...
func (u UserChat) CreatedAt() int64 {
	return u.createdAt
}

// UpdatedAt ...
func (u UserChat) UpdatedAt() int64 {
	return u.updatedAt
}
