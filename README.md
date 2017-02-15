## wechat-go
go version wechat web api

## Install
	go get -u -v github.com/songtianyi/wechat-go

## golang.org/x dep install
	mkdir $GOPATH/src/golang.org/x
	cd $GOPATH/src/golang.org/x
	git clone https://github.com/golang/net.git

## Example code
```go
import (
	"github.com/songtianyi/wechat-go"
)
func main() {
	wxbot.AutoLogin()
	wxbot.Run()
}
```

## Show
![example](http://p1.bpimg.com/567571/374325070b2a9042.jpg)
