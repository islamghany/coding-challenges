package parser

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"strconv"
)

type CmdName string

const (
	PUB     CmdName = "PUB"
	SUB     CmdName = "SUB"
	UNSUB   CmdName = "UNSUB"
	MSG     CmdName = "MSG"
	PONG    CmdName = "PONG"
	PING    CmdName = "PING"
	INFO    CmdName = "INFO"
	CONNECT CmdName = "CONNECT"
)

type Cmd struct {
	Name        CmdName
	Bytes       []byte
	Subject     []byte
	ID          []byte
	ConnectData ConnectCommand
}

type ConnectCommand struct {
	Name      string `json:"name,omitempty"`
	Protocol  int    `json:"protocol,omitempty"`
	AuthToken string `json:"auth_token,omitempty"`
	Username  string `json:"username,omitempty"`
	Password  string `json:"password,omitempty"`
	Verbose   bool   `json:"verbose,omitempty"`
}

// parse takes a io.Reader and returns a cmd, and error
func Parse(reader io.Reader) (*Cmd, error) {
	buffReader := bufio.NewReader(reader)
	// Read until newline without allocation
	line, err := buffReader.ReadSlice('\n')
	if err != nil {
		return nil, err
	}

	// Find the first space to separate the command name
	spaceIndex := bytes.IndexByte(line, ' ')

	if spaceIndex < 0 {
		return nil, fmt.Errorf("Empty Command")
	}

	cmdName := line[:spaceIndex]
	cmd := &Cmd{Name: CmdName(cmdName)} // No allocation on command name

	// Switch over the command name as a []byte to avoid string conversion
	switch {
	case bytes.Equal(cmdName, []byte("PUB")):
		return cmd.parsePUB(line[spaceIndex+1:], buffReader) // Pass rest of line
	case bytes.Equal(cmdName, []byte("CONNECT")):
		return cmd.parseCONNECT(line[spaceIndex+1:]) // Pass JSON part only
	default:
		return nil, fmt.Errorf("Unknown Command: %s", cmdName)
	}
}

func (c *Cmd) parsePUB(fields []byte, reader *bufio.Reader) (*Cmd, error) {
	parts := bytes.Fields(fields) // Parses without creating strings
	if len(parts) < 2 {
		return nil, fmt.Errorf("Insufficient arguments for PUB")
	}
	c.Subject = parts[0]
	bytesLength := -1
	var err error
	if len(parts) == 3 {
		c.ID = parts[1]
		bytesLength, err = strconv.Atoi(string(parts[2]))

		if err != nil {
			return nil, fmt.Errorf("Error parsing bytes length: %w", err)
		}
	} else {
		bytesLength, err = strconv.Atoi(string(parts[1]))
		if err != nil {
			return nil, fmt.Errorf("Error parsing bytes length: %w", err)
		}
	}
	c.Bytes = make([]byte, bytesLength)
	if _, err := io.ReadFull(reader, c.Bytes); err != nil {
		return nil, fmt.Errorf("reading value: %w", err)
	}

	// Read and discard the trailing \r\n
	if _, err := reader.ReadString('\n'); err != nil {
		return nil, fmt.Errorf("reading trailing newline: %w", err)
	}

	return c, nil

}

func (c *Cmd) parseCONNECT(payload []byte) (*Cmd, error) {
	var connectData ConnectCommand

	// Parse JSON payload into connectData struct
	err := json.Unmarshal(payload, &connectData)
	if err != nil {
		return nil, fmt.Errorf("Failed to parse CONNECT command: %w", err)
	}

	c.Name = CONNECT            // Set the command name
	c.ConnectData = connectData // Store the connection info (add ConnectData to Cmd struct)
	return c, nil
}

func (c *Cmd) String() string {
	return fmt.Sprintf("Name: %s, Subject: %s, ID: %s, Bytes: %s", c.Name, string(c.Subject), c.ID, c.Bytes)
}
