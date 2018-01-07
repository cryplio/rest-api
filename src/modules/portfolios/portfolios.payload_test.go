package portfolios_test

import (
	"testing"

	"github.com/cryplio/rest-api/src/modules/portfolios"
	"github.com/cryplio/rest-api/src/modules/portfolios/testportfolios"

	"github.com/stretchr/testify/assert"
)

func TestExportPublicPortfolios(t *testing.T) {
	list := portfolios.Portfolios{
		testportfolios.NewPortfolio(),
		testportfolios.NewPortfolio(),
		testportfolios.NewPortfolio(),
	}
	pld := list.ExportPublic()
	assert.Equal(t, len(list), len(pld.Results), "Wrong number of portfolios returned")

	for _, p := range pld.Results {
		assert.NotEmpty(t, p.Name, "Nane should have been set")
	}
}

func TestExportPrivatePortfolios(t *testing.T) {
	list := portfolios.Portfolios{
		testportfolios.NewPortfolio(),
		testportfolios.NewPortfolio(),
		testportfolios.NewPortfolio(),
	}
	pld := list.ExportPrivate()
	assert.Equal(t, len(list), len(pld.Results), "Wrong number of portfolios returned")

	for i, p := range pld.Results {
		assert.Equal(t, list[i].Name, p.Name, "Nane should have been set")
	}
}
