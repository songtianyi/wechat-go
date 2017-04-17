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
[go-aida](https://www.github.com/songtianyi/go-aida)

## Example code for creating your own chatbot
```go
package main

import (
	"github.com/songtianyi/rrframework/logs"
	"github.com/songtianyi/wechat-go/plugins/faceplusplus"
	"github.com/songtianyi/wechat-go/wxweb"
	"github.com/songtianyi/wechat-go/plugins/wxweb/gifer"
	"github.com/songtianyi/wechat-go/plugins/wxweb/replier"
	"github.com/songtianyi/wechat-go/plugins/wxweb/switcher"
)

func main() {
	// create session
	session, err := wxweb.CreateSession(nil, wxweb.TERMINAL_MODE)
	if err != nil {
		logs.Error(err)
		return
	}

	// add plugins for this session, they are disabled by default
	faceplusplus.Register(session)
	replier.Register(session)
	switcher.Register(session)
	gifer.Register(session)

	// enable plugin
	session.HandlerRegister.EnableByName("switcher")
	session.HandlerRegister.EnableByName("faceplusplus")

	if err := session.LoginAndServe(); err != nil {
		logs.Error("session exit, %s", err)
	}
}
```
## Plugins
#### switcher
一个管理插件的插件
```
#关闭某个插件
disable faceplusplus
#开启某个插件
enable faceplusplus
```
#### faceplusplus
对收到的图片做面部识别，返回性别和年龄
#### gifer
以收到的文字消息做gif搜索，返回gif图, 注意返回的gif可能尺度较大，比如文字消息中包含“污”等关键词。
#### replier
对收到的文字/图片消息，做自动应答，回复固定文字消息

## Show
![example](http://p1.bpimg.com/567571/374325070b2a9042.jpg)
