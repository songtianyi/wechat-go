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
)

type handler func(map[string]interface)

type HandlerWrapper struct {
	handle handler
}

func(s *HandlerWrapper) Run(msg map[string]interface) {
	s.handle(msg)
}

var (
	
)

func init() {
	createHandlerRegister()
}

type handlerRegister struct {
	mu   *sync.RWMutex
	hmap map[string]*HandlerWrapper
}

func createHandlerRegister() (error, *HandlerRegister) {
	return nil, &handlerRegister{
		mu:   new(sync.RWMutex),
		hmap: make(map[int][]*HandlerWrapper),
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
