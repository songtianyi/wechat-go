package main

import (
	"log"

	"github.com/songtianyi/wechat-go/plugins/wxweb/replier"
	"github.com/songtianyi/wechat-go/wxweb"
)

func main() {
	session, err := wxweb.CreateSession(nil, nil, wxweb.TERMINAL_MODE)
	if err != nil {
		log.Fatal(err)
		return
	}
	replier.Register(session)

	if err := session.LoginAndServe(false); err != nil {
		log.Fatal("session exit, %s", err)
	}
}
