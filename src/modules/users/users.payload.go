package users

import "github.com/Nivl/go-rest-tools/security/auth"

// Payload represents a user payload with non public field
type Payload struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Email   string `json:"email,omitempty"`
	IsAdmin bool   `json:"is_admin,omitempty"`
}

// NewPrivatePayload turns a user into an object that is safe to be
// returned by the API
func NewPrivatePayload(u *auth.User) *Payload {
	pld := NewPayload(u)
	pld.Email = u.Email
	return pld
}

// NewPayload turns a user into an object that is safe to be
// returned by the API
func NewPayload(u *auth.User) *Payload {
	return &Payload{
		ID:      u.ID,
		Name:    u.Name,
		IsAdmin: u.IsAdmin,
	}
}
