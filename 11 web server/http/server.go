package http

import (
	"bufio"
	"fmt"
	"net"
	"strconv"
	"strings"
)

type ServerConfig struct {
	Addr string
}

type Server struct {
	addr string
	// Handler  Handler
	listener net.Listener
	routes   Routes
}

type Routes struct {
	routes map[string]Handler
}

type Request struct {
	Method        string
	Path          string
	Header        map[string]string
	Body          []byte
	Proto         string
	ContentLength int
}

func (w Response) Write(b []byte) (int, error) {
	data := EncodeResponse(b, w.Status, w.StatusText, w.Header)
	return w.conn.Write(data)
}

type Response struct {
	Header     map[string]string
	conn       net.Conn
	Status     int
	StatusText string
}

type Handler interface {
	ServeHTTP(Response, *Request)
}

type HandlerFunc func(Response, *Request)

func (f HandlerFunc) ServeHTTP(w Response, r *Request) {
	f(w, r)
}

func NewServer(cfg ServerConfig) *Server {
	return &Server{
		addr: cfg.Addr,
		routes: Routes{
			routes: make(map[string]Handler),
		},
	}
}

func (s *Server) Start() error {
	ln, err := net.Listen("tcp", s.addr)
	if err != nil {
		return err
	}
	s.listener = ln
	for {
		conn, err := ln.Accept()
		if err != nil {
			return err
		}
		go s.handleConnection(conn)
	}
}

func (s *Server) handleConnection(conn net.Conn) {
	defer conn.Close()
	// 1- read request
	req, err := s.handleReadRequest(conn)
	if err != nil {
		fmt.Println("error reading request: ", err)
		return
	}
	fmt.Printf("Request Body: %+v\n", string(req.Body))
	// 2- create response
	res := Response{
		conn:       conn,
		Header:     make(map[string]string),
		Status:     200,
		StatusText: "OK",
	}
	// 3- find route
	route, ok := s.routes.routes["/v1/hello"]
	if !ok {
		// 404
	}

	// 4- call handler
	route.ServeHTTP(res, req)
}

func (s *Server) handleReadRequest(conn net.Conn) (*Request, error) {
	req := &Request{
		Header: make(map[string]string),
	}
	reader := bufio.NewReader(conn)
	// 1- read the request line
	line, err := reader.ReadString('\n')
	if err != nil {
		return nil, fmt.Errorf("error reading request line: %v", err)
	}

	parts := strings.Split(line, " ")
	if len(parts) != 3 {
		return nil, fmt.Errorf("invalid request line: %s", line)
	}
	req.Method = parts[0]
	req.Path = parts[1]
	req.Proto = parts[2]
	fmt.Println(req.Method, req.Path, req.Proto)
	// 2- read the headers
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			return nil, fmt.Errorf("error reading header: %v", err)
		}
		if line == "\r\n" {
			break
		}
		parts := strings.SplitN(line, ":", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid header: %s", line)
		}
		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])
		req.Header[key] = value
		if strings.ToLower(key) == "content-length" {
			req.ContentLength, _ = strconv.Atoi(value)
		}
	}
	fmt.Println("content length: ", req.ContentLength)
	// 3- read the body if exists
	contentLength := req.ContentLength
	for contentLength > 0 {
		b, err := reader.ReadByte()
		if err != nil {
			return nil, fmt.Errorf("error reading body: %v", err)
		}
		req.Body = append(req.Body, b)
		contentLength--
		fmt.Println("content length: ", contentLength, "body: ", string(req.Body))
	}

	fmt.Printf("Request: %+v\n", req)
	return req, nil
}

func (s *Server) Get(path string, h HandlerFunc) {
	s.routes.routes[path] = h
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
