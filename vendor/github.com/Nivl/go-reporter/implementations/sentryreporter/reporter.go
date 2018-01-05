// Package sentryreporter is a Reporter implementation using Sentry
package sentryreporter

import (
	reporter "github.com/Nivl/go-reporter"
	"github.com/getsentry/raven-go"
)

// Makes sure the Reporter object implement Reporter
var _ reporter.Reporter = (*Reporter)(nil)

// New creates a new sentry client
func New(con string) (*Reporter, error) {
	c, err := raven.New(con)
	if err != nil {
		return nil, err
	}

	return &Reporter{
		client: c,
		tags:   map[string]string{},
	}, nil
}

// Reporter represents a client used to report errors to Reporter
type Reporter struct {
	client *raven.Client
	tags   map[string]string
}

// SetUser attaches the provided user to report
func (r *Reporter) SetUser(u *reporter.User) {
	if u == nil {
		return
	}
	sentryUser := &raven.User{
		ID:       u.ID,
		Email:    u.Email,
		Username: u.Username,
	}
	r.client.SetUserContext(sentryUser)
}

// AddTag attaches the provided data to the report
func (r *Reporter) AddTag(key, value string) {
	r.tags[key] = value
}

// AddTags attaches the provided data to the report
func (r *Reporter) AddTags(tags map[string]string) {
	for key, value := range tags {
		r.tags[key] = value
	}
}

// ReportError sends a report with the specified error
func (r *Reporter) ReportError(err error) {
	r.client.CaptureError(err, r.tags)
}

// ReportErrorAndWait sends a report with the specified error, and wait for
// the reports to be sent
func (r *Reporter) ReportErrorAndWait(err error) {
	r.client.CaptureErrorAndWait(err, r.tags)
}
