package tickers

// Code generated; DO NOT EDIT.

import (
	"errors"
	
	"fmt"
	"strings"
	

	"github.com/Nivl/go-rest-tools/types/apierror"
	"github.com/Nivl/go-types/datetime"
	"github.com/Nivl/go-sqldb"
	
)

// JoinSQL returns a string ready to be embed in a JOIN query
func JoinSQL(prefix string) string {
	fields := []string{ "id", "created_at", "updated_at", "deleted_at", "name", "unit", "marketcap", "volume_24h", "max_supply", "current_supply", "logo_url", "website", "price_usd", "percent_change_1h", "percent_change_24h", "percent_change_7d", "coinmarketcap_id" }
	output := ""

	for _, field := range fields {
		fullName := fmt.Sprintf("%s.%s", prefix, field)
		output += fmt.Sprintf("%s \"%s\", ", fullName, fullName)
	}
	return strings.TrimSuffix(output, ", ")
}

// GetByID finds and returns an active ticker by ID
// Deleted object are not returned
func GetByID(q sqldb.Queryable, id string) (*Ticker, error) {
	ti := &Ticker{}
	stmt := "SELECT * from tickers WHERE id=$1 and deleted_at IS NULL LIMIT 1"
	err := q.Get(ti, stmt, id)
	return ti, apierror.NewFromSQL(err)
}

// GetAnyByID finds and returns an ticker by ID.
// Deleted object are returned
func GetAnyByID(q sqldb.Queryable, id string) (*Ticker, error) {
	ti := &Ticker{}
	stmt := "SELECT * from tickers WHERE id=$1 LIMIT 1"
	err := q.Get(ti, stmt, id)
	return ti, apierror.NewFromSQL(err)
}




// Create persists a ticker in the database
func (ti *Ticker) Create(q sqldb.Queryable) error {
	

	return ti.doCreate(q)
}

// doCreate persists a ticker in the database using a Node
func (ti *Ticker) doCreate(q sqldb.Queryable) error {
	
	ti.UpdatedAt = datetime.Now()
	if ti.CreatedAt == nil {
		ti.CreatedAt = datetime.Now()
	}

	stmt := "INSERT INTO tickers (id, created_at, updated_at, deleted_at, name, unit, marketcap, volume_24h, max_supply, current_supply, logo_url, website, price_usd, percent_change_1h, percent_change_24h, percent_change_7d, coinmarketcap_id) VALUES (:id, :created_at, :updated_at, :deleted_at, :name, :unit, :marketcap, :volume_24h, :max_supply, :current_supply, :logo_url, :website, :price_usd, :percent_change_1h, :percent_change_24h, :percent_change_7d, :coinmarketcap_id)"
	_, err := q.NamedExec(stmt, ti)

  return apierror.NewFromSQL(err)
}

// Update updates most of the fields of a persisted ticker
// Excluded fields are id, created_at, deleted_at, etc.
func (ti *Ticker) Update(q sqldb.Queryable) error {
	if ti.ID == "" {
		return errors.New("cannot update a non-persisted ticker")
	}

	return ti.doUpdate(q)
}

// doUpdate updates a ticker in the database
func (ti *Ticker) doUpdate(q sqldb.Queryable) error {
	if ti.ID == "" {
		return errors.New("cannot update a non-persisted ticker")
	}

	ti.UpdatedAt = datetime.Now()

	stmt := "UPDATE tickers SET id=:id, created_at=:created_at, updated_at=:updated_at, deleted_at=:deleted_at, name=:name, unit=:unit, marketcap=:marketcap, volume_24h=:volume_24h, max_supply=:max_supply, current_supply=:current_supply, logo_url=:logo_url, website=:website, price_usd=:price_usd, percent_change_1h=:percent_change_1h, percent_change_24h=:percent_change_24h, percent_change_7d=:percent_change_7d, coinmarketcap_id=:coinmarketcap_id WHERE id=:id"
	_, err := q.NamedExec(stmt, ti)

	return apierror.NewFromSQL(err)
}

// Delete removes a ticker from the database
func (ti *Ticker) Delete(q sqldb.Queryable) error {
	if ti.ID == "" {
		return errors.New("ticker has not been saved")
	}

	stmt := "DELETE FROM tickers WHERE id=$1"
	_, err := q.Exec(stmt, ti.ID)

	return err
}

// IsZero checks if the object is either nil or don't have an ID
func (ti *Ticker) IsZero() bool {
	return ti == nil || ti.ID == ""
}