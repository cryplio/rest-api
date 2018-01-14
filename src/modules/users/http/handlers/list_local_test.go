package handlers

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestListParamsGetSort(t *testing.T) {
	// sugar
	shouldFail := true

	testCases := []struct {
		description string
		fields      string
		expected    string
		shouldFail  bool
	}{
		{
			"No fields should return the default sorting",
			"",
			"name ASC,created_at ASC",
			!shouldFail,
		},
		{
			"Order by ,, should return the default sorting",
			",,",
			"name ASC,created_at ASC",
			!shouldFail,
		},
		{
			"Order by ,,,,,,, should return the default sorting",
			",,,,,,,",
			"name ASC,created_at ASC",
			!shouldFail,
		},
		{
			"Order by ,,,name,,,, should sort by name",
			",,,name,,,,",
			"name ASC",
			!shouldFail,
		},
		{
			"Order by -name should work",
			"-name",
			"name DESC",
			!shouldFail,
		},
		{
			"Order by created_at and -name should work",
			"created_at,-name",
			"created_at ASC,name DESC",
			!shouldFail,
		},
		{
			"Order by not_a_field should fail",
			"not_a_field",
			"",
			shouldFail,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.description, func(t *testing.T) {
			t.Parallel()

			output, err := listParamsGetSort(tc.fields)
			if tc.shouldFail {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expected, output)
			}
		})
	}
}
