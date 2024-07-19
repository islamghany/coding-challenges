package algorithms

import (
	"fmt"
	"lb/serversmanager"
)

const (
	ROUNDROBIN         = "round_robin"
	LEASTCONNECTION    = "least_connection"
	WEIGHTEDROUNDROBIN = "weighted_round_robin"
	LEASTRESPONSETIME  = "least_response_time"
)

var availableAlgorithms []string = []string{
	ROUNDROBIN,
	LEASTCONNECTION,
	WEIGHTEDROUNDROBIN,
	LEASTRESPONSETIME,
}

func IsAlgorithmAvailable(algorithm string) error {
	for _, a := range availableAlgorithms {
		if a == algorithm {
			return nil
		}
	}
	return fmt.Errorf("algorithm name isnot available, only available alogorithms are: %+v", availableAlgorithms)
}

type LoadBalancerAlgorithm interface {
	Next(servers []*serversmanager.ServerManager) *serversmanager.ServerManager
}
