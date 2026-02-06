package response

import (
	"fmt"
	"net"
	"strings"
)

func (w Response) Write(b []byte) (int, error) {
	data := EncodeResponse(b, w.Status, w.StatusText, w.Header)
	return w.Conn.Write(data)
}

type Response struct {
	Header     map[string]string
	Conn       net.Conn
	Status     int
	StatusText string
}

func EncodeResponse(body []byte, status int, statusText string, header map[string]string) []byte {
	var res strings.Builder
	res.WriteString(fmt.Sprintf("HTTP/1.1 %d %s\r\n", status, statusText))
	for k, v := range header {
		res.WriteString(fmt.Sprintf("%s: %s\r\n", k, v))
	}
	res.WriteString(fmt.Sprintf("Content-Length: %d\r\n", len(body)))
	res.WriteString("\r\n")
	res.Write(body)
	return []byte(res.String())
}
