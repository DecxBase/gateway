package server

import (
	"sync"

	"github.com/DecxBase/gateway/registry"
)

type lcServerPool struct {
	backends []registry.Backend
	mux      sync.RWMutex
}

func (s *lcServerPool) GetNextValidPeer() registry.Backend {
	var leastConnectedPeer registry.Backend
	for _, b := range s.backends {
		if b.IsAlive() {
			leastConnectedPeer = b
			break
		}
	}

	for _, b := range s.backends {
		if !b.IsAlive() {
			continue
		}
		if leastConnectedPeer.GetActiveConnections() > b.GetActiveConnections() {
			leastConnectedPeer = b
		}
	}
	return leastConnectedPeer
}

func (s *lcServerPool) AddBackend(b registry.Backend) {
	s.backends = append(s.backends, b)
}

func (s *lcServerPool) GetServerPoolSize() int {
	return len(s.backends)
}

func (s *lcServerPool) GetBackends() []registry.Backend {
	return s.backends
}
