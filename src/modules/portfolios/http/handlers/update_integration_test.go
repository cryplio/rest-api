// +build integration

package handlers_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/cryplio/rest-api/src/modules/portfolios"
	"github.com/cryplio/rest-api/src/modules/portfolios/http/handlers"
	"github.com/cryplio/rest-api/src/modules/portfolios/testportfolios"

	"github.com/Nivl/go-rest-tools/dependencies"
	"github.com/Nivl/go-rest-tools/network/http/httptests"
	"github.com/Nivl/go-rest-tools/testing/integration"
	"github.com/cryplio/rest-api/src/modules/api"
	"github.com/cryplio/rest-api/src/modules/users/testusers"
	"github.com/stretchr/testify/assert"
)

func TestIntegrationUpdate(t *testing.T) {
	t.Parallel()

	helper, err := integration.New(NewDeps(), migrationFolder)
	if err != nil {
		panic(err)
	}
	defer helper.Close()
	dbCon := helper.Deps.DB()

	u1, s1 := testusers.NewAuth(t, dbCon)
	toUpdate := testportfolios.NewPersistedPortfolio(t, dbCon, u1, nil)
	_, s2 := testusers.NewAuth(t, dbCon)

	tests := []struct {
		description string
		code        int
		params      *handlers.UpdateParams
		auth        *httptests.RequestAuth
	}{
		{
			"Not logged",
			http.StatusUnauthorized,
			&handlers.UpdateParams{ID: toUpdate.ID},
			nil,
		},
		{
			"Updating someone's portfolio",
			http.StatusNotFound,
			&handlers.UpdateParams{ID: toUpdate.ID, Name: "should not change"},
			httptests.NewRequestAuth(s2),
		},
		{
			"Updating portfolio",
			http.StatusOK,
			&handlers.UpdateParams{ID: toUpdate.ID, Name: "New name"},
			httptests.NewRequestAuth(s1),
		},
		{
			"Updating unexisting portolio",
			http.StatusNotFound,
			&handlers.UpdateParams{ID: "9cef9a62-f320-4efe-b52c-6038c1e4668c", Name: "New name"},
			httptests.NewRequestAuth(s1),
		},
	}

	t.Run("parallel", func(t *testing.T) {
		for _, tc := range tests {
			tc := tc
			t.Run(tc.description, func(t *testing.T) {
				t.Parallel()
				defer helper.RecoverPanic()

				rec := callUpdate(t, tc.params, tc.auth, helper.Deps)
				assert.Equal(t, tc.code, rec.Code)

				if rec.Code == http.StatusOK {
					var u portfolios.Payload
					if err := json.NewDecoder(rec.Body).Decode(&u); err != nil {
						t.Fatal(err)
					}

					if tc.params.Name != "" {
						assert.Equal(t, tc.params.Name, u.Name)
					}
				}
			})
		}
	})
}

func callUpdate(t *testing.T, params *handlers.UpdateParams, auth *httptests.RequestAuth, deps dependencies.Dependencies) *httptest.ResponseRecorder {
	ri := &httptests.RequestInfo{
		Endpoint: handlers.Endpoints[handlers.EndpointUpdate],
		Params:   params,
		Auth:     auth,
		Router:   api.GetRouter(deps),
	}
	return httptests.NewRequest(t, ri)
}
