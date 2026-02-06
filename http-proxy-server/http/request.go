package http

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"net"
	"strconv"
	"strings"
)

var (
	ErrInvalidRequest       = errors.New("Invalid request")
	ErrMissingContentLength = errors.New("Missing Content-Length header")
)

type Request struct {
	Method   string
	URI      string
	Protocol string
	Headers  map[string]string
	Body     []byte
}

func (r *Request) GetHeader(key string) string {
	return r.Headers[key]
}

func (r *Request) SetHeader(key, value string) {
	r.Headers[key] = value
}

func (r *Request) RemoveHeader(key string) {
	delete(r.Headers, key)
}

func ParseIncomingRequest(conn net.Conn) (*Request, error) {
	request := &Request{
		Headers: make(map[string]string), // Initialize the Headers map
	}
	reader := bufio.NewReader(conn)
	// 1- first parse the request line, (Method, URI, Protocol)
	requestLine, err := reader.ReadString('\n')
	if err != nil {
		return nil, err
	}
	requestLine = strings.TrimSpace(requestLine)
	parts := strings.Split(requestLine, " ")

	if len(parts) != 3 {
		return nil, ErrInvalidRequest
	}

	request.Method = parts[0]
	request.URI = parts[1]
	request.Protocol = parts[2]

	// 2- parse the headers
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			return nil, err
		}
		line = strings.TrimSpace(line)
		if line == "" {
			break
		}
		parts := strings.SplitN(line, ":", 2)
		if len(parts) != 2 {
			fmt.Println("Invalid header line: ", line)
			continue
		}
		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])
		request.Headers[key] = value
	}

	// 3- parse the body

	if request.Method == "POST" || request.Method == "PUT" || request.Method == "PATCH" {
		contentLengthStr, ok := request.Headers["Content-Length"]
		if !ok {
			return nil, ErrMissingContentLength
		}
		// convert contentLength to int
		contentLength, err := strconv.Atoi(contentLengthStr)
		if err != nil {
			return nil, err
		}

		// read contentLength bytes from the reader
		body := make([]byte, contentLength)
		_, err = io.ReadFull(reader, body)
		if err != nil {
			return nil, err
		}
		request.Body = body
	}

	return request, nil
}

func WriteRequestToConn(conn net.Conn, request *Request) error {
	// Write the request line
	requestLine := fmt.Sprintf("%s %s %s\r\n", request.Method, request.URI, request.Protocol)
	_, err := conn.Write([]byte(requestLine))
	if err != nil {
		return err
	}

	// Write the headers
	for key, value := range request.Headers {
		header := fmt.Sprintf("%s: %s\r\n", key, value)
		_, err = conn.Write([]byte(header))
		if err != nil {
			return err
		}
	}

	// Add an empty line to separate headers from the body
	_, err = conn.Write([]byte("\r\n"))
	if err != nil {
		return err
	}

	// Write the body
	if len(request.Body) > 0 {
		_, err = conn.Write(request.Body)
		if err != nil {
			return err
		}
	}

	return nil
}
