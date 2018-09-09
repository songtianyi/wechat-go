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

	rwLock sync.RWMutex
}

func (s *KVStorage) Put(key string, value interface{}) {
	s.rwLock.Lock()
	s.meta[key] = value
	s.rwLock.Unlock()
}

func (s *KVStorage) Get(key string) interface{} {
	s.rwLock.RLock()
	if v, ok := s.meta[key]; ok {
		return v
	}
	s.rwLock.RUnlock()
	return nil
}
