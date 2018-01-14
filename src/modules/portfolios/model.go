package portfolios

import (
	"github.com/Nivl/go-rest-tools/security/auth"
	"github.com/Nivl/go-types/datetime"
)

// Portfolio represents a coin/token collection of a user
//go:generate api-cli generate model Portfolio -t user_portfolios --single=false
type Portfolio struct {
	ID        string             `db:"id"`
	CreatedAt *datetime.DateTime `db:"created_at"`
	UpdatedAt *datetime.DateTime `db:"updated_at"`
	DeletedAt *datetime.DateTime `db:"deleted_at"`

	UserID string `db:"user_id"`
	Name   string `db:"name"`

	// Embedded models
	*auth.User `db:"users"`
}

// Portfolios represents a list of Portfolio
type Portfolios []*Portfolio
