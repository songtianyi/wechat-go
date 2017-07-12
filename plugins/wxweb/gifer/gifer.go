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

package gifer

import (
	"github.com/songtianyi/laosj/spider"
	"github.com/songtianyi/rrframework/logs"
	"github.com/songtianyi/wechat-go/wxweb"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"time"
)

// Register plugin
func Register(session *wxweb.Session) {
	session.HandlerRegister.Add(wxweb.MSG_TEXT, wxweb.Handler(gifer), "gifer")
	if err := session.HandlerRegister.EnableByName("gifer"); err != nil {
		logs.Error(err)
	}
}

func gifer(session *wxweb.Session, msg *wxweb.ReceivedMessage) {
	// contact filter
	contact := session.Cm.GetContactByUserName(msg.FromUserName)
	if contact == nil {
		logs.Error("no this contact, ignore", msg.FromUserName)
		return
	}
	uri := "http://www.gifmiao.com/search/" + url.QueryEscape(msg.Content) + "/3"
	s, err := spider.CreateSpiderFromUrl(uri)
	if err != nil {
		logs.Error(err)
		return
	}
	srcs, _ := s.GetAttr("div.wrap>div#main>ul#waterfall>li.item>div.img_block>a>img.gifImg", "xgif")
	if len(srcs) < 1 {
		logs.Error("no gif result for", msg.Content)
		return
	}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	gif := srcs[r.Intn(len(srcs))]
	resp, err := http.Get(gif)
	if err != nil {
		logs.Error(err)
		return
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	if msg.FromUserName == session.Bot.UserName {
		session.SendEmotionFromBytes(body, session.Bot.UserName, msg.ToUserName)
	} else {
		session.SendEmotionFromBytes(body, session.Bot.UserName, msg.FromUserName)
	}
}
