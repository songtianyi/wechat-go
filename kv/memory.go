package kv

import (
	"sync"
)

var (
	KVStorageInstance *KVStorage
)

func init() {
	KVStorageInstance = &KVStorage{
		meta: make(map[string]interface{}),
	}
}

type KVStorage struct {
	meta map[string]interface{}

	rLock sync.RWMutex
	wLock sync.RWMutex
}

func (s *KVStorage) Put(key string, value interface{}) {
	s.wLock.Lock()
	s.meta[key] = value
	s.wLock.Unlock()
}

func (s *KVStorage) Get(key string) interface{} {
	s.rLock.RLock()
	if v, ok := s.meta[key]; ok {
		return v
	}
	s.rLock.RUnlock()
	return nil
}
