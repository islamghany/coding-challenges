package http

import (
	"fmt"
	"net"
)

var StatusReasons = map[int]string{
	200: "OK",
	201: "Created",
	202: "Accepted",
	204: "No Content",
	400: "Bad Request",
	401: "Unauthorized",
	403: "Forbidden",
	404: "Not Found",
	500: "Internal Server Error",
}

type Response struct {
	Protocol string
	Status   int
	Reason   string
	Headers  map[string]string
	Body     []byte
}

type ResponseConfig struct {
	Protocol string
	Status   int
}

func (r *Response) SetStatus(code int) {
	r.Status = code
	r.Reason = StatusReasons[code]
}

func (r *Response) GetHeader(key string) string {
	return r.Headers[key]
}

func (r *Response) SetHeader(key, value string) {
	r.Headers[key] = value
}

func (r *Response) RemoveHeader(key string) {

	delete(r.Headers, key)
}

func NewResponse(cfg *ResponseConfig) *Response {
	response := &Response{
		Headers: make(map[string]string),
	}
	if cfg == nil || cfg.Protocol == "" {
		response.Protocol = "HTTP/1.1"
	} else {
		response.Protocol = cfg.Protocol
	}

	if cfg == nil || cfg.Status == 0 {
		response.Status = 200
		response.Reason = "OK"
	} else {
		response.Status = cfg.Status
		response.Reason = StatusReasons[cfg.Status]
	}

	return response
}

// func ParseIncomingResponse

func WriteResponseToConn(conn net.Conn, response *Response) error {
	// Write the response line
	responseLine := fmt.Sprintf("%s %d %s\r\n", response.Protocol, response.Status, response.Reason)
	_, err := conn.Write([]byte(responseLine))
	if err != nil {
		return err
	}

	// Write the headers
	for key, value := range response.Headers {
		header := fmt.Sprintf("%s: %s\r\n", key, value)
		_, err := conn.Write([]byte(header))
		if err != nil {
			return err
		}
	}

	// Write the final CRLF
	_, err = conn.Write([]byte("\r\n"))
	if err != nil {
		return err
	}

	// Write the body if present
	if len(response.Body) > 0 {
		response.SetHeader("Content-Length", fmt.Sprintf("%d", len(response.Body)))
		_, err = conn.Write(response.Body)
		if err != nil {
			return err
		}
	}

	return nil
}
