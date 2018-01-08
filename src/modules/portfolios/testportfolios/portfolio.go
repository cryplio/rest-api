package testportfolios

import (
	"testing"

	"github.com/Nivl/go-rest-tools/security/auth"
	"github.com/Nivl/go-rest-tools/security/auth/testauth"
	sqldb "github.com/Nivl/go-sqldb"
	"github.com/cryplio/rest-api/src/modules/portfolios"
	uuid "github.com/satori/go.uuid"
)

// NewPortfolio returns a non persisted portfolio
func NewPortfolio() *portfolios.Portfolio {
	u := testauth.NewUser()
	return &portfolios.Portfolio{
		ID:     uuid.NewV4().String(),
		UserID: u.ID,
		Name:   "My Portfolio",
		User:   u,
	}
}

// NewPersistedPortfolio creates and persists a new portfolio
func NewPersistedPortfolio(t *testing.T, q sqldb.Queryable, u *auth.User, p *portfolios.Portfolio) *portfolios.Portfolio {
	if p == nil {
		p = &portfolios.Portfolio{}
	}
	if p.Name == "" {
		p.Name = "My Portfolio"
	}

	if u == nil {
		u = testauth.NewPersistedUser(t, q, nil)
	}
	p.User = u
	p.UserID = u.ID

	if err := p.Create(q); err != nil {
		t.Fatalf("failed to create portfolio: %s", err)
	}
	return p
}
