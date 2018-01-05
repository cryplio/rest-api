package auth

import (
	"github.com/Nivl/go-types/datetime"
)

// User is a structure representing a user that can be saved in the database
//go:generate api-cli generate model User -t users --single=false
type User struct {
	ID        string             `db:"id"`
	CreatedAt *datetime.DateTime `db:"created_at"`
	UpdatedAt *datetime.DateTime `db:"updated_at"`
	DeletedAt *datetime.DateTime `db:"deleted_at"`

	Name     string `db:"name"`
	Email    string `db:"email"`
	Password string `db:"password"`
	IsAdmin  bool   `db:"is_admin"`
}

// IsLogged checks if the user object belong to a logged in user
// Works on nil object
func (u *User) IsLogged() bool {
	return u != nil
}

// IsAdm checks if the user object belong to a logged in admin
// Works on nil object
func (u *User) IsAdm() bool {
	return u.IsLogged() && u.IsAdmin
}
