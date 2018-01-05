package mailer

import "fmt"

// Message represents a message to send
type Message struct {
	TemplateID string
	From       string
	To         string
	Subject    string
	Body       string
	Vars       map[string]string
}

// NewMessage creates a new message from a template
func NewMessage(templateID string) *Message {
	return &Message{
		TemplateID: templateID,
		Vars:       map[string]string{},
	}
}

// SetVar set a variable to the message
func (msg *Message) SetVar(name string, value string) {
	fullName := fmt.Sprintf("-%s-", name)
	msg.Vars[fullName] = value
}
