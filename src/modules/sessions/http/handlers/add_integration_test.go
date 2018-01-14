// +build integration

package handlers_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Nivl/go-rest-tools/dependencies"
	"github.com/Nivl/go-rest-tools/network/http/httptests"
	"github.com/Nivl/go-rest-tools/security/auth/testauth"
	"github.com/Nivl/go-rest-tools/testing/integration"
	"github.com/cryplio/rest-api/src/modules/api"
	"github.com/cryplio/rest-api/src/modules/sessions"
	"github.com/cryplio/rest-api/src/modules/sessions/http/handlers"
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

	u1 := testauth.NewPersistedUser(t, dbCon, nil)

	tests := []struct {
		description string
		code        int
		params      *handlers.AddParams
	}{
		{
			"Unexisting email should fail",
			http.StatusBadRequest,
			&handlers.AddParams{Email: "invalid@fake.com", Password: "fake"},
		},
		{
			"Invalid password should fail",
			http.StatusBadRequest,
			&handlers.AddParams{Email: u1.Email, Password: "invalid"},
		},
		{
			"Valid Request should work",
			http.StatusCreated,
			&handlers.AddParams{Email: u1.Email, Password: "fake"},
		},
	}

	t.Run("parallel", func(t *testing.T) {
		for _, tc := range tests {
			tc := tc
			t.Run(tc.description, func(t *testing.T) {
				t.Parallel()
				defer helper.RecoverPanic()

				rec := callAdd(t, tc.params, helper.Deps)
				assert.Equal(t, tc.code, rec.Code)

				if rec.Code == http.StatusCreated {
					var session sessions.Payload
					if err := json.NewDecoder(rec.Body).Decode(&session); err != nil {
						t.Fatal(err)
					}

					assert.NotEmpty(t, session.Token)
					assert.Equal(t, u1.ID, session.UserID)
				}
			})
		}
	})
}

func callAdd(t *testing.T, params *handlers.AddParams, deps dependencies.Dependencies) *httptest.ResponseRecorder {
	ri := &httptests.RequestInfo{
		Endpoint: handlers.Endpoints[handlers.EndpointAdd],
		Params:   params,
		Router:   api.GetRouter(deps),
	}
	return httptests.NewRequest(t, ri)
}
