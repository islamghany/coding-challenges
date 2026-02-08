package command

import (
	"fmt"
	"strings"
)

// Parse parses an SMTP command line into a Command struct
func Parse(line string) (*Command, error) {
	line = strings.TrimSpace(line)
	if line == "" {
		return nil, fmt.Errorf("empty command")
	}

	fields := strings.Fields(line)
	cmdName := strings.ToUpper(fields[0])

	// Handle two-word commands: MAIL FROM, RCPT TO
	switch cmdName {
	case "MAIL":
		return parseMailFrom(line)
	case "RCPT":
		return parseRcptTo(line)
	default:
		return &Command{
			Name: Name(cmdName),
			Args: fields[1:],
		}, nil
	}
}

func parseMailFrom(line string) (*Command, error) {
	parts := strings.SplitN(line, ":", 2)
	if len(parts) != 2 || strings.ToUpper(strings.TrimSpace(parts[0])) != "MAIL FROM" {
		return nil, fmt.Errorf("invalid MAIL FROM syntax")
	}
	return &Command{
		Name: MAIL,
		Args: []string{strings.TrimSpace(parts[1])},
	}, nil
}

func parseRcptTo(line string) (*Command, error) {
	parts := strings.SplitN(line, ":", 2)
	if len(parts) != 2 || strings.ToUpper(strings.TrimSpace(parts[0])) != "RCPT TO" {
		return nil, fmt.Errorf("invalid RCPT TO syntax")
	}
	return &Command{
		Name: RCPT,
		Args: []string{strings.TrimSpace(parts[1])},
	}, nil
}
