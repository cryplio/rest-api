package testauth

import (
	"fmt"
	"testing"

	"github.com/satori/go.uuid"

	"github.com/Nivl/go-hasher/implementations/bcrypt"
	"github.com/Nivl/go-rest-tools/security/auth"
	db "github.com/Nivl/go-sqldb"
	"github.com/dchest/uniuri"
)

// NewUser returns a non-persisted user with "fake" as password
func NewUser() *auth.User {
	hashr := bcrypt.Bcrypt{}
	password, _ := hashr.Hash("fake")

	return &auth.User{
		ID:       uuid.NewV4().String(),
		Name:     uniuri.New(),
		Email:    fmt.Sprintf("%s@domain.tld", uniuri.New()),
		Password: password,
	}
}

// NewAdmin returns a non-persisted admin user with "fake" as password
func NewAdmin() *auth.User {
	user := NewUser()
	user.IsAdmin = true
	return user
}

// NewPersistedUser creates and persists a new user with "fake" as password
func NewPersistedUser(t *testing.T, q db.Queryable, u *auth.User) *auth.User {
	hashr := bcrypt.Bcrypt{}
	if u == nil {
		u = &auth.User{}
	}

	if u.Email == "" {
		u.Email = fmt.Sprintf("fake+%s@domain.tld", uniuri.New())
	}

	if u.Name == "" {
		u.Name = "Fake Account"
	}

	if u.Password == "" {
		var err error
		u.Password, err = hashr.Hash("fake")
		if err != nil {
			t.Fatalf("failed to create password: %s", err)
		}
	}

	if err := u.Create(q); err != nil {
		t.Fatalf("failed to create user: %s", err)
	}
	return u
}
