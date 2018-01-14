package portfolios

// Code generated; DO NOT EDIT.

import (
	"errors"
	
	"fmt"
	"strings"
	

	"github.com/Nivl/go-rest-tools/types/apierror"
	"github.com/Nivl/go-types/datetime"
	"github.com/Nivl/go-sqldb"
	uuid "github.com/satori/go.uuid"
)

// JoinPortfolioSQL returns a string ready to be embed in a JOIN query
func JoinPortfolioSQL(prefix string) string {
	fields := []string{ "id", "created_at", "updated_at", "deleted_at", "user_id", "name" }
	output := ""

	for _, field := range fields {
		fullName := fmt.Sprintf("%s.%s", prefix, field)
		output += fmt.Sprintf("%s \"%s\", ", fullName, fullName)
	}
	return strings.TrimSuffix(output, ", ")
}

// GetPortfolioByID finds and returns an active portfolio by ID
// Deleted object are not returned
func GetPortfolioByID(q sqldb.Queryable, id string) (*Portfolio, error) {
	p := &Portfolio{}
	stmt := "SELECT * from user_portfolios WHERE id=$1 and deleted_at IS NULL LIMIT 1"
	err := q.Get(p, stmt, id)
	return p, apierror.NewFromSQL(err)
}

// GetAnyPortfolioByID finds and returns an portfolio by ID.
// Deleted object are returned
func GetAnyPortfolioByID(q sqldb.Queryable, id string) (*Portfolio, error) {
	p := &Portfolio{}
	stmt := "SELECT * from user_portfolios WHERE id=$1 LIMIT 1"
	err := q.Get(p, stmt, id)
	return p, apierror.NewFromSQL(err)
}


// Save creates or updates the article depending on the value of the id using
// a transaction
func (p *Portfolio) Save(q sqldb.Queryable) error {
	if p.ID == "" {
		return p.Create(q)
	}

	return p.Update(q)
}

// Create persists a portfolio in the database
func (p *Portfolio) Create(q sqldb.Queryable) error {
	if p.ID != "" {
		return errors.New("cannot persist a portfolio that already has an ID")
	}

	return p.doCreate(q)
}

// doCreate persists a portfolio in the database using a Node
func (p *Portfolio) doCreate(q sqldb.Queryable) error {
	p.ID = uuid.NewV4().String()
	p.UpdatedAt = datetime.Now()
	if p.CreatedAt == nil {
		p.CreatedAt = datetime.Now()
	}

	stmt := "INSERT INTO user_portfolios (id, created_at, updated_at, deleted_at, user_id, name) VALUES (:id, :created_at, :updated_at, :deleted_at, :user_id, :name)"
	_, err := q.NamedExec(stmt, p)

  return apierror.NewFromSQL(err)
}

// Update updates most of the fields of a persisted portfolio
// Excluded fields are id, created_at, deleted_at, etc.
func (p *Portfolio) Update(q sqldb.Queryable) error {
	if p.ID == "" {
		return errors.New("cannot update a non-persisted portfolio")
	}

	return p.doUpdate(q)
}

// doUpdate updates a portfolio in the database
func (p *Portfolio) doUpdate(q sqldb.Queryable) error {
	if p.ID == "" {
		return errors.New("cannot update a non-persisted portfolio")
	}

	p.UpdatedAt = datetime.Now()

	stmt := "UPDATE user_portfolios SET id=:id, created_at=:created_at, updated_at=:updated_at, deleted_at=:deleted_at, user_id=:user_id, name=:name WHERE id=:id"
	_, err := q.NamedExec(stmt, p)

	return apierror.NewFromSQL(err)
}

// Delete removes a portfolio from the database
func (p *Portfolio) Delete(q sqldb.Queryable) error {
	if p.ID == "" {
		return errors.New("portfolio has not been saved")
	}

	stmt := "DELETE FROM user_portfolios WHERE id=$1"
	_, err := q.Exec(stmt, p.ID)

	return err
}

// GetID returns the ID field
func (p *Portfolio) GetID() string {
	return p.ID
}

// SetID sets the ID field
func (p *Portfolio) SetID(id string) {
	p.ID = id
}

// IsZero checks if the object is either nil or don't have an ID
func (p *Portfolio) IsZero() bool {
	return p == nil || p.ID == ""
}