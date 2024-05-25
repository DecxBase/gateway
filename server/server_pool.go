package server

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/DecxBase/gateway/registry"
)

type ServerPool interface {
	GetBackends() []registry.Backend
	GetNextValidPeer() registry.Backend
	AddBackend(registry.Backend)
	GetServerPoolSize() int
}

type roundRobinServerPool struct {
	backends []registry.Backend
	mux      sync.RWMutex
	current  int
}

func (s *roundRobinServerPool) Rotate() registry.Backend {
	s.mux.Lock()
	s.current = (s.current + 1) % s.GetServerPoolSize()
	s.mux.Unlock()
	return s.backends[s.current]
}

func (s *roundRobinServerPool) GetNextValidPeer() registry.Backend {
	for i := 0; i < s.GetServerPoolSize(); i++ {
		nextPeer := s.Rotate()
		if nextPeer.IsAlive() {
			return nextPeer
		}
	}
	return nil
}

func (s *roundRobinServerPool) GetBackends() []registry.Backend {
	return s.backends
}

func (s *roundRobinServerPool) AddBackend(b registry.Backend) {
	s.backends = append(s.backends, b)
}

func (s *roundRobinServerPool) GetServerPoolSize() int {
	return len(s.backends)
}

func HealthCheck(ctx context.Context, s ServerPool) {
	aliveChannel := make(chan bool, 1)

	for _, b := range s.GetBackends() {
		b := b
		requestCtx, stop := context.WithTimeout(ctx, 10*time.Second)
		defer stop()
		status := "up"
		go registry.IsBackendAlive(requestCtx, aliveChannel, b.GetURL())

		select {
		case <-ctx.Done():
			registry.Logger.Info().Msg("Gracefully shutting down health check")
			return
		case alive := <-aliveChannel:
			b.SetAlive(alive)
			if !alive {
				status = "down"
			}
		}

		registry.Logger.Debug().Str("url", b.GetURL().String()).Str("status", status).Msg("URL Status")
	}
}

func NewServerPool(strategy registry.LBStrategy) (ServerPool, error) {
	switch strategy {
	case registry.RoundRobin:
		return &roundRobinServerPool{
			backends: make([]registry.Backend, 0),
			current:  0,
		}, nil
	case registry.LeastConnected:
		return &lcServerPool{
			backends: make([]registry.Backend, 0),
		}, nil
	default:
		return nil, fmt.Errorf("Invalid strategy")
	}
}
