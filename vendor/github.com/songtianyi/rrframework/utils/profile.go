package rrutils

import (
	"github.com/songtianyi/rrframework/logs"
	"net/http"
	_ "net/http/pprof"
)

func StartProfiling() {
	go func() {
		if err := http.ListenAndServe("localhost:6060", nil); err != nil {
			logs.Error("Start profiling fail, %s", err)
			return
		}
	}()
}
