package resp

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"strconv"
)

var (
	ErrMalformedType  = errors.New("malformed type")
	ErrInvalidInteger = errors.New("invalid integer")
	ErrUnknownType    = errors.New("unknown type")
)

type RESPType int

const (
	SimpleString RESPType = iota
	Error
	Integer
	BulkString
	Array
)

type RESPData struct {
	Data interface{}
	Type RESPType
}

func Deserialize(reader *bufio.Reader) (*RESPData, error) {
	firstByte, err := reader.ReadByte()
	if err != nil {
		return nil, err
	}
	switch firstByte {
	case '+':
		return readSimpleString(reader)
	case '-':
		return readSimpleError(reader)
	case ':':
		return readInteger(reader)
	case '$':
		return readBulkString(reader)
	case '*':
		return readArray(reader)
	default:
		return nil, ErrUnknownType

	}
}

// readSimpleString is a helper function that deserializes a simple string
// from the RESP protocol.
func readSimpleString(reader *bufio.Reader) (*RESPData, error) {
	//  Simple String: Starts with "+" and ends with "\r\n"
	//  Example: +OK\r\n

	line, err := reader.ReadString('\n')
	if err != nil {
		return nil, err
	}

	return &RESPData{
		Data: line[:len(line)-2],
		Type: SimpleString,
	}, nil
}

// readSimpleError is a helper function that deserializes a simple error
// from the RESP protocol.
func readSimpleError(reader *bufio.Reader) (*RESPData, error) {
	//  Simple Error: Starts with "-" and ends with "\r\n"
	//  Example: -ERR\r\n

	line, err := reader.ReadString('\n')

	if err != nil {
		fmt.Println("Erssror", err, errors.Is(err, io.EOF))
		return nil, ErrMalformedType
	}

	return &RESPData{
		Data: line[:len(line)-2],
		Type: Error,
	}, nil
}

// readInteger is a helper function that deserializes an integer
// from the RESP protocol.
func readInteger(reader *bufio.Reader) (*RESPData, error) {
	//  Integer: Starts with ":" and ends with "\r\n"
	//  Example: :1000\r\n

	line, err := reader.ReadString('\n')
	if err != nil {
		return nil, ErrMalformedType
	}
	val, err := strconv.Atoi(line[:len(line)-2])
	if err != nil {
		return nil, ErrInvalidInteger
	}
	return &RESPData{
		Data: val,
		Type: Integer,
	}, nil
}

func readBulkString(reader *bufio.Reader) (*RESPData, error) {
	// Bulk String: Starts with "$" then the length of then "\r\n", then the string itself and ends with "\r\n"
	// Example: $6\r\nfoobar\r\n

	line, err := reader.ReadString('\n')
	if err != nil {
		return nil, err
	}
	length, err := strconv.Atoi(line[:len(line)-2])
	if err != nil {
		return nil, ErrInvalidInteger
	}
	// if length is -1, then the string is nil
	if length == -1 {
		return &RESPData{
			Data: nil,
			Type: BulkString,
		}, nil
	}
	buf := make([]byte, length+2)
	// then we need to read the string itself
	_, err = io.ReadFull(reader, buf)
	if err != nil {
		return nil, ErrMalformedType
	}

	return &RESPData{
		Data: string(buf[:length]),
		Type: BulkString,
	}, nil

}

func readArray(reader *bufio.Reader) (*RESPData, error) {
	// Array: start with "*" followed by the number of elements, each element is then desrialized Individually
	// e.g *<number_of_elements>\r\n<element1><element2>...<elementN>"
	// example: "*2\r\n$4\r\necho\r\n$11\r\nhello world\r\n"

	// get the length of the array
	line, err := reader.ReadString('\n')
	length, err := strconv.Atoi(line[:len(line)-2])
	if err != nil {
		return nil, ErrInvalidInteger
	}
	// if length is -1, then the array is nil
	if length == -1 {
		return &RESPData{
			Data: nil,
			Type: Array,
		}, nil
	}
	respArray := make([]RESPData, length)
	for i := 0; i < length; i++ {
		element, err := Deserialize(reader)
		if err != nil {
			return nil, err
		}
		respArray[i] = *element
	}
	return &RESPData{
		Data: respArray,
		Type: Array,
	}, nil

}

func NewError(msg string) RESPData {
	return RESPData{
		Data: msg,
		Type: Error,
	}
}
func NewSimpleString(msg string) RESPData {
	return RESPData{
		Data: msg,
		Type: SimpleString,
	}
}

func NewBulkString(msg string) RESPData {
	return RESPData{
		Data: msg,
		Type: BulkString,
	}
}

func NewInteger(num string) RESPData {
	return RESPData{
		Data: num,
		Type: Integer,
	}
}

func NewArray(elems []string) []RESPData {
	length := len(elems)
	respArray := make([]RESPData, length)
	for i := 0; i < length; i++ {
		respArray[i] = NewBulkString(elems[i])
	}

	return respArray
}

func NewNil() RESPData {
	return RESPData{
		Data: nil,
		Type: BulkString,
	}
}
