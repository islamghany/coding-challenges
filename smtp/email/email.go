package email

import (
	"fmt"
	"time"
)

// Email represents a received email message
type Email struct {
	ID       string
	From     string
	To       []string
	Content  string
	Received time.Time
}

// New creates a new Email with auto-generated ID
func New(from string, to []string, content string) *Email {
	return &Email{
		ID:       generateID(),
		From:     from,
		To:       to,
		Content:  content,
		Received: time.Now(),
	}
}

func generateID() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}

// String returns a formatted string representation of the email
func (e *Email) String() string {
	return fmt.Sprintf("Email{ID: %s, From: %s, To: %v, Received: %s}",
		e.ID, e.From, e.To, e.Received.Format(time.RFC3339))
}
