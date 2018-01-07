package testguard

import (
	"testing"

	"github.com/Nivl/go-rest-tools/router/guard"
	"github.com/Nivl/go-rest-tools/security/auth"
	"github.com/stretchr/testify/assert"
)

// AccessTestCase represents the data needed to test params
type AccessTestCase struct {
	Description string
	User        *auth.User
	ErrCode     int // <= 0 for no error
}

// AccessTest checks if the auth are correctly set
func AccessTest(t *testing.T, g *guard.Guard, testCases []AccessTestCase) {
	for _, tc := range testCases {
		tc := tc
		t.Run(tc.Description, func(t *testing.T) {
			t.Parallel()

			_, err := g.HasAccess(tc.User)
			if tc.ErrCode > 0 {
				assert.Error(t, err)
				assert.Equal(t, tc.ErrCode, err.HTTPStatus())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
