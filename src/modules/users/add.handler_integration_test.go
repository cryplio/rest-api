// +build integration

package users_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/cryplio/rest-api/src/modules/users/testusers"

	"github.com/Nivl/go-rest-tools/dependencies"
	"github.com/Nivl/go-rest-tools/network/http/httptests"
	"github.com/Nivl/go-rest-tools/testing/integration"
	"github.com/cryplio/rest-api/src/modules/api"
	"github.com/cryplio/rest-api/src/modules/users"
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

	existingUser := testusers.NewPersistedProfile(t, dbCon, nil)

	tests := []struct {
		description string
		code        int
		params      *users.AddParams
	}{
		{
			"It should fail to add an empty user",
			http.StatusBadRequest,
			&users.AddParams{},
		},
		{
			"It should add a valid user",
			http.StatusCreated,
			&users.AddParams{Name: "Name", Email: "email+TestAdd@fake.com", Password: "password"},
		},
		{
			"It should fail adding a user with an email already taken",
			http.StatusConflict,
			&users.AddParams{Name: "Name", Email: existingUser.User.Email, Password: "password"},
		},
	}

	t.Run("parallel", func(t *testing.T) {
		for _, tc := range tests {
			tc := tc
			t.Run(tc.description, func(t *testing.T) {
				t.Parallel()
				defer helper.RecoverPanic()

				rec := callHandlerAdd(t, tc.params, helper.Deps)
				assert.Equal(t, tc.code, rec.Code)

				if rec.Code == http.StatusCreated {
					var u users.ProfilePayload
					if err := json.NewDecoder(rec.Body).Decode(&u); err != nil {
						t.Fatal(err)
					}

					assert.NotEmpty(t, u.ID)
					assert.Equal(t, tc.params.Email, u.Email)

					// Let's make sure we have a portfolio created
					var totalPortfolios int
					stmt := `SELECT COUNT(*) from user_portfolios WHERE user_id=$1`
					err := dbCon.NamedGet(&totalPortfolios, stmt, u.ID)
					require.NoError(t, err, "could not get the number of portfolios")
					assert.Equal(t, 1, totalPortfolios, "1 portfolio should have been created")
				}
			})
		}
	})
}

func callHandlerAdd(t *testing.T, params *users.AddParams, deps dependencies.Dependencies) *httptest.ResponseRecorder {
	ri := &httptests.RequestInfo{
		Endpoint: users.Endpoints[users.EndpointAdd],
		Params:   params,
		Router:   api.GetRouter(deps),
	}
	return httptests.NewRequest(t, ri)
}
