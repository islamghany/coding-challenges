package serversmanager

import (
	"fmt"
	"net"
	"sync"
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
	Weight  int
}

type ServerManager struct {
	status            string
	url               string
	timeout           time.Duration
	totalRequests     int
	totalResponseTime int
	weight            int
	mux               *sync.Mutex
}

func NewServerManager(cfg Config) *ServerManager {
	sv := &ServerManager{
		status:            DEAD,
		timeout:           time.Second * 5,
		url:               cfg.Url,
		totalRequests:     0,
		totalResponseTime: 0,
		weight:            0,
		mux:               &sync.Mutex{},
	}
	if cfg.Timeout != 0 {
		sv.timeout = cfg.Timeout
	}
	if cfg.Weight != 0 {
		sv.weight = cfg.Weight
	}

	return sv
}

func (s *ServerManager) Check() {
	conn, err := net.DialTimeout("tcp", s.url, s.timeout)
	if err != nil {
		fmt.Println("Server is dead", s.url)
		s.mux.Lock()
		s.status = DEAD
		s.mux.Unlock()
		return
	}
	defer conn.Close()
	fmt.Println("Server is live")
	s.mux.Lock()
	s.status = LIVE
	s.mux.Unlock()
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
	s.mux.Lock()
	s.status = DEAD
	s.mux.Unlock()
	return nil, ErrorServerDead

}

func (s *ServerManager) IncrementTotalRequests() {
	s.mux.Lock()
	defer s.mux.Unlock()
	s.totalRequests++
}
func (s *ServerManager) DecrementTotalRequests() {
	s.mux.Lock()
	defer s.mux.Unlock()
	s.totalRequests--
}

func (s *ServerManager) GetTotalRequests() int {
	s.mux.Lock()
	defer s.mux.Unlock()
	return s.totalRequests
}

func (s *ServerManager) GeURL() string {
	return s.url
}
