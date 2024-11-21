package domain

// User is a domain User
type UserChat struct {
	id        int
	userExtID int
	name      string
	surname   string
	email     string
	userType  string
	createdAt uint64
	updatedAt uint64
}

// NewUserData is a domain User
type NewUserChatData struct {
	ID        int    `json:"id" db:"id"`
	UserExtID int    `json:"user_ext_id" db:"user_ext_id"`
	Name      string `json:"name" db:"name"`
	Surname   string `json:"surname" db:"surname"`
	Email     string `json:"email" db:"email"`
	UserType  string `json:"type" db:"type"`
	CreatedAt uint64 `json:"created_at" db:"created_at"`
	UpdatedAt uint64 `json:"updated_at" db:"updated_at"`
}

// NewUser ...
func NewUserChat(data NewUserChatData) UserChat {
	return UserChat{
		id:        data.ID,
		userExtID: data.UserExtID,
		name:      data.Name,
		surname:   data.Surname,
		email:     data.Email,
		userType:  data.UserType,
		createdAt: data.CreatedAt,
		updatedAt: data.UpdatedAt,
	}
}

// ID ...
func (u UserChat) GetID() int {
	return u.id
}

// Login ...
func (u UserChat) GetUserExtID() int {
	return u.userExtID
}

// Password ...
func (u UserChat) GetName() string {
	return u.name
}

// Role ...
func (u UserChat) GetSurname() string {
	return u.surname
}

// Token ...
func (u UserChat) GetEmail() string {
	return u.email
}

// Token ...
func (u UserChat) GetUserType() string {
	return u.userType
}

// CreratedAt ...
func (u UserChat) GetCreatedAt() uint64 {
	return u.createdAt
}

// UpdatedAt ...
func (u UserChat) GetUpdatedAt() uint64 {
	return u.updatedAt
}
