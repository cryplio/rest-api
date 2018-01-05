package sentryreporter

import "github.com/Nivl/go-reporter"

// Makes sure Creator is a Creator
var _ reporter.Creator = (*Creator)(nil)

// Creator is a Creator use to get sentry clients
type Creator struct {
	con string
}

// New returns a new sentry client
func (c *Creator) New() (reporter.Reporter, error) {
	return New(c.con)
}

// NewCreator creates a client's creator
func NewCreator(con string) (*Creator, error) {
	return &Creator{
		con: con,
	}, nil
}
