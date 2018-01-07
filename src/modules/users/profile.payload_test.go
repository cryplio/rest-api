package users_test

import (
	"testing"

	"github.com/cryplio/rest-api/src/modules/users/testusers"

	"github.com/cryplio/rest-api/src/modules/users"
	"github.com/stretchr/testify/assert"
)

func TestExportPublicProfileNil(t *testing.T) {
	var profile *users.Profile
	pld := profile.ExportPublic()
	assert.Nil(t, pld, "nil should export as nil")
}

func TestExportPrivateProfileNil(t *testing.T) {
	var profile *users.Profile
	pld := profile.ExportPrivate()
	assert.Nil(t, pld, "nil should export as nil")
}

func TestExportPublicProfiles(t *testing.T) {
	profiles := users.Profiles{
		testusers.NewProfile(),
		testusers.NewProfile(),
		testusers.NewProfile(),
	}
	pld := profiles.ExportPublic()
	assert.Equal(t, len(profiles), len(pld.Results), "Wong number of profile returned")

	for _, p := range pld.Results {
		assert.Empty(t, p.Email, "Email should not have been set")
	}
}

func TestExportPrivateProfiles(t *testing.T) {
	profiles := users.Profiles{
		testusers.NewProfile(),
		testusers.NewProfile(),
		testusers.NewProfile(),
	}
	pld := profiles.ExportPrivate()
	assert.Equal(t, len(profiles), len(pld.Results), "Wong number of profile returned")

	for i, p := range pld.Results {
		assert.Equal(t, profiles[i].User.Email, p.Email, "Email should have been set")
	}
}
