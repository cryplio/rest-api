package testauth

import (
	"testing"

	"github.com/Nivl/go-rest-tools/security/auth"
	db "github.com/Nivl/go-sqldb"
)

// NewAuth creates a non-persisted user and their session
func NewAuth() (*auth.User, *auth.Session) {
	user := NewUser()
	session := NewSession(user)
	return user, session
}

// NewPersistedAuth creates a persisted new user and their session
func NewPersistedAuth(t *testing.T, q db.Queryable) (*auth.User, *auth.Session) {
	user := NewPersistedUser(t, q, nil)
	session := NewPersistedSession(t, q, user)
	return user, session
}

// NewAdminAuth creates a new non-persisted admin and their session
func NewAdminAuth() (*auth.User, *auth.Session) {
	user, session := NewAuth()
	user.IsAdmin = true
	return user, session
}

// NewPersistedAdminAuth creates a new admin and their session
func NewPersistedAdminAuth(t *testing.T, q db.Queryable) (*auth.User, *auth.Session) {
	user := NewPersistedUser(t, q, &auth.User{IsAdmin: true})
	session := NewPersistedSession(t, q, user)
	return user, session
}

// NewSession creates a non-persisted session for the given user
func NewSession(user *auth.User) *auth.Session {
	return &auth.Session{
		UserID: user.ID,
	}
}

// NewPersistedSession creates and persists a new session for the given user
func NewPersistedSession(t *testing.T, q db.Queryable, user *auth.User) *auth.Session {
	session := &auth.Session{
		UserID: user.ID,
	}
	if err := session.Create(q); err != nil {
		t.Fatal(err)
	}
	return session
}
