// +build integration

package portfolios_test

import (
	"os"
	"path"

	"github.com/Nivl/go-rest-tools/dependencies"
	"github.com/cryplio/rest-api/src/modules/api"
)

var (
	migrationFolder string
)

func NewDeps() dependencies.Dependencies {
	var err error
	_, deps, err := api.DefaultSetup()
	if err != nil {
		panic(err)
	}
	return deps
}

func init() {
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	migrationFolder = path.Join(wd, "..", "..", "..", "db", "migrations")
}
