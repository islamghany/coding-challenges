package session

import (
	"strings"

	"smtp/command"
	"smtp/email"
)

// handleHelo handles HELO and EHLO commands
func (s *Session) handleHelo(cmd *command.Command) error {
	if len(cmd.Args) < 1 {
		return s.reply(501, "Missing domain argument")
	}

	s.state = StateGreeted
	s.resetEnvelope()

	return s.reply(250, "Hello "+cmd.Args[0])
}

// handleMailFrom handles the MAIL FROM command
func (s *Session) handleMailFrom(cmd *command.Command) error {
	if s.state != StateGreeted {
		return s.reply(503, "Bad sequence of commands")
	}

	if len(cmd.Args) < 1 {
		return s.reply(501, "Missing sender address")
	}

	s.mailFrom = cmd.Args[0]
	s.state = StateMail

	return s.reply(250, "OK")
}

// handleRcptTo handles the RCPT TO command
func (s *Session) handleRcptTo(cmd *command.Command) error {
	if s.state != StateMail && s.state != StateRcpt {
		return s.reply(503, "Bad sequence of commands")
	}

	if len(cmd.Args) < 1 {
		return s.reply(501, "Missing recipient address")
	}

	s.rcptTo = append(s.rcptTo, cmd.Args[0])
	s.state = StateRcpt

	return s.reply(250, "OK")
}

// handleData handles the DATA command and reads the message body
func (s *Session) handleData() error {
	if s.state != StateRcpt {
		return s.reply(503, "Bad sequence of commands")
	}

	// Send intermediate response
	if err := s.reply(354, "Start mail input; end with <CRLF>.<CRLF>"); err != nil {
		return err
	}

	s.state = StateData

	// Read message body
	var body strings.Builder
	for {
		line, err := s.reader.ReadString('\n')
		if err != nil {
			return err
		}
		line = strings.TrimRight(line, "\r\n")

		// Check for end of message
		if line == "." {
			break
		}

		// Dot-unstuffing: if line starts with "..", remove one dot
		if strings.HasPrefix(line, ".") {
			line = line[1:]
		}

		body.WriteString(line + "\n")
	}

	// Create and store the email
	msg := email.New(s.mailFrom, s.rcptTo, body.String())
	if err := s.store.Save(msg); err != nil {
		s.state = StateGreeted
		s.resetEnvelope()
		return s.reply(451, "Failed to store message")
	}

	// Reset for next transaction
	s.state = StateGreeted
	s.resetEnvelope()

	return s.reply(250, "OK: message queued as "+msg.ID)
}

// handleQuit handles the QUIT command
func (s *Session) handleQuit() error {
	s.state = StateQuit
	return s.reply(221, "Bye")
}

// handleRset handles the RSET command
func (s *Session) handleRset() error {
	if s.state >= StateGreeted {
		s.state = StateGreeted
	}
	s.resetEnvelope()
	return s.reply(250, "OK")
}
