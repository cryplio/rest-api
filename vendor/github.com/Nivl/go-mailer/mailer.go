// Package mailer contains interfaces to deal with mailers
package mailer

//go:generate mockgen -destination implementations/mockmailer/mailer.go -package mockmailer github.com/Nivl/go-mailer Mailer

// Mailer is an object used to send email
type Mailer interface {
	// SendStackTrace emails the current stacktrace to the default FROM
	SendStackTrace(trace []byte, message string, context map[string]string) error

	// Send is used to send an email
	Send(msg *Message) error
}
