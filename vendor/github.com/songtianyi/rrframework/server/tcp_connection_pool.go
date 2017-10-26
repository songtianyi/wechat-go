package rrserver

import (
	"fmt"
	"sync"
)

type TCPConnectionPool struct {
	pool     map[string]*channelPool
	mu       sync.RWMutex
	poolSize int
}

func CreateTCPConnectionPool(poolSize int) *TCPConnectionPool {
	return &TCPConnectionPool{
		pool:     make(map[string]*channelPool),
		poolSize: 10,
	}
}

func (s *TCPConnectionPool) Get(addr string) (error, *TCPConnection) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.pool[addr]; !ok {
		// first ever
		s.pool[addr] = &channelPool{
			conns: make(chan *TCPConnection, s.poolSize),
		}
	}
	return s.pool[addr].get(addr)
}

func (s *TCPConnectionPool) Add(addr string, c *TCPConnection) error {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if _, ok := s.pool[addr]; !ok {
		return fmt.Errorf("the connection you passed may not belonged to me")
	}
	return s.pool[addr].add(c)
}

func (s *TCPConnectionPool) CloseAll(addr string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.pool[addr]; !ok {
		return
	}
	s.pool[addr].closePool()
	delete(s.pool, addr)
	return
}
