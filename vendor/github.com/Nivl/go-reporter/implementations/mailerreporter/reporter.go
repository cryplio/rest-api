// Package mailerreporter is an implementation of a Reporter that sends emails
package mailerreporter

import (
	"runtime/debug"

	mailer "github.com/Nivl/go-mailer"
	reporter "github.com/Nivl/go-reporter"
)

var (
	// Makes sure the Email object implement Reporter
	_ reporter.Reporter = (*Reporter)(nil)
)

// New creates a new Mailer reporter
func New(m mailer.Mailer) (*Reporter, error) {
	return &Reporter{
		client: m,
		tags:   map[string]string{},
	}, nil
}

// Reporter represents a client used to report errors by email
type Reporter struct {
	client mailer.Mailer
	user   *reporter.User
	tags   map[string]string
}

// SetUser attaches the provided user to report
func (r *Reporter) SetUser(u *reporter.User) {
	r.user = u
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

// ReportError sends a report with this specific error
func (r *Reporter) ReportError(err error) {
	// We copy the tags so we can add more data without affecting the source
	finalTags := map[string]string{}
	if r.user != nil {
		finalTags["User ID"] = r.user.ID
		finalTags["User email"] = r.user.Email
		finalTags["User name"] = r.user.Username
	}
	for k, v := range r.tags {
		finalTags[k] = v
	}

	sendEmail := func(stacktrace []byte) {
		r.client.SendStackTrace(stacktrace, err.Error(), finalTags)
	}
	go sendEmail(debug.Stack())
}

// ReportErrorAndWait sends a report with this specific error
func (r *Reporter) ReportErrorAndWait(err error) {
	// We copy the tags so we can add more data without affecting the source
	finalTags := map[string]string{}
	if r.user != nil {
		finalTags["User ID"] = r.user.ID
		finalTags["User email"] = r.user.Email
		finalTags["User name"] = r.user.Username
	}
	for k, v := range r.tags {
		finalTags[k] = v
	}
	r.client.SendStackTrace(debug.Stack(), err.Error(), finalTags)
}
