package tickers

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

// JoinSQL returns a string ready to be embed in a JOIN query
func JoinSQL(prefix string) string {
	fields := []string{ "id", "created_at", "updated_at", "deleted_at", "name", "symbol", "unit", "marketcap", "volume_24h", "max_supply", "current_supply", "logo_url", "website", "price_usd", "percent_change_1h", "percent_change_24h", "percent_change_7d", "coinmarketcap_id" }
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
	t := &Ticker{}
	stmt := "SELECT * from tickers WHERE id=$1 and deleted_at IS NULL LIMIT 1"
	err := q.Get(t, stmt, id)
	return t, apierror.NewFromSQL(err)
}

// GetAnyByID finds and returns an ticker by ID.
// Deleted object are returned
func GetAnyByID(q sqldb.Queryable, id string) (*Ticker, error) {
	t := &Ticker{}
	stmt := "SELECT * from tickers WHERE id=$1 LIMIT 1"
	err := q.Get(t, stmt, id)
	return t, apierror.NewFromSQL(err)
}


// Save creates or updates the article depending on the value of the id using
// a transaction
func (t *Ticker) Save(q sqldb.Queryable) error {
	if t.ID == "" {
		return t.Create(q)
	}

	return t.Update(q)
}

// Create persists a ticker in the database
func (t *Ticker) Create(q sqldb.Queryable) error {
	if t.ID != "" {
		return errors.New("cannot persist a ticker that already has an ID")
	}

	return t.doCreate(q)
}

// doCreate persists a ticker in the database using a Node
func (t *Ticker) doCreate(q sqldb.Queryable) error {
	t.ID = uuid.NewV4().String()
	t.UpdatedAt = datetime.Now()
	if t.CreatedAt == nil {
		t.CreatedAt = datetime.Now()
	}

	stmt := "INSERT INTO tickers (id, created_at, updated_at, deleted_at, name, symbol, unit, marketcap, volume_24h, max_supply, current_supply, logo_url, website, price_usd, percent_change_1h, percent_change_24h, percent_change_7d, coinmarketcap_id) VALUES (:id, :created_at, :updated_at, :deleted_at, :name, :symbol, :unit, :marketcap, :volume_24h, :max_supply, :current_supply, :logo_url, :website, :price_usd, :percent_change_1h, :percent_change_24h, :percent_change_7d, :coinmarketcap_id)"
	_, err := q.NamedExec(stmt, t)

  return apierror.NewFromSQL(err)
}

// Update updates most of the fields of a persisted ticker
// Excluded fields are id, created_at, deleted_at, etc.
func (t *Ticker) Update(q sqldb.Queryable) error {
	if t.ID == "" {
		return errors.New("cannot update a non-persisted ticker")
	}

	return t.doUpdate(q)
}

// doUpdate updates a ticker in the database
func (t *Ticker) doUpdate(q sqldb.Queryable) error {
	if t.ID == "" {
		return errors.New("cannot update a non-persisted ticker")
	}

	t.UpdatedAt = datetime.Now()

	stmt := "UPDATE tickers SET id=:id, created_at=:created_at, updated_at=:updated_at, deleted_at=:deleted_at, name=:name, symbol=:symbol, unit=:unit, marketcap=:marketcap, volume_24h=:volume_24h, max_supply=:max_supply, current_supply=:current_supply, logo_url=:logo_url, website=:website, price_usd=:price_usd, percent_change_1h=:percent_change_1h, percent_change_24h=:percent_change_24h, percent_change_7d=:percent_change_7d, coinmarketcap_id=:coinmarketcap_id WHERE id=:id"
	_, err := q.NamedExec(stmt, t)

	return apierror.NewFromSQL(err)
}

// Delete removes a ticker from the database
func (t *Ticker) Delete(q sqldb.Queryable) error {
	if t.ID == "" {
		return errors.New("ticker has not been saved")
	}

	stmt := "DELETE FROM tickers WHERE id=$1"
	_, err := q.Exec(stmt, t.ID)

	return err
}

// GetID returns the ID field
func (t *Ticker) GetID() string {
	return t.ID
}

// SetID sets the ID field
func (t *Ticker) SetID(id string) {
	t.ID = id
}

// IsZero checks if the object is either nil or don't have an ID
func (t *Ticker) IsZero() bool {
	return t == nil || t.ID == ""
}