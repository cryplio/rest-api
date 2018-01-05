package auth

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

// JoinUserSQL returns a string ready to be embed in a JOIN query
func JoinUserSQL(prefix string) string {
	fields := []string{ "id", "created_at", "updated_at", "deleted_at", "name", "email", "password", "is_admin" }
	output := ""

	for _, field := range fields {
		fullName := fmt.Sprintf("%s.%s", prefix, field)
		output += fmt.Sprintf("%s \"%s\", ", fullName, fullName)
	}
	return strings.TrimSuffix(output, ", ")
}

// GetUserByID finds and returns an active user by ID
// Deleted object are not returned
func GetUserByID(q sqldb.Queryable, id string) (*User, error) {
	u := &User{}
	stmt := "SELECT * from users WHERE id=$1 and deleted_at IS NULL LIMIT 1"
	err := q.Get(u, stmt, id)
	return u, apierror.NewFromSQL(err)
}

// GetAnyUserByID finds and returns an user by ID.
// Deleted object are returned
func GetAnyUserByID(q sqldb.Queryable, id string) (*User, error) {
	u := &User{}
	stmt := "SELECT * from users WHERE id=$1 LIMIT 1"
	err := q.Get(u, stmt, id)
	return u, apierror.NewFromSQL(err)
}


// Save creates or updates the article depending on the value of the id using
// a transaction
func (u *User) Save(q sqldb.Queryable) error {
	if u.ID == "" {
		return u.Create(q)
	}

	return u.Update(q)
}

// Create persists a user in the database
func (u *User) Create(q sqldb.Queryable) error {
	if u.ID != "" {
		return errors.New("cannot persist a user that already has an ID")
	}

	return u.doCreate(q)
}

// doCreate persists a user in the database using a Node
func (u *User) doCreate(q sqldb.Queryable) error {
	u.ID = uuid.NewV4().String()
	u.UpdatedAt = datetime.Now()
	if u.CreatedAt == nil {
		u.CreatedAt = datetime.Now()
	}

	stmt := "INSERT INTO users (id, created_at, updated_at, deleted_at, name, email, password, is_admin) VALUES (:id, :created_at, :updated_at, :deleted_at, :name, :email, :password, :is_admin)"
	_, err := q.NamedExec(stmt, u)

  return apierror.NewFromSQL(err)
}

// Update updates most of the fields of a persisted user
// Excluded fields are id, created_at, deleted_at, etc.
func (u *User) Update(q sqldb.Queryable) error {
	if u.ID == "" {
		return errors.New("cannot update a non-persisted user")
	}

	return u.doUpdate(q)
}

// doUpdate updates a user in the database
func (u *User) doUpdate(q sqldb.Queryable) error {
	if u.ID == "" {
		return errors.New("cannot update a non-persisted user")
	}

	u.UpdatedAt = datetime.Now()

	stmt := "UPDATE users SET id=:id, created_at=:created_at, updated_at=:updated_at, deleted_at=:deleted_at, name=:name, email=:email, password=:password, is_admin=:is_admin WHERE id=:id"
	_, err := q.NamedExec(stmt, u)

	return apierror.NewFromSQL(err)
}

// Delete removes a user from the database
func (u *User) Delete(q sqldb.Queryable) error {
	if u.ID == "" {
		return errors.New("user has not been saved")
	}

	stmt := "DELETE FROM users WHERE id=$1"
	_, err := q.Exec(stmt, u.ID)

	return err
}

// GetID returns the ID field
func (u *User) GetID() string {
	return u.ID
}

// SetID sets the ID field
func (u *User) SetID(id string) {
	u.ID = id
}

// IsZero checks if the object is either nil or don't have an ID
func (u *User) IsZero() bool {
	return u == nil || u.ID == ""
}