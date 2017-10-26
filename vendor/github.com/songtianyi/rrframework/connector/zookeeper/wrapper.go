package rrzk

import (
	"fmt"
	"github.com/samuel/go-zookeeper/zk"
	"github.com/songtianyi/rrframework/logs"
	"strings"
	"sync"
	"time"
)

var (
	defaultTimeout = 1 * time.Second
)

type ZKClient struct {
	conn *zk.Conn
	ev   <-chan zk.Event // one-way channel
}

type clientPool struct {
	pool map[string]*ZKClient
	mu   sync.RWMutex
}

var (
	cp = &clientPool{pool: make(map[string]*ZKClient)}
)

func init() {
	// TODO
	// set zk logger to rrframework/logs
}

func (s *clientPool) add(sv string, c *ZKClient) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.pool[sv]; !ok {
		// not exist
		s.pool[sv] = c
	}
	// exist
}

func (s *clientPool) get(sv string) *ZKClient {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if v, ok := s.pool[sv]; ok {
		return v
	}
	return nil
}

func (s *clientPool) closeClient(sv string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.pool[sv]; !ok {
		// closed
		return
	}
	c, _ := s.pool[sv]
	c.conn.Close()
	delete(s.pool, sv)
}

func watchZkConn(sv string, ev <-chan zk.Event) {
	for {
		e := <-ev
		switch e.State {
		case zk.StateDisconnected:
			logs.Critical("zk [%s] disconnected", sv)
			// close old client
			cp.closeClient(sv)
			// reconnect
			Connect(sv, defaultTimeout)
			return
		case zk.StateConnected:
			logs.Debug("Zk cluster [%s] connected", sv)
			continue
		case zk.StateConnecting:
			logs.Info("Connecting to zk [%s]", sv)
		default:
			continue
		}
	}
}

// string sv: servers split by character ','
func Connect(sv string, timeout time.Duration) error {
	servers := strings.Split(sv, ",")
	conn, ev, err := zk.Connect(servers, timeout)
	if err != nil {
		return fmt.Errorf("Connect to servers [%s] fail, %s", sv, err)
	}
	cp.add(sv, &ZKClient{conn: conn, ev: ev})
	go watchZkConn(sv, ev)
	return nil
}

func GetZkClient(sv string) (error, *ZKClient) {
	if c := cp.get(sv); c != nil {
		return nil, c
	}
	if err := Connect(sv, defaultTimeout); err != nil {
		return err, nil
	}
	return GetZkClient(sv)
}

//func (s *ZkConn)
