// Package sendgridmailer contains an implementation of the mailer interface that
// uses sendgrid
package sendgridmailer

import (
	"errors"
	"fmt"
	"strings"

	mailer "github.com/Nivl/go-mailer"
	sendgrid "github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

// Makes sure Sendgrid implements Mailer
var _ mailer.Mailer = (*Sendgrid)(nil)

// Sendgrid is an object used to send email through sendgrid
type Sendgrid struct {
	APIKey               string
	DefaultFrom          string
	DefaultTo            string
	StacktraceTemplateID string
}

// SendStackTrace emails the current stacktrace to the default FROM
func (s *Sendgrid) SendStackTrace(trace []byte, message string, context map[string]string) error {
	if s.StacktraceTemplateID == "" {
		return errors.New("StacktraceTemplateID not set")
	}

	msg := mailer.NewMessage(s.StacktraceTemplateID)
	stacktrace := string(trace[:])

	msg.Body = strings.Replace(stacktrace, "\n", "<br>", -1)

	htmlContext := fmt.Sprintf("<li><strong>Error</strong>: %s</li>\n", message)
	for k, v := range context {
		htmlContext += fmt.Sprintf("<li><strong>%s</strong>: %s</li>\n", k, v)
	}
	msg.SetVar("context", htmlContext)
	return s.Send(msg)
}

// Send is used to send an email
func (s *Sendgrid) Send(msg *mailer.Message) error {
	from := mail.NewEmail("No Reply", msg.From)
	if msg.From == "" {
		from = mail.NewEmail("No Reply", s.DefaultFrom)
	}

	to := mail.NewEmail(msg.To, msg.To)
	if msg.From == "" {
		to = mail.NewEmail(s.DefaultTo, s.DefaultTo)
	}

	content := mail.NewContent("text/html", msg.Body)
	email := mail.NewV3MailInit(from, msg.Subject, to, content)
	email.SetTemplateID(msg.TemplateID)

	for k, v := range msg.Vars {
		email.Personalizations[0].SetSubstitution(k, v)
	}

	request := sendgrid.GetRequest(s.APIKey, "/v3/mail/send", "https://api.sendgrid.com")
	request.Method = "POST"
	_, err := sendgrid.API(request)

	return err
}

// New creates and returns a new mailer using sendgrid
func New(APIKey, defaultFrom, defaultTo, stacktraceUUID string) *Sendgrid {
	return &Sendgrid{
		APIKey:               APIKey,
		DefaultFrom:          defaultFrom,
		DefaultTo:            defaultTo,
		StacktraceTemplateID: stacktraceUUID,
	}
}
