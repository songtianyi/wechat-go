package rrserver

import (
	"fmt"
	"net"
	"sync"
)

var (
	ErrClosed = fmt.Errorf("pool not created or closed")
)

type channelPool struct {
	mu    sync.Mutex
	conns chan *TCPConnection
}

func (c *channelPool) factory(addr string) (error, *TCPConnection) {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return err, nil
	}
	return nil, NewTCPConnection(conn)
}

func (c *channelPool) get(addr string) (error, *TCPConnection) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.conns == nil {
		return ErrClosed, nil
	}

	select {
	case conn, ok := <-c.conns:
		if !ok || conn == nil {
			return ErrClosed, nil
		}
		return nil, conn
	default:
		// create a new connection
		return c.factory(addr)
	}
}

func (c *channelPool) add(conn *TCPConnection) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	if conn == nil {
		return fmt.Errorf("connection is nil")
	}

	if c.conns == nil {
		// pool is closed, cann't put it into pool
		return conn.Close()
	}

	// put the resource back into the pool. If the pool is full, this will
	// block and the default case will be executed.
	select {
	case c.conns <- conn:
		return nil
	default:
		// pool is full, close passed connection
		return conn.Close()
	}
}

func (c *channelPool) closePool() {
	c.mu.Lock()
	defer c.mu.Unlock()
	// called when remote server down
	if c.conns == nil {
		return
	}

	close(c.conns)
	for conn := range c.conns {
		conn.Close()
	}
}
