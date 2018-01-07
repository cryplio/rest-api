// Package bcrypt is an implementation of the hasher interface for bcrypt.
package bcrypt

import (
	hasher "github.com/Nivl/go-hasher"
	"golang.org/x/crypto/bcrypt"
)

var _ hasher.Hasher = (*Bcrypt)(nil)

// Bcrypt is an Hasher implementation for bcrypt
type Bcrypt struct {
}

// Hash returns a hash for the provided string
func (h Bcrypt) Hash(raw string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(raw), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

// IsValid checks if a hash and a string match
func (h Bcrypt) IsValid(hash string, raw string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(raw))
	return err == nil
}
