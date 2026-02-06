package serversmanager

import (
	"fmt"
	"sync"
	"time"
)

var maxRetry int = 10

type ServerPool struct {
	servers          []*ServerManager
	recovringServers map[*ServerManager]bool
	retry            int
	pointer          int
	mutex            sync.Mutex
}

func NewServerPool(servers []*ServerManager, retry int) *ServerPool {
	return &ServerPool{
		servers:          servers,
		retry:            retry,
		pointer:          0,
		recovringServers: make(map[*ServerManager]bool),
	}
}

func (sp *ServerPool) recoverServer(srv *ServerManager) {
	// defer the function to remove the server from the recovering servers
	defer func(srv *ServerManager) {
		delete(sp.recovringServers, srv)
	}(srv)
	backoff := time.Second
	retryTimes := 0
	for {
		time.Sleep(backoff)
		srv.Check()
		if srv.Active() {
			fmt.Println("Server is back", srv.url)
			return
		}
		retryTimes++
		backoff *= 3
		if retryTimes >= maxRetry {
			// return if the server is not back after maxRetry
			fmt.Println("Server is not back", srv.url)
			return
		}
	}

}

func (sp *ServerPool) GetAllActiveServers() []*ServerManager {
	sp.mutex.Lock()
	defer sp.mutex.Unlock()
	var activeServers []*ServerManager
	for _, server := range sp.servers {
		server.Check()
		if server.Active() {
			activeServers = append(activeServers, server)
			continue
		}
		if _, ok := sp.recovringServers[server]; !ok {
			fmt.Println("sever", server.url, "will enter the recovering state")
			sp.recovringServers[server] = true
			go sp.recoverServer(server)
		}
	}
	return activeServers
}

// Function to bootup all servers
func (sp *ServerPool) BootupServers() {
	for i := 0; i < len(sp.servers); i++ {
		server := sp.servers[i]
		server.Check()
	}
}

// Function to delete an element at a specific index
func (sp *ServerPool) deleteServer(srv *ServerManager) {
	sp.mutex.Lock()         // Lock before modifying the slice
	defer sp.mutex.Unlock() // Unlock after modification
	for i := 0; i < len(sp.servers); i++ {
		if sp.servers[i] == srv {
			fmt.Println("delete the server no:", srv.url)
			sp.servers = append(sp.servers[:i], sp.servers[i+1:]...)
			break
		}
	}
}
