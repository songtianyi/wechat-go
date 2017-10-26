package rrhandler

import (
	"fmt"
	"sync"
	"time"
)

type HandlerRegister struct {
	mu   *sync.RWMutex
	hmap map[string]*HandlerWrapper
}

func CreateHandlerRegister() (error, *HandlerRegister) {
	return nil, &HandlerRegister{
		mu:   new(sync.RWMutex),
		hmap: make(map[string]*HandlerWrapper),
	}
}

func (hr *HandlerRegister) Add(key string, h Handler, t time.Duration) {
	hr.mu.Lock()
	defer hr.mu.Unlock()
	hr.hmap[key] = &HandlerWrapper{
		handle:  h,
		timeout: t,
	}
}

func (hr *HandlerRegister) Get(key string) (error, *HandlerWrapper) {
	hr.mu.RLock()
	defer hr.mu.RUnlock()
	if v, ok := hr.hmap[key]; ok {
		return nil, v
	}
	return fmt.Errorf("value for key [%d] does not exist in map", key), nil
}
