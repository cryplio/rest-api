// +build integration

package portfolios_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/cryplio/rest-api/src/modules/portfolios"

	"github.com/Nivl/go-rest-tools/dependencies"
	"github.com/Nivl/go-rest-tools/network/http/httptests"
	"github.com/Nivl/go-rest-tools/security/auth/testauth"
	"github.com/Nivl/go-rest-tools/testing/integration"
	"github.com/cryplio/rest-api/src/modules/api"
	"github.com/stretchr/testify/assert"
)

func TestAdd(t *testing.T) {
	t.Parallel()

	helper, err := integration.New(NewDeps(), migrationFolder)
	if err != nil {
		panic(err)
	}
	defer helper.Close()
	dbCon := helper.Deps.DB()

	_, userSession := testauth.NewPersistedAuth(t, dbCon)

	tests := []struct {
		description string
		code        int
		params      *portfolios.AddParams
		auth        *httptests.RequestAuth
	}{
		{
			"Not logged",
			http.StatusUnauthorized,
			&portfolios.AddParams{Name: "My portfolio"},
			nil,
		},
		{
			"It should add a new portfolio",
			http.StatusCreated,
			&portfolios.AddParams{Name: "My portfolio"},
			httptests.NewRequestAuth(userSession),
		},
	}

	t.Run("parallel", func(t *testing.T) {
		for _, tc := range tests {
			tc := tc
			t.Run(tc.description, func(t *testing.T) {
				t.Parallel()
				defer helper.RecoverPanic()

				rec := callHandlerAdd(t, tc.params, tc.auth, helper.Deps)
				assert.Equal(t, tc.code, rec.Code)

				if rec.Code == http.StatusCreated {
					var p portfolios.Payload
					if err := json.NewDecoder(rec.Body).Decode(&p); err != nil {
						t.Fatal(err)
					}

					assert.NotEmpty(t, p.ID)
					assert.Equal(t, tc.params.Name, p.Name)
				}
			})
		}
	})
}

func callHandlerAdd(t *testing.T, params *portfolios.AddParams, auth *httptests.RequestAuth, deps dependencies.Dependencies) *httptest.ResponseRecorder {
	ri := &httptests.RequestInfo{
		Endpoint: portfolios.Endpoints[portfolios.EndpointAdd],
		Params:   params,
		Auth:     auth,
		Router:   api.GetRouter(deps),
	}
	return httptests.NewRequest(t, ri)
}
