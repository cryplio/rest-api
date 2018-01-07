package httptests

import (
	"encoding/base64"
	"fmt"

	"github.com/Nivl/go-rest-tools/security/auth"
)

// RequestAuth represents the auth data for a request
type RequestAuth struct {
	SessionID string
	UserID    string
}

// ToBasicAuth returns the data using the basic auth format
func (ra *RequestAuth) ToBasicAuth() string {
	authValue := fmt.Sprintf("%s:%s", ra.UserID, ra.SessionID)
	encoded := base64.StdEncoding.EncodeToString([]byte(authValue))
	return "basic " + encoded
}

// NewRequestAuth creates a new request auth
func NewRequestAuth(s *auth.Session) *RequestAuth {
	return &RequestAuth{
		SessionID: s.ID,
		UserID:    s.UserID,
	}
}
