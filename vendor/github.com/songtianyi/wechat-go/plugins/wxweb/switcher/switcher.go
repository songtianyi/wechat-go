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

package switcher

import (
	"github.com/songtianyi/rrframework/logs"
	"github.com/songtianyi/wechat-go/wxweb"
	"strings"
)

// Register plugin
func Register(session *wxweb.Session) {
	session.HandlerRegister.Add(wxweb.MSG_TEXT, wxweb.Handler(switcher), "switcher")
	if err := session.HandlerRegister.EnableByName("switcher"); err != nil {
		logs.Error(err)
	}
}

func switcher(session *wxweb.Session, msg *wxweb.ReceivedMessage) {
	// contact filter
	contact := session.Cm.GetContactByUserName(msg.FromUserName)
	if contact == nil {
		logs.Error("no this contact, ignore", msg.FromUserName)
		return
	}

	if strings.ToLower(msg.Content) == "dump" {
		session.SendText(session.HandlerRegister.Dump(), session.Bot.UserName, wxweb.RealTargetUserName(session, msg))
		return
	}

	if !strings.Contains(strings.ToLower(msg.Content), "enable") &&
		!strings.Contains(strings.ToLower(msg.Content), "disable") {
		return
	}

	ss := strings.Split(msg.Content, " ")
	if len(ss) < 2 {
		return
	}
	if strings.ToLower(ss[1]) == "switcher" {
		session.SendText("hehe, you think too much", session.Bot.UserName, wxweb.RealTargetUserName(session, msg))
		return
	}

	var (
		err error
	)
	if strings.ToLower(ss[0]) == "enable" {
		if err = session.HandlerRegister.EnableByName(ss[1]); err == nil {
			session.SendText(msg.Content+" [DONE]", session.Bot.UserName, wxweb.RealTargetUserName(session, msg))
		} else {
			session.SendText(err.Error(), session.Bot.UserName, wxweb.RealTargetUserName(session, msg))
		}
	} else if strings.ToLower(ss[0]) == "disable" {
		if err = session.HandlerRegister.DisableByName(ss[1]); err == nil {
			session.SendText(msg.Content+" [DONE]", session.Bot.UserName, wxweb.RealTargetUserName(session, msg))
		} else {
			session.SendText(err.Error(), session.Bot.UserName, wxweb.RealTargetUserName(session, msg))
		}
	}
}
