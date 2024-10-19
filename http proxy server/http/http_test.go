package http

import (
	"net"
	"testing"
)

// Helper function to simulate a client-server connection
func mockConnection(input string) net.Conn {
	client, server := net.Pipe()
	go func() {
		server.Write([]byte(input))
		server.Close()
	}()
	return client
}
func TestParseIncomingRequest(t *testing.T) {
	conn := mockConnection("POST /hello HTTP/1.1\r\nHost: localhost\r\nContent-Length: 5\r\n\r\nHello")
	defer conn.Close()

	request, err := ParseIncomingRequest(conn)
	if err != nil {
		t.Fatal(err)
	}

	if request.Method != "POST" {
		t.Errorf("Expected GET, got %s", request.Method)
	}
	if request.URI != "/hello" {
		t.Errorf("Expected /hello, got %s", request.URI)
	}
	if request.Protocol != "HTTP/1.1" {
		t.Errorf("Expected HTTP/1.1, got %s", request.Protocol)
	}
	if request.GetHeader("Host") != "localhost" {
		t.Errorf("Expected localhost, got %s", request.GetHeader("Host"))
	}
	if request.GetHeader("Content-Length") != "5" {
		t.Errorf("Expected 5, got %s", request.GetHeader("Content-Length"))
	}
	if string(request.Body) != "Hello" {
		t.Errorf("Expected Hello, got %s", string(request.Body))
	}

}
