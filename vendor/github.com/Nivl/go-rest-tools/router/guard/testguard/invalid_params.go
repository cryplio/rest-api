package testguard

import (
	"net/url"
	"strings"
	"testing"

	"github.com/Nivl/go-rest-tools/types/apierror"
	"github.com/Nivl/go-params/formfile"
	"github.com/Nivl/go-rest-tools/router/guard"
	"github.com/stretchr/testify/assert"
)

// InvalidParamsTestCase represents the data needed to test params
type InvalidParamsTestCase struct {
	Description string
	MsgMatch    string
	FieldName   string
	Sources     map[string]url.Values
	FileHolder  formfile.FileHolder
}

// InvalidParams checks if the params are correctly failing
func InvalidParams(t *testing.T, g *guard.Guard, testCases []InvalidParamsTestCase) {
	for _, tc := range testCases {
		tc := tc
		t.Run(tc.Description, func(t *testing.T) {
			t.Parallel()

			_, err := g.ParseParams(tc.Sources, tc.FileHolder)
			if assert.Error(t, err, "expected the guard to fail") {
				assert.True(t, strings.Contains(err.Error(), tc.MsgMatch),
					"the error \"%s\" should contain the string \"%s\"", err.Error(), tc.MsgMatch)

				e, casted := err.(apierror.Error)
				if assert.True(t, casted, "expected the error to be a apierror.Error") {
					assert.Equal(t, tc.FieldName, e.Field(), "expected %s to be the failing param", tc.FieldName)
				}
			}
		})
	}
}
