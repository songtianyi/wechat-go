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

package cleaner

import (
	"fmt"
	"github.com/songtianyi/rrframework/logs"
	"github.com/songtianyi/wechat-go/wxweb"
	"math"
	"strings"
)

// Register plugin
func Register(session *wxweb.Session) {
	session.HandlerRegister.Add(wxweb.MSG_TEXT, wxweb.Handler(listenCmd), "cleaner")
	if err := session.HandlerRegister.EnableByName("cleaner"); err != nil {
		logs.Error(err)
	}
}

func listenCmd(session *wxweb.Session, msg *wxweb.ReceivedMessage) {
	// from myself
	if msg.FromUserName != session.Bot.UserName {
		return
	}
	// command filter
	if !strings.Contains(strings.ToLower(msg.Content), "run cleaner") {
		return
	}
	// do clean
	users := session.Cm.GetStrangers()
	for _, v := range users {
		fmt.Println(v.NickName)
	}
	max := 25
	for i := range users {
		if i%max == 0 {
			b, err := wxweb.WebWxCreateChatroom(session.WxWebCommon, session.WxWebXcg, session.Cookies, users[i:int(math.Max(float64(i+max), float64(len(users))))], "test")
			if err != nil {
				logs.Error(err)
				return
			}
			logs.Debug(string(b.([]byte)))
		}
	}
}
