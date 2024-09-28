package commandsparser

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
)

type CommandName string

const (
	// Commands
	SetCommand    CommandName = "set"
	GetCommand    CommandName = "get"
	DeleteCommand CommandName = "delete"
	ExitCommand   CommandName = "exit"
	EndCommand    CommandName = "end"
)

var (
	ErrInvalidCommand = fmt.Errorf("invalid command")
	ErrInvalidFormat  = fmt.Errorf("invalid command format")
)

// StorageCommand represents a storage command
// <command name> <key> <flags> <exptime> <bytes> [noreply]\r\n
// <data block>\r\n
type Command struct {
	Name    CommandName
	Key     string
	Value   []byte
	Flags   uint32
	Expiry  int64
	Bytes   uint32
	Noreply bool
}

type Parser struct{}

func NewParser() *Parser {
	return &Parser{}
}

func (cp *Parser) Parse(reader io.Reader) (*Command, error) {
	buffReader := bufio.NewReader(reader)
	line, err := buffReader.ReadString('\n')
	if err != nil {
		return nil, fmt.Errorf("reading command line: %w", err)
	}

	fields := strings.Fields(strings.TrimSpace(line))
	if len(fields) == 0 {
		return nil, ErrInvalidCommand
	}
	cmd := &Command{Name: CommandName(fields[0])}

	switch cmd.Name {
	case SetCommand:
		return cp.parseSet(fields, buffReader)
	case GetCommand:
		return cp.parseGet(fields)
	case DeleteCommand:
		return cp.parseDelete(fields)
	case ExitCommand, EndCommand:
		return cmd, nil
	default:
		return nil, fmt.Errorf("Unknown command")
	}

}

func (cp *Parser) parseSet(fields []string, reader *bufio.Reader) (*Command, error) {

	if len(fields) < 5 {
		return nil, ErrInvalidCommand
	}
	cmd := &Command{
		Name: SetCommand,
		Key:  fields[1],
	}
	var err error

	if cmd.Flags, err = parseUint32(fields[2]); err != nil {
		return nil, fmt.Errorf("parsing flags: %w", err)
	}
	if cmd.Expiry, err = parseInt64(fields[3]); err != nil {
		return nil, fmt.Errorf("parsing expiry: %w", err)
	}
	if cmd.Bytes, err = parseUint32(fields[4]); err != nil {
		return nil, fmt.Errorf("parsing bytes: %w", err)
	}
	cmd.Noreply = len(fields) == 6 && fields[5] == "noreply"
	cmd.Value = make([]byte, cmd.Bytes)
	if _, err := io.ReadFull(reader, cmd.Value); err != nil {
		return nil, fmt.Errorf("reading value: %w", err)
	}

	// Read and discard the trailing \r\n
	if _, err := reader.ReadString('\n'); err != nil {
		return nil, fmt.Errorf("reading trailing newline: %w", err)
	}

	return cmd, nil
}

func (p *Parser) parseGet(fields []string) (*Command, error) {
	if len(fields) < 2 {
		return nil, ErrInvalidFormat
	}

	return &Command{
		Name: GetCommand,
		Key:  fields[1],
	}, nil

}
func (p *Parser) parseDelete(fields []string) (*Command, error) {
	if len(fields) < 2 {
		return nil, ErrInvalidFormat
	}

	return &Command{
		Name:    DeleteCommand,
		Key:     fields[1],
		Noreply: len(fields) == 3 && fields[2] == "noreply",
	}, nil
}

func parseUint32(s string) (uint32, error) {
	v, err := strconv.ParseUint(s, 10, 32)
	return uint32(v), err
}

func parseInt64(s string) (int64, error) {
	return strconv.ParseInt(s, 10, 64)
}

func parseUint16(s string) (uint16, error) {
	v, err := strconv.ParseUint(s, 10, 16)
	return uint16(v), err
}
