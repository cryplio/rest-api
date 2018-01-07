package users

// Code generated; DO NOT EDIT.

import (
	"errors"
	

	"github.com/Nivl/go-rest-tools/types/apierror"
	"github.com/Nivl/go-types/datetime"
	"github.com/Nivl/go-sqldb"
	uuid "github.com/satori/go.uuid"
)








// Save creates or updates the article depending on the value of the id using
// a transaction
func (p *Profile) Save(q sqldb.Queryable) error {
	if p.ID == "" {
		return p.Create(q)
	}

	return p.Update(q)
}

// Create persists a profile in the database
func (p *Profile) Create(q sqldb.Queryable) error {
	if p.ID != "" {
		return errors.New("cannot persist a profile that already has an ID")
	}

	return p.doCreate(q)
}

// doCreate persists a profile in the database using a Node
func (p *Profile) doCreate(q sqldb.Queryable) error {
	p.ID = uuid.NewV4().String()
	p.UpdatedAt = datetime.Now()
	if p.CreatedAt == nil {
		p.CreatedAt = datetime.Now()
	}

	stmt := "INSERT INTO user_profiles (id, created_at, updated_at, deleted_at, user_id) VALUES (:id, :created_at, :updated_at, :deleted_at, :user_id)"
	_, err := q.NamedExec(stmt, p)

  return apierror.NewFromSQL(err)
}

// Update updates most of the fields of a persisted profile
// Excluded fields are id, created_at, deleted_at, etc.
func (p *Profile) Update(q sqldb.Queryable) error {
	if p.ID == "" {
		return errors.New("cannot update a non-persisted profile")
	}

	return p.doUpdate(q)
}

// doUpdate updates a profile in the database
func (p *Profile) doUpdate(q sqldb.Queryable) error {
	if p.ID == "" {
		return errors.New("cannot update a non-persisted profile")
	}

	p.UpdatedAt = datetime.Now()

	stmt := "UPDATE user_profiles SET id=:id, created_at=:created_at, updated_at=:updated_at, deleted_at=:deleted_at, user_id=:user_id WHERE id=:id"
	_, err := q.NamedExec(stmt, p)

	return apierror.NewFromSQL(err)
}

// Delete removes a profile from the database
func (p *Profile) Delete(q sqldb.Queryable) error {
	if p.ID == "" {
		return errors.New("profile has not been saved")
	}

	stmt := "DELETE FROM user_profiles WHERE id=$1"
	_, err := q.Exec(stmt, p.ID)

	return err
}

// GetID returns the ID field
func (p *Profile) GetID() string {
	return p.ID
}

// SetID sets the ID field
func (p *Profile) SetID(id string) {
	p.ID = id
}

// IsZero checks if the object is either nil or don't have an ID
func (p *Profile) IsZero() bool {
	return p == nil || p.ID == ""
}