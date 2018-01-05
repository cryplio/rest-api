package noopreporter

import "github.com/Nivl/go-reporter"

// Makes sure Creator is a Creator
var _ reporter.Creator = (*Creator)(nil)

// NewCreator creates a client's creator
func NewCreator() (*Creator, error) {
	return &Creator{}, nil
}

// Creator is a Creator used to get no-op reporter
type Creator struct{}

// New returns a new mailer Reporter
func (c *Creator) New() (reporter.Reporter, error) {
	return New()
}
