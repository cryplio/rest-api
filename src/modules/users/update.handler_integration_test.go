// +build integration

package users_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Nivl/go-hasher/implementations/bcrypt"
	"github.com/Nivl/go-rest-tools/dependencies"
	"github.com/Nivl/go-rest-tools/network/http/httptests"
	"github.com/Nivl/go-rest-tools/security/auth"
	"github.com/Nivl/go-rest-tools/testing/integration"
	"github.com/cryplio/rest-api/src/modules/api"
	"github.com/cryplio/rest-api/src/modules/users"
	"github.com/cryplio/rest-api/src/modules/users/testusers"
	"github.com/stretchr/testify/assert"
)

func TestUpdate(t *testing.T) {
	t.Parallel()

	helper, err := integration.New(NewDeps(), migrationFolder)
	if err != nil {
		panic(err)
	}
	defer helper.Close()
	dbCon := helper.Deps.DB()

	u1, s1 := testusers.NewAuth(t, dbCon)
	u2, s2 := testusers.NewAuth(t, dbCon)

	tests := []struct {
		description string
		code        int
		params      *users.UpdateParams
		auth        *httptests.RequestAuth
	}{
		{
			"Not logged",
			http.StatusUnauthorized,
			&users.UpdateParams{ID: u1.ID},
			nil,
		},
		{
			"Updating an other user",
			http.StatusForbidden,
			&users.UpdateParams{ID: u1.ID},
			httptests.NewRequestAuth(s2),
		},
		{
			"Updating email without providing password",
			http.StatusUnauthorized,
			&users.UpdateParams{ID: u1.ID, Email: "melvin@fake.io"},
			httptests.NewRequestAuth(s1),
		},
		{
			"Updating password without providing current Password",
			http.StatusUnauthorized,
			&users.UpdateParams{ID: u1.ID, NewPassword: "TestUpdateUser"},
			httptests.NewRequestAuth(s1),
		},
		{
			"Updating email to a used one",
			http.StatusConflict,
			&users.UpdateParams{ID: u1.ID, CurrentPassword: "fake", Email: u2.Email},
			httptests.NewRequestAuth(s1),
		},
		{
			"Updating password",
			http.StatusOK,
			&users.UpdateParams{ID: u2.ID, CurrentPassword: "fake", NewPassword: "TestUpdateUser"},
			httptests.NewRequestAuth(s2),
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
					var u users.ProfilePayload
					if err := json.NewDecoder(rec.Body).Decode(&u); err != nil {
						t.Fatal(err)
					}

					if tc.params.Name != "" {
						assert.Equal(t, tc.params.Name, u.Name)
					}
					if tc.params.Email != "" {
						assert.Equal(t, tc.params.Email, u.Email)
					}

					if tc.params.NewPassword != "" {
						// To check the password has been updated with need to get the
						// encrypted version, and compare it to the raw one
						updatedUser, err := auth.GetUserByID(dbCon, tc.params.ID)
						if err != nil {
							t.Fatal(err)
						}

						hashr := bcrypt.Bcrypt{}
						hash := updatedUser.Password
						assert.True(t, hashr.IsValid(hash, tc.params.NewPassword))
					}
				}
			})
		}
	})
}

func callUpdate(t *testing.T, params *users.UpdateParams, auth *httptests.RequestAuth, deps dependencies.Dependencies) *httptest.ResponseRecorder {
	ri := &httptests.RequestInfo{
		Endpoint: users.Endpoints[users.EndpointUpdate],
		Params:   params,
		Auth:     auth,
		Router:   api.GetRouter(deps),
	}
	return httptests.NewRequest(t, ri)
}
