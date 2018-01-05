// Package printmailer contains an implementation of the mailer interface that
// prints
package printmailer

import (
	"fmt"

	mailer "github.com/Nivl/go-mailer"
)

// Makes sure Mailer implements mailer.Mailer
var _ mailer.Mailer = (*Mailer)(nil)

// Mailer is a mailer that just print emails
type Mailer struct {
}

// SendStackTrace emails the current stacktrace to the default FROM
func (m *Mailer) SendStackTrace(trace []byte, message string, context map[string]string) error {
	fmt.Printf("%s,%#v\n%s", message, context, trace)
	return nil
}

// Send is used to send an email
func (m *Mailer) Send(msg *mailer.Message) error {
	fmt.Printf("FROM: %s\nTO: %s\nSUBJECT: %s\n%s\n", msg.From, msg.To, msg.Subject, msg.Body)
	return nil
}
