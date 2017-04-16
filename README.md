## wechat-go
go version wechat web api

## Install
	go get -u -v github.com/songtianyi/wechat-go

## golang.org/x dep install
	mkdir $GOPATH/src/golang.org/x
	cd $GOPATH/src/golang.org/x
	git clone https://github.com/golang/net.git
	git clone https://github.com/golang/text.git

## Example project
	![go-aida](https://www.github.com/songtianyi/go-aida)

## Example code for creating your own chatbot
```go
package main

import (
	"github.com/songtianyi/rrframework/logs"
	"github.com/songtianyi/wechat-go/plugins/faceplusplus"
	"github.com/songtianyi/wechat-go/wxweb"
)

func main() {
	// create session
	session, err := wxweb.CreateSession(nil, wxweb.TERMINAL_MODE)
	if err != nil {
		logs.Error(err)
		return
	}
	// add plugins for this session
	faceplusplus.Register(session)

	if err := session.LoginAndServe(); err != nil {
		logs.Error("session exit, %s", err)
	}
}
```

## Show
![example](http://p1.bpimg.com/567571/374325070b2a9042.jpg)
