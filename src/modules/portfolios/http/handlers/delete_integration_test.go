// +build integration

package handlers_test

import (
	"database/sql"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/cryplio/rest-api/src/modules/api"
	"github.com/cryplio/rest-api/src/modules/portfolios/http/handlers"
	"github.com/cryplio/rest-api/src/modules/portfolios/testportfolios"

	"github.com/Nivl/go-rest-tools/dependencies"
	"github.com/Nivl/go-rest-tools/network/http/httptests"
	"github.com/Nivl/go-rest-tools/security/auth"
	"github.com/Nivl/go-rest-tools/security/auth/testauth"
	"github.com/Nivl/go-rest-tools/testing/integration"
	"github.com/stretchr/testify/assert"
)

func TestDelete(t *testing.T) {
	t.Parallel()

	helper, err := integration.New(NewDeps(), migrationFolder)
	if err != nil {
		panic(err)
	}
	defer helper.Close()
	dbCon := helper.Deps.DB()

	u1, s1 := testauth.NewPersistedAuth(t, dbCon)
	portfolioOfU1ToDelete := testportfolios.NewPersistedPortfolio(t, dbCon, u1, nil)
	portfolioOfU1ToKeep := testportfolios.NewPersistedPortfolio(t, dbCon, u1, nil)

	_, s2 := testauth.NewPersistedAuth(t, dbCon)

	tests := []struct {
		description string
		code        int
		params      *handlers.DeleteParams
		auth        *httptests.RequestAuth
	}{
		{
			"Not logged",
			http.StatusUnauthorized,
			&handlers.DeleteParams{ID: portfolioOfU1ToKeep.ID},
			nil,
		},
		{
			"Deleting a porfolio that does not exist",
			http.StatusNotFound,
			&handlers.DeleteParams{ID: portfolioOfU1ToKeep.ID},
			httptests.NewRequestAuth(s2),
		},
		{
			"Deleting the porfolio of an other user",
			http.StatusNotFound,
			&handlers.DeleteParams{ID: "37dc1575-a296-48f5-8345-99e0649959c9"},
			httptests.NewRequestAuth(s2),
		},
		{
			"Deleting a portfolio",
			http.StatusNoContent,
			&handlers.DeleteParams{ID: portfolioOfU1ToDelete.ID},
			httptests.NewRequestAuth(s1),
		},
	}

	t.Run("parallel", func(t *testing.T) {
		for _, tc := range tests {
			tc := tc
			t.Run(tc.description, func(t *testing.T) {
				t.Parallel()
				defer helper.RecoverPanic()

				rec := callDelete(t, tc.params, tc.auth, helper.Deps)
				assert.Equal(t, tc.code, rec.Code)

				if rec.Code == http.StatusNoContent {
					// We check that the user has been deleted
					var user auth.User
					stmt := "SELECT * FROM user_portfolios WHERE id=$1 LIMIT 1"
					err := dbCon.Get(&user, stmt, tc.params.ID)
					assert.Equal(t, sql.ErrNoRows, err, "Portfolio not deleted")
				}
			})
		}
	})
}

func callDelete(t *testing.T, params *handlers.DeleteParams, auth *httptests.RequestAuth, deps dependencies.Dependencies) *httptest.ResponseRecorder {
	ri := &httptests.RequestInfo{
		Endpoint: handlers.Endpoints[handlers.EndpointDelete],
		Params:   params,
		Auth:     auth,
		Router:   api.GetRouter(deps),
	}
	return httptests.NewRequest(t, ri)
}
