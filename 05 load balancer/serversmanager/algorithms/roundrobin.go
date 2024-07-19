package algorithms

import (
	"lb/serversmanager"
	"sync"
)

type RoundRobin struct {
	pointer int
	mux     sync.Mutex
}

func NewRoundRobin() *RoundRobin {
	return &RoundRobin{
		pointer: 0,
	}
}

func (rr *RoundRobin) Next(servers []*serversmanager.ServerManager) *serversmanager.ServerManager {
	rr.mux.Lock()
	defer rr.mux.Unlock()
	server := servers[rr.pointer]
	rr.pointer = (rr.pointer + 1) % len(servers)
	return server
}
