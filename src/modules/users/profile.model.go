package users

import (
	"github.com/Nivl/go-rest-tools/security/auth"
	"github.com/Nivl/go-rest-tools/types/apierror"
	sqldb "github.com/Nivl/go-sqldb"
	"github.com/Nivl/go-types/datetime"
)

// Profile represents the public information of a user
//go:generate api-cli generate model Profile -t user_profiles -e Get,GetAny,JoinSQL --single=false
type Profile struct {
	ID        string             `db:"id"`
	CreatedAt *datetime.DateTime `db:"created_at"`
	UpdatedAt *datetime.DateTime `db:"updated_at"`
	DeletedAt *datetime.DateTime `db:"deleted_at"`

	UserID string `db:"user_id"`

	// Embedded models
	*auth.User `db:"users"`
}

// Profiles represents a list of Profile
type Profiles []*Profile

// GetByIDWithProfile finds and returns an active user with their profile by ID
// Deleted object are not returned
func GetByIDWithProfile(q sqldb.Queryable, id string) (*Profile, error) {
	u := &Profile{}
	stmt := `
	SELECT profile.*, ` + auth.JoinUserSQL("users") + `
	FROM user_profiles profile
	JOIN users
	  ON users.id = profile.user_id
	WHERE users.id=$1
	  AND users.deleted_at IS NULL
	LIMIT 1`
	err := q.Get(u, stmt, id)
	return u, apierror.NewFromSQL(err)
}
