package mailerreporter

import (
	"github.com/Nivl/go-mailer"
	"github.com/Nivl/go-reporter"
)

// Makes sure Creator is a Creator
var _ reporter.Creator = (*Creator)(nil)

// NewCreator creates a client's creator
func NewCreator(m mailer.Mailer) (*Creator, error) {
	return &Creator{
		mailer: m,
	}, nil
}

// Creator is a Creator used to get mailers
type Creator struct {
	mailer mailer.Mailer
}

// New returns a new mailer Reporter
func (c *Creator) New() (reporter.Reporter, error) {
	return New(c.mailer)
}
