// email/sender.go
package email

import (
	"fmt"
	"net/smtp"
	"strings"
)

type SMTPConfig struct {
	Host     string
	Port     string
	Username string
	Password string
}

// GmailConfig returns config for Gmail SMTP
func GmailConfig(email, appPassword string) *SMTPConfig {
	return &SMTPConfig{
		Host:     "smtp.gmail.com",
		Port:     "587",
		Username: email,
		Password: appPassword,
	}
}

// Send sends an email via SMTP
func Send(cfg *SMTPConfig, from string, to []string, subject, body string) error {
	auth := smtp.PlainAuth("", cfg.Username, cfg.Password, cfg.Host)

	// Build message
	msg := strings.Builder{}
	msg.WriteString(fmt.Sprintf("From: %s\r\n", from))
	msg.WriteString(fmt.Sprintf("To: %s\r\n", strings.Join(to, ", ")))
	msg.WriteString(fmt.Sprintf("Subject: %s\r\n", subject))
	msg.WriteString("MIME-Version: 1.0\r\n")
	msg.WriteString("Content-Type: text/plain; charset=utf-8\r\n")
	msg.WriteString("\r\n")
	msg.WriteString(body)

	addr := fmt.Sprintf("%s:%s", cfg.Host, cfg.Port)
	return smtp.SendMail(addr, auth, from, to, []byte(msg.String()))
}

// USAGE

// cfg := email.GmailConfig("your@gmail.com", "xxxx-xxxx-xxxx-xxxx")
// err := email.Send(cfg,
//     "your@gmail.com",
//     []string{"recipient@example.com"},
//     "Hello!",
//     "This is a test email from Go!",
// )

// HOW TO GET THE APP PASSWORD FOR GMAIL
// Step-by-Step Guide
// Step 1: Enable 2-Step Verification (if not already)
// Go to Google Account Security
// Under "How you sign in to Google", click 2-Step Verification
// Follow the prompts to enable it

// Step 2: Generate App Password
// Go to: https://myaccount.google.com/apppasswords
// Or navigate: Google Account → Security → 2-Step Verification → App Passwords
// You'll see a page like this:
//    App passwords      App passwords let you sign in to your Google Account from apps    that don't support 2-Step Verification.      App name: [________________]      [Create]
// Enter an app name (e.g., "Go SMTP Client" or "My App")
// Click Create
// Google will show you a 16-character password:
//    Your app password for "Go SMTP Client"      xxxx xxxx xxxx xxxx      Copy this password. You won't see it again.
// Copy it immediately! (You won't see it again)
