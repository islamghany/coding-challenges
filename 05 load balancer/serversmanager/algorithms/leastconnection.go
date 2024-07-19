package algorithms

import (
	"lb/serversmanager"
	"sync"
)

type LeastConnection struct {
	mux sync.Mutex
}

func NewLeastConnection() *LeastConnection {
	return &LeastConnection{}
}

func (lc *LeastConnection) Next(servers []*serversmanager.ServerManager) *serversmanager.ServerManager {
	lc.mux.Lock()
	defer lc.mux.Unlock()
	var server *serversmanager.ServerManager
	// Find the server with the least number of requests
	for _, s := range servers {
		if server == nil {
			server = s
			continue
		}
		if s.GetTotalRequests() < server.GetTotalRequests() {
			server = s
		}

	}
	return server
}
