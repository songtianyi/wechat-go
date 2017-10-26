package rrserver

import (
	"encoding/binary"
	"io"
	"net"
	"sync"
)

const (
	PT_SIZE          = uint32(512) // packet len in byte
	PT_SIZE_BYTE_LEN = 4           // packet len value in bits
)

type TCPConnection struct {
	conn  net.Conn
	r     io.Reader
	w     io.Writer
	rLock sync.Mutex
	wLock sync.Mutex
}

func NewTCPConnection(conn net.Conn) *TCPConnection {
	return &TCPConnection{
		conn: conn,
		r:    conn,
		w:    conn,
	}
}

func (c *TCPConnection) Close() error {
	return c.conn.Close()
}

func (c *TCPConnection) SetKeepAlive(v bool) error {
	return c.conn.(*net.TCPConn).SetKeepAlive(v)
}

func (c *TCPConnection) RemoteAddr() string {
	return c.conn.RemoteAddr().String()
}
func (c *TCPConnection) LocalAddr() string {
	return c.conn.LocalAddr().String()
}

func (c *TCPConnection) Read() (error, []byte) {
	c.rLock.Lock()
	defer c.rLock.Unlock()
	buf := make([]byte, PT_SIZE)
	if _, err := io.ReadFull(c.r, buf[:PT_SIZE_BYTE_LEN]); err != nil {
		return err, buf
	}
	pl := binary.BigEndian.Uint32(buf[:PT_SIZE_BYTE_LEN])
	if pl > PT_SIZE {
		buf = make([]byte, pl)
	}
	if _, err := io.ReadFull(c.r, buf[:pl]); err != nil {
		return err, buf
	}
	return nil, buf[:pl]
}

func (c *TCPConnection) Write(msg []byte) error {
	c.wLock.Lock()
	defer c.wLock.Unlock()
	buf := make([]byte, PT_SIZE_BYTE_LEN)
	binary.BigEndian.PutUint32(buf[:PT_SIZE_BYTE_LEN], uint32(len(msg)))
	if _, err := c.w.Write(buf[:PT_SIZE_BYTE_LEN]); err != nil {
		return err
	}
	if _, err := c.w.Write(msg); err != nil {
		return err
	}
	return nil
}
