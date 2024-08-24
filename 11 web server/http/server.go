package http

import (
	"bufio"
	"context"
	"fmt"
	"net"
	"strconv"
	"strings"
	"webserver/http/request"
	"webserver/http/response"
	"webserver/http/router"
)

type ServerConfig struct {
	Addr   string
	Router *router.Router
}

type Server struct {
	addr string
	// Handler  Handler
	listener net.Listener
	routes   *router.Router
}

func NewServer(cfg ServerConfig) *Server {
	if cfg.Router == nil {
		cfg.Router = router.NewRouter()
	}
	return &Server{
		addr:   cfg.Addr,
		routes: cfg.Router,
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
	res := response.Response{
		Conn:       conn,
		Header:     make(map[string]string),
		Status:     200,
		StatusText: "OK",
	}
	// serve the request
	s.routes.ServeHTTP(res, req)
}

func (s *Server) handleReadRequest(conn net.Conn) (*request.Request, error) {
	req := &request.Request{
		Header:  make(map[string]string),
		Context: context.Background(),
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
