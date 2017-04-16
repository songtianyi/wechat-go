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

type Handler func(*Session, *ReceivedMessage)

type HandlerWrapper struct {
	handle Handler
	enabled  bool
	name string
}

func (s *HandlerWrapper) Run(session *Session, msg *ReceivedMessage) {
	if s.enabled {
		s.handle(session, msg)
	}
}

func (s *HandlerWrapper) Name() string {
	return s.name
}

func (s *HandlerWrapper) Enable() {
	s.enabled = true
	return
}

func (s *HandlerWrapper) Disable() {
	s.enabled = false
	return
}


type HandlerRegister struct {
	mu   sync.RWMutex
	hmap map[int][]*HandlerWrapper
}

func CreateHandlerRegister() *HandlerRegister {
	return &HandlerRegister{
		hmap: make(map[int][]*HandlerWrapper),
	}
}

func (hr *HandlerRegister) Add(key int, h Handler, name string) error {
	hr.mu.Lock()
	defer hr.mu.Unlock()
	for _, v := range hr.hmap {
		for _, handle := range hr.hmap {
			if handle.Name() == name {
				return fmt.Errorf("handler  name %s has been registered",  name)
			}
		}
	}
	hr.hmap[key] = append(hr.hmap[key], &HandlerWrapper{handle: h, enabled: true, name: name})
}

func (hr *HandlerRegister) Get(key int) (error, []*HandlerWrapper) {
	hr.mu.RLock()
	defer hr.mu.RUnlock()
	if v, ok := hr.hmap[key]; ok {
		return nil, v
	}
	return fmt.Errorf("handlers for key [%d] not registered", key), nil
}

func  (hr *HandlerRegister) Enable(key int, name string) error {
	err, handles :=  s.Get(key)
	if err != nil {
		return err
	}
	hr.mu.RLock()
	defer hr.mu.RUnlock()
	for _, v := range handles  {
		if  v.Name() == name  {
			v.Enable()
			return
		}
	}
}

func  (hr *HandlerRegister) Enable(key int, name string) error {
	err, handles :=  s.Get(key)
	if err != nil {
		return err
	}
	hr.mu.RLock()
	defer hr.mu.RUnlock()
	for _, v := range handles  {
		if  v.Name() == name  {
			v.Disable()
			return
		}
	}
}

