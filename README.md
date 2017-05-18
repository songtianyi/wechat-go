## wechat-go

[![Build Status](https://travis-ci.org/songtianyi/wechat-go.svg?branch=master)](https://travis-ci.org/songtianyi/wechat-go)
[![Join the chat at https://gitter.im/wechat-go/Lobby](https://badges.gitter.im/wechat-go/Lobby.svg)](https://gitter.im/wechat-go/Lobby?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge&utm_content=badge)
[![Go Report Card](https://goreportcard.com/badge/github.com/songtianyi/wechat-go)](https://goreportcard.com/report/github.com/songtianyi/wechat-go)
[![codebeat badge](https://codebeat.co/badges/4f78bcb2-bf75-477d-a8f4-b09fde3dae80)](https://codebeat.co/projects/github-com-songtianyi-wechat-go-master)
[![license](https://img.shields.io/github/license/mashape/apistatus.svg)]()

go version wechat web api

* 支持多用户(多开)
* 支持掉线后免扫码重登
* 功能以插件的形式提供，可以根据用户(比如付费情况）选择加载或者不加载某插件
* 对于加载的插件, 使用机器人的用户可以通过微信聊天界面动态开启/关闭.
* 目前已提供头像识别, gif搜索, 笑话大全, "阅后即焚", 消息跨群转发等多个有趣插件
* 插件编写简单, 二次开发极为方便
* 可以发送图片/文字/gif/视频/表情等多种消息


## 安装
	mkdir $GOPATH/src/golang.org/x
	cd $GOPATH/src/golang.org/x
	git clone https://github.com/golang/net.git

	go get -u -v github.com/songtianyi/wechat-go

## 示例项目
[go-aida](https://www.github.com/songtianyi/go-aida)

## 示例代码
##### 基于wechat-go创建自定义的机器人
```go
package main

import (
	"github.com/songtianyi/rrframework/logs"
	"github.com/songtianyi/wechat-go/plugins/wxweb/faceplusplus"
	"github.com/songtianyi/wechat-go/plugins/wxweb/gifer"
	"github.com/songtianyi/wechat-go/plugins/wxweb/replier"
	"github.com/songtianyi/wechat-go/plugins/wxweb/switcher"
	"github.com/songtianyi/wechat-go/wxweb"
)

func main() {
	// 创建session, 一个session对应一个机器人
	session, err := wxweb.CreateSession(nil, nil, wxweb.TERMINAL_MODE)
	if err != nil {
		logs.Error(err)
		return
	}

	// 注册插件, 所有插件默认是关闭的
	faceplusplus.Register(session)
	replier.Register(session)
	switcher.Register(session)
	gifer.Register(session)

	// 开启一些常用插件
	session.HandlerRegister.EnableByName("switcher")
	session.HandlerRegister.EnableByName("faceplusplus")

	// 登录并接收消息
	if err := session.LoginAndServe(false); err != nil {
		logs.Error("session exit, %s", err)
	}
}
```
## 插件列表
###### switcher
一个管理插件的插件
```
#关闭某个插件, 在微信聊天窗口输入
disable faceplusplus
#开启某个插件, 在微信聊天窗口输入
enable faceplusplus
#查看所有插件信息, 在微信聊天窗口输入
dump
```
###### faceplusplus
对收到的图片做面部识别，返回性别和年龄
###### gifer
以收到的文字消息为关键字做gif搜索，返回gif图, 注意返回的gif可能尺度较大，比如文字消息中包含“污”等关键词。
###### replier
对收到的文字/图片消息，做自动应答，回复固定文字消息
###### laosj
随机获取一张美女图片, 在聊天窗口输入
```
美女
```
###### joker
获取一则笑话, 在聊天窗口输入
```
笑话
```
###### revoker
消息撤回插件, 3s后自动撤回手机端所发的文本消息. 机器人发出的消息需要自己在对应插件里写撤回逻辑.

###### system
处理消息撤回/红包等系统提示

###### forwarder
消息跨群转发, 在插件里修改群名的全拼即可.

###### youdao
中英互译插件, 基于有道的翻译API

## 制作自己的插件
```go
package demo // 以插件名命令包名

import (
	"github.com/songtianyi/rrframework/logs" // 导入日志包
	"github.com/songtianyi/wechat-go/wxweb"  // 导入协议包
)

// 必须有的插件注册函数
// 指定session, 可以对不同用户注册不同插件
func Register(session *wxweb.Session) {
	// 将插件注册到session
	// 第一个参数: 指定消息类型, 所有该类型的消息都会被转发到此插件
	// 第二个参数: 指定消息处理函数, 消息会进入此函数
	// 第三个参数: 自定义插件名，不能重名，switcher插件会用到此名称
	session.HandlerRegister.Add(wxweb.MSG_TEXT, wxweb.Handler(demo), "textdemo")

	// 可以多个消息类型使用同一个处理函数，也可以分开
	session.HandlerRegister.Add(wxweb.MSG_IMG, wxweb.Handler(demo), "imgdemo")
}

// 消息处理函数
func demo(session *wxweb.Session, msg *wxweb.ReceivedMessage) {

	// 可选: 可以用contact manager来过滤, 比如过滤掉没有保存到通讯录的群
	contact := session.Cm.GetContactByUserName(msg.FromUserName)
	if contact == nil {
		logs.Error("ignore the messages from", msg.FromUserName)
		return
	}

	// 可选: 根据消息类型来过滤
	if msg.MsgType == wxweb.MSG_IMG {
		return
	}

	// 可选: 根据wxweb.User数据结构中的数据来过滤
	if contact.PYQuanPin != "songtianyi" {
		// 比如根据用户昵称的拼音全拼来过滤
		return
	}

	// 可选: 过滤和自己无关的群组消息
	if msg.IsGroup && msg.Who != session.Bot.UserName {
		return
	}

	// 取出收到的内容
	// 取text
	logs.Info(msg.Content)
	//// 取img
	//if b, err := session.GetImg(msg.MsgId); err == nil {
	//	logs.Debug(string(b))
	//}

	// anything

	// 回复消息
	// 第一个参数: 回复的内容
	// 第二个参数: 机器人ID
	// 第三个参数: 联系人/群组/特殊账号ID
	session.SendText("plugin demo", session.Bot.UserName, wxweb.RealTargetUserName(session, msg))
	// 回复图片和gif 参见wxweb/session.go

}
```
## 展示
![example](http://p1.bpimg.com/567571/374325070b2a9042.jpg)
