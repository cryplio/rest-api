// Package noopreporter is a no-op implementation of a Reporter
package noopreporter

import reporter "github.com/Nivl/go-reporter"

var (
	// Makes sure the Email object implement Reporter
	_ reporter.Reporter = (*Reporter)(nil)
)

// New creates a new Mailer reporter
func New() (*Reporter, error) {
	return &Reporter{}, nil
}

// Reporter represents reporter that does nothin
type Reporter struct {
}

// SetUser does nothing
func (r *Reporter) SetUser(u *reporter.User) {}

// AddTag does nothing
func (r *Reporter) AddTag(key, value string) {}

// AddTags does nothing
func (r *Reporter) AddTags(tags map[string]string) {}

// ReportError does nothing
func (r *Reporter) ReportError(err error) {}

// ReportErrorAndWait does nothing
func (r *Reporter) ReportErrorAndWait(err error) {}
