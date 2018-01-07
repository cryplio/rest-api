package testusers

import (
	"testing"

	"github.com/Nivl/go-rest-tools/security/auth"
	"github.com/Nivl/go-rest-tools/security/auth/testauth"
	sqldb "github.com/Nivl/go-sqldb"
	"github.com/cryplio/rest-api/src/modules/users"
	uuid "github.com/satori/go.uuid"
)

// NewProfile returns a non persisted profile and user
func NewProfile() *users.Profile {
	u := testauth.NewUser()
	return &users.Profile{
		ID:     uuid.NewV4().String(),
		UserID: u.ID,
		User:   u,
	}
}

// NewPersistedProfile creates and persists a new user and user profile with "fake" as password
func NewPersistedProfile(t *testing.T, q sqldb.Queryable, p *users.Profile) *users.Profile {
	if p == nil {
		p = &users.Profile{}
	}

	p.User = testauth.NewPersistedUser(t, q, p.User)
	p.UserID = p.User.ID

	if err := p.Create(q); err != nil {
		t.Fatalf("failed to create user: %s", err)
	}
	return p
}

// NewAuth creates a new user and their session
func NewAuth(t *testing.T, q sqldb.Queryable) (*auth.User, *auth.Session) {
	p, session := NewAuthProfile(t, q)
	return p.User, session
}

// NewAuthProfile creates a new profile, user, and their session
func NewAuthProfile(t *testing.T, q sqldb.Queryable) (*users.Profile, *auth.Session) {
	user, session := testauth.NewPersistedAuth(t, q)
	p := NewProfile()
	p.ID = ""
	p.User = user
	p.UserID = user.ID
	if err := p.Create(q); err != nil {
		t.Fatalf("failed to create a new auth with profile: %s", err)
	}
	return p, session
}

// NewAdminAuth creates a new admin and their session
func NewAdminAuth(t *testing.T, q sqldb.Queryable) (*auth.User, *auth.Session) {
	user, session := testauth.NewPersistedAdminAuth(t, q)
	p := NewProfile()
	p.ID = ""
	p.User = user
	p.UserID = user.ID
	if err := p.Create(q); err != nil {
		t.Fatalf("failed to create a new admin auth with profile: %s", err)
	}
	return user, session
}
