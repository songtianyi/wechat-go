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

package youdao // 以插件名命令包名

import (
	"github.com/songtianyi/rrframework/config"
	"github.com/songtianyi/rrframework/logs" // 导入日志包
	"github.com/songtianyi/wechat-go/wxweb"  // 导入协议包
	"io/ioutil"
	"net/http"
	"net/url"
)

// Register plugin
// 必须有的插件注册函数
func Register(session *wxweb.Session) {
	// 将插件注册到session
	session.HandlerRegister.Add(wxweb.MSG_TEXT, wxweb.Handler(youdao), "youdao")
	if err := session.HandlerRegister.EnableByName("youdao"); err != nil {
		logs.Error(err)
	}
}

// 消息处理函数
func youdao(session *wxweb.Session, msg *wxweb.ReceivedMessage) {

	// 可选:避免此插件对未保存到通讯录的群生效 可以用contact manager来过滤
	contact := session.Cm.GetContactByUserName(msg.FromUserName)
	if contact == nil {
		logs.Error("ignore the messages from", msg.FromUserName)
		return
	}

	uri := "http://fanyi.youdao.com/openapi.do?keyfrom=go-aida&key=145986666&type=data&doctype=json&version=1.1&q=" + url.QueryEscape(msg.Content)
	response, err := http.Get(uri)
	if err != nil {
		logs.Error(err)
		return
	}
	defer response.Body.Close()
	body, _ := ioutil.ReadAll(response.Body)
	jc, err := rrconfig.LoadJsonConfigFromBytes(body)
	if err != nil {
		logs.Error(err)
		return
	}
	errorCode, err := jc.GetInt("errorCode")
	if err != nil {
		logs.Error(err)
		return
	}
	if errorCode != 0 {
		logs.Error("youdao API", errorCode)
		return
	}
	trans, err := jc.GetSliceString("translation")
	if err != nil {
		logs.Error(err)
		return
	}
	if len(trans) < 1 {
		return
	}

	session.SendText(msg.At+trans[0], session.Bot.UserName, wxweb.RealTargetUserName(session, msg))
}
