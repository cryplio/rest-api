// +build integration

package handlers_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Nivl/go-rest-tools/dependencies"
	"github.com/Nivl/go-rest-tools/network/http/httptests"
	"github.com/Nivl/go-rest-tools/testing/integration"
	"github.com/cryplio/rest-api/src/modules/api"
	"github.com/cryplio/rest-api/src/modules/users"
	"github.com/cryplio/rest-api/src/modules/users/http/handlers"
	"github.com/cryplio/rest-api/src/modules/users/testusers"
	"github.com/stretchr/testify/assert"
)

func TestGet(t *testing.T) {
	t.Parallel()

	helper, err := integration.New(NewDeps(), migrationFolder)
	if err != nil {
		panic(err)
	}
	defer helper.Close()
	dbCon := helper.Deps.DB()

	u1, s1 := testusers.NewAuth(t, dbCon)
	_, s2 := testusers.NewAuth(t, dbCon)

	tests := []struct {
		description string
		code        int
		params      *handlers.GetParams
		auth        *httptests.RequestAuth
	}{
		{
			"Not logged",
			http.StatusOK,
			&handlers.GetParams{ID: u1.ID},
			nil,
		},
		{
			"Getting an other user",
			http.StatusOK,
			&handlers.GetParams{ID: u1.ID},
			httptests.NewRequestAuth(s2),
		},
		{
			"Getting own data",
			http.StatusOK,
			&handlers.GetParams{ID: u1.ID},
			httptests.NewRequestAuth(s1),
		},
		{
			"Getting un-existing user with valid ID",
			http.StatusNotFound,
			&handlers.GetParams{ID: "f76700e7-988c-4ae9-9f02-ac3f9d7cd88e"},
			nil,
		},
		{
			"Getting un-existing user with invalid ID",
			http.StatusBadRequest,
			&handlers.GetParams{ID: "invalidID"},
			nil,
		},
	}

	t.Run("parallel", func(t *testing.T) {
		for _, tc := range tests {
			tc := tc
			t.Run(tc.description, func(t *testing.T) {
				t.Parallel()
				defer helper.RecoverPanic()

				rec := callGet(t, tc.params, tc.auth, helper.Deps)
				assert.Equal(t, tc.code, rec.Code)

				if rec.Code == http.StatusOK {
					var u users.ProfilePayload
					if err := json.NewDecoder(rec.Body).Decode(&u); err != nil {
						t.Fatal(err)
					}

					if assert.Equal(t, tc.params.ID, u.ID, "Not the same user") {
						// User access their own data
						if tc.auth != nil && u.ID == tc.auth.UserID {
							assert.NotEmpty(t, u.Email, "Same user needs their private data")
						} else { // user access an other user data
							assert.Empty(t, u.Email, "Should not return private data")
						}
					}
				}
			})
		}
	})
}

func callGet(t *testing.T, params *handlers.GetParams, auth *httptests.RequestAuth, deps dependencies.Dependencies) *httptest.ResponseRecorder {
	ri := &httptests.RequestInfo{
		Endpoint: handlers.Endpoints[handlers.EndpointGet],
		Params:   params,
		Auth:     auth,
		Router:   api.GetRouter(deps),
	}
	return httptests.NewRequest(t, ri)
}
