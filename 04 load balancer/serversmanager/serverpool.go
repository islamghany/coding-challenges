package serversmanager

import (
	"fmt"
	"time"
)

var maxRetry int = 5

type ServerPool struct {
	servers          []*ServerManager
	recovringServers map[*ServerManager]bool
	retry            int
	pointer          int
	dead             int
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
			sp.dead--
			return
		}
		retryTimes++
		backoff *= 3
		if retryTimes >= maxRetry {
			// delete the server
			sp.deleteServer(srv)
			sp.dead--
			return
		}
	}

}

// Function to get the next server
// If the server is dead, it will be recovering
func (sp *ServerPool) GetNextServer() *ServerManager {
	for {
		// if the pointer is greater than the length of the servers, reset the pointer
		if sp.pointer >= len(sp.servers) {
			sp.pointer = len(sp.servers) - 1
		}
		// if the length of the servers is 0, this means that there are no servers available
		if len(sp.servers) == 0 {
			return nil
		}
		server := sp.servers[sp.pointer]
		sp.pointer = (sp.pointer + 1) % len(sp.servers)
		// check if the server is active
		server.Check()
		if server.Active() {
			return server
		}
		if _, ok := sp.recovringServers[server]; !ok {
			sp.dead++
			fmt.Println("sever", server.url, "will enter the recovering state")
			sp.recovringServers[server] = true
			go sp.recoverServer(server)
		}
		// if the number of dead servers is equal to the number of servers, return nil
		if sp.dead == len(sp.servers) {
			return nil
		}

	}
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
	for i := 0; i < len(sp.servers); i++ {
		if sp.servers[i] == srv {
			fmt.Println("delete the server no:", srv.url)
			sp.servers = append(sp.servers[:i], sp.servers[i+1:]...)
			break
		}
	}
}
