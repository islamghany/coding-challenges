package session

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"strings"

	"smtp/command"
	"smtp/email"
)

// State represents the current state of the SMTP session
type State int

const (
	StateConnected State = iota // Just connected, waiting for HELO
	StateGreeted                // HELO received, waiting for MAIL FROM
	StateMail                   // MAIL FROM received, waiting for RCPT TO
	StateRcpt                   // RCPT TO received, waiting for DATA
	StateData                   // DATA received, reading message
	StateQuit                   // QUIT received, session ending
)

// Session represents an SMTP session with a client
type Session struct {
	conn   net.Conn
	reader *bufio.Reader
	state  State
	store  email.Store

	// Envelope data
	mailFrom string
	rcptTo   []string
}

// New creates a new SMTP session
func New(conn net.Conn, store email.Store) *Session {
	return &Session{
		conn:   conn,
		reader: bufio.NewReader(conn),
		state:  StateConnected,
		store:  store,
	}
}

// Run starts the session's main loop
func (s *Session) Run() error {
	for s.state != StateQuit {
		line, err := s.reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				return nil // Client disconnected gracefully
			}
			return err
		}

		line = strings.TrimRight(line, "\r\n")

		if err := s.handleLine(line); err != nil {
			return err
		}
	}
	return nil
}

func (s *Session) handleLine(line string) error {
	cmd, err := command.Parse(line)
	if err != nil {
		return s.reply(500, "Syntax error: "+err.Error())
	}

	switch cmd.Name {
	case command.HELO, command.EHLO:
		return s.handleHelo(cmd)
	case command.MAIL:
		return s.handleMailFrom(cmd)
	case command.RCPT:
		return s.handleRcptTo(cmd)
	case command.DATA:
		return s.handleData()
	case command.QUIT:
		return s.handleQuit()
	case command.RSET:
		return s.handleRset()
	case command.NOOP:
		return s.reply(250, "OK")
	default:
		return s.reply(502, "Command not implemented")
	}
}

// reply sends an SMTP response to the client
func (s *Session) reply(code int, msg string) error {
	_, err := fmt.Fprintf(s.conn, "%d %s\r\n", code, msg)
	return err
}

// resetEnvelope clears the current mail transaction
func (s *Session) resetEnvelope() {
	s.mailFrom = ""
	s.rcptTo = nil
}
