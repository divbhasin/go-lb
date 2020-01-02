package main

import (
	"log"
	"net/url"
	"sync/atomic"
)

// ServerPool : Tracks the array of servers and the current active one
type ServerPool struct {
	backends []*Backend
	current  uint64
}

// NextIndex : Gets the index of the next server in line
func (s *ServerPool) NextIndex() int {
	return int(atomic.AddUint64(&s.current, uint64(1)) % uint64(len(s.backends)))
}

// GetNextPeer : Gets the next server to send a request to and sets the current of the server pool
func (s *ServerPool) GetNextPeer() *Backend {
	next := s.NextIndex()
	l := len(s.backends) + next
	for i := next; i < l; i++ {
		idx := i % len(s.backends)
		if s.backends[idx].IsAlive() {
			if i != next {
				atomic.StoreUint64(&s.current, uint64(idx))
			}
			return s.backends[idx]
		}
	}
	return nil
}

// MarkBackendStatus : Marks the status of a backend to be alive or not alive
func (s *ServerPool) MarkBackendStatus(u *url.URL, status bool) {
	for _, backend := range s.backends {
		if backend.URL.Host == u.Host {
			backend.SetAlive(false)
		}
	}
}

// HealthCheck : does a health check on all of the servers by sending a TCP request
func (s *ServerPool) HealthCheck() {
	for _, b := range s.backends {
		status := "up"
		alive := isBackendAlive(b.URL)
		b.SetAlive(alive)
		if !alive {
			status = "down"
		}
		log.Printf("%s [%s]\n", b.URL, status)
	}
}
