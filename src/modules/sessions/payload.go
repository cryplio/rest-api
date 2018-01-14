package sessions

import "github.com/Nivl/go-rest-tools/security/auth"

// Payload represents a Session that can be safely returned by the API
type Payload struct {
	Token  string `json:"token"`
	UserID string `json:"user_id"`
}

// NewPayload turns a Session into an object that is safe to be
// returned by the API
func NewPayload(s *auth.Session) *Payload {
	return &Payload{
		Token:  s.ID,
		UserID: s.UserID,
	}
}
