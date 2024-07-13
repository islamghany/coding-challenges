package serversmanager

import (
	"fmt"
	"net"
	"time"
)

const (
	LIVE = "live"
	DEAD = "dead"
)

var (
	ErrorServerDead = fmt.Errorf("Server is dead")
)

type Config struct {
	Url     string
	Timeout time.Duration
}

type ServerManager struct {
	status  string
	url     string
	timeout time.Duration
}

func NewServerManager(cfg Config) *ServerManager {
	sv := &ServerManager{
		status:  DEAD,
		timeout: time.Second * 5,
		url:     cfg.Url,
	}
	if cfg.Timeout != 0 {
		sv.timeout = cfg.Timeout
	}

	return sv
}

func (s *ServerManager) Check() {
	conn, err := net.DialTimeout("tcp", s.url, s.timeout)
	if err != nil {
		fmt.Println("Server is dead", s.url)
		s.status = DEAD
		return
	}
	defer conn.Close()
	fmt.Println("Server is live")
	s.status = LIVE
}

func (s *ServerManager) Active() bool {
	return s.status == LIVE
}

func (s *ServerManager) Dial(retry int) (net.Conn, error) {
	backoff := time.Second
	for i := 0; i < retry; i++ {
		conn, err := net.DialTimeout("tcp", s.url, s.timeout)
		if err != nil {
			fmt.Println("Error connecting to server")
			time.Sleep(backoff)
			backoff *= 2
			continue
		}
		return conn, nil
	}
	s.status = DEAD
	return nil, ErrorServerDead

}
