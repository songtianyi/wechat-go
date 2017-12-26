/*
Copyright 2017 wechat-go Authors. All Rights Reserved.
MIT License

Copyright (c) 2017

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
*/

package wxweb

import (
	"fmt"
	"sync"
)

// Handler: message function wrapper
type Handler func(*Session, *ReceivedMessage)

// HandlerWrapper: message handler wrapper
type HandlerWrapper struct {
	handle  Handler
	enabled bool
	name    string
}

// Run: run message handler
func (s *HandlerWrapper) Run(session *Session, msg *ReceivedMessage) {
	if s.enabled {
		s.handle(session, msg)
	}
}

func (s *HandlerWrapper) getName() string {
	return s.name
}

// GetName: export name for external projects
func (s *HandlerWrapper) GetName() string {
	return s.name
}

// GetEnabled: export enabled for external projects
func (s *HandlerWrapper) GetEnabled() bool {
	return s.enabled
}

func (s *HandlerWrapper) enableHandle() {
	s.enabled = true
	return
}

func (s *HandlerWrapper) disableHandle() {
	s.enabled = false
	return
}

// HandlerRegister: message handler manager
type HandlerRegister struct {
	mu   sync.RWMutex
	hmap map[int][]*HandlerWrapper
}

// CreateHandlerRegister: create handler register
func CreateHandlerRegister() *HandlerRegister {
	return &HandlerRegister{
		hmap: make(map[int][]*HandlerWrapper),
	}
}

// Add: add message handler to handler register
func (hr *HandlerRegister) Add(key int, h Handler, name string) error {
	hr.mu.Lock()
	defer hr.mu.Unlock()
	for _, v := range hr.hmap {
		for _, handle := range v {
			if handle.getName() == name {
				return fmt.Errorf("handler name %s has been registered", name)
			}
		}
	}
	hr.hmap[key] = append(hr.hmap[key], &HandlerWrapper{handle: h, enabled: false, name: name})
	return nil
}

// Get: get message handler
func (hr *HandlerRegister) Get(key int) (error, []*HandlerWrapper) {
	hr.mu.RLock()
	defer hr.mu.RUnlock()
	if v, ok := hr.hmap[key]; ok {
		return nil, v
	}
	return fmt.Errorf("no handlers for key [%d]", key), nil
}

// GetAll: get all message handlers
func (hr *HandlerRegister) GetAll() []*HandlerWrapper {
	hr.mu.RLock()
	defer hr.mu.RUnlock()
	result := make([]*HandlerWrapper, 0)
	for _, v := range hr.hmap {
		result = append(result, v...)
	}
	return result
}

// EnableByType: enable handler by message type
func (hr *HandlerRegister) EnableByType(key int) error {
	err, handles := hr.Get(key)
	if err != nil {
		return err
	}
	hr.mu.Lock()
	defer hr.mu.Unlock()
	// all
	for _, v := range handles {
		v.enableHandle()
	}
	return nil
}

// DisableByType: disable handler by message type
func (hr *HandlerRegister) DisableByType(key int) error {
	err, handles := hr.Get(key)
	if err != nil {
		return err
	}
	hr.mu.Lock()
	defer hr.mu.Unlock()
	// all
	for _, v := range handles {
		v.disableHandle()
	}
	return nil
}

// EnableByName: enable message handler by name
func (hr *HandlerRegister) EnableByName(name string) error {
	hr.mu.Lock()
	defer hr.mu.Unlock()
	for _, handles := range hr.hmap {
		for _, v := range handles {
			if v.getName() == name {
				v.enableHandle()
				return nil
			}
		}
	}
	return fmt.Errorf("cannot find handler %s", name)
}

// DisableByName: disable message handler by name
func (hr *HandlerRegister) DisableByName(name string) error {
	hr.mu.Lock()
	defer hr.mu.Unlock()
	for _, handles := range hr.hmap {
		for _, v := range handles {
			if v.getName() == name {
				v.disableHandle()
				return nil
			}
		}
	}
	return fmt.Errorf("cannot find handler %s", name)
}

// Dump: output all message handlers
func (hr *HandlerRegister) Dump() string {
	hr.mu.RLock()
	defer hr.mu.RUnlock()
	str := "[plugins dump]\n"
	for k, handles := range hr.hmap {
		for _, v := range handles {
			str += fmt.Sprintf("%d %s [%v]\n", k, v.getName(), v.enabled)
		}
	}
	return str
}
