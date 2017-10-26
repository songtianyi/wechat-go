package rrhandler

import (
	"time"
)

// handler for requests
type Handler func(interface{}, interface{})

type HandlerWrapper struct {
	handle  Handler
	timeout time.Duration
}

func (h HandlerWrapper) Run(c interface{}, msg interface{}) {
	h.handle(c, msg)
}

// handler for timer jobs
