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

package wxbot

import (
	"fmt"
	"sync"
)

type Handler func(map[string]interface{})

type HandlerWrapper struct {
	handle Handler
}

func (s *HandlerWrapper) Run(msg map[string]interface{}) {
	s.handle(msg)
}

var (
	HandlerRegister *handlerRegister
)

func init() {
	HandlerRegister, _ = createHandlerRegister()
}

type handlerRegister struct {
	mu   *sync.RWMutex
	hmap map[int][]*HandlerWrapper
}

func createHandlerRegister() (*handlerRegister, error) {
	return &handlerRegister{
		mu:   new(sync.RWMutex),
		hmap: make(map[int][]*HandlerWrapper),
	}, nil
}

func (hr *handlerRegister) Add(key int, h Handler) {
	hr.mu.Lock()
	defer hr.mu.Unlock()
	if _, ok := hr.hmap[key]; !ok {
		hr.hmap[key] = make([]*HandlerWrapper, 0)
	}
	hr.hmap[key] = append(hr.hmap[key], &HandlerWrapper{handle: h})
}

func (hr *handlerRegister) Get(key int) (error, []*HandlerWrapper) {
	hr.mu.RLock()
	defer hr.mu.RUnlock()
	if v, ok := hr.hmap[key]; ok {
		return nil, v
	}
	return fmt.Errorf("handlers for key [%d] not registered", key), nil
}
