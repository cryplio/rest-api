package integration

import (
	"fmt"
	"strings"

	"github.com/Nivl/go-rest-tools/dependencies"
	db "github.com/Nivl/go-sqldb"
	sqlx "github.com/Nivl/go-sqldb/implementations/sqlxdb"
	"github.com/pkg/errors"
	"github.com/pressly/goose"
	uuid "github.com/satori/go.uuid"
)

// Wrapper is an helper to simplify the encapsulation of integration
// tests
type Wrapper struct {
	Deps      dependencies.Dependencies
	masterDB  db.Connection
	tmpDBName string
}

// New creates a new Wrapper
func New(deps dependencies.Dependencies, migrationFolder string) (*Wrapper, error) {
	it := &Wrapper{
		Deps: deps,

		// masterDB contains a DB connection to the default database.
		// this is needed because postgres won't allow us to drop the
		// current database, so we use the default one to drop the
		// temporary one
		masterDB: deps.DB(),

		// for the new table name we use a uuid without "-"
		tmpDBName: strings.Replace(uuid.NewV4().String(), "-", "", -1),
	}

	// We get the current database name
	masterDbName := ""
	if err := it.masterDB.Get(&masterDbName, "SELECT current_database();"); err != nil {
		return nil, errors.Wrap(err, "failed getting master database name")
	}

	// We create a new Database to avoid races between tests
	stmt := fmt.Sprintf(`CREATE DATABASE "%s";`, it.tmpDBName)
	if _, err := it.masterDB.Exec(stmt); err != nil {
		return nil, errors.Wrap(err, "failed creating tmp database")
	}

	// We set this new database as the new current one, to do that we
	// get the DSN of the current connection and swap the table name by
	// the new one
	masterDBString := fmt.Sprintf("dbname=%s", masterDbName)
	tmpDBString := fmt.Sprintf("dbname=%s", it.tmpDBName)
	tmpDBDSN := strings.Replace(it.masterDB.DSN(), masterDBString, tmpDBString, -1)
	tmpDB, err := sqlx.New(tmpDBDSN)
	if err != nil {
		return nil, errors.Wrap(err, "could not connect to the tmp table")
	}
	it.Deps.SetDB(tmpDB)

	// We apply the migration to the newly created database
	if err := goose.Up(it.Deps.DB().SQL(), migrationFolder); err != nil {
		return nil, errors.Wrap(err, "could not execute the migrations")
	}

	return it, nil
}

// Close cleans up the tests by deleting the database
func (it *Wrapper) Close() error {
	if err := it.Deps.DB().Close(); err != nil {
		return err
	}
	stmt := fmt.Sprintf(`DROP DATABASE "%s";`, it.tmpDBName)
	if _, err := it.masterDB.Exec(stmt); err != nil {
		return err
	}
	return it.masterDB.Close()
}

// RecoverPanic prevents a panic from not calling the defer in the othr goroutines
func (it *Wrapper) RecoverPanic() {
	recover()
}

// CloseOnPanic cleans up the tests and re-panic
func (it *Wrapper) CloseOnPanic() {
	if rec := recover(); rec != nil {
		it.Close()
		panic(fmt.Sprintf("%v", rec))
	}
}
