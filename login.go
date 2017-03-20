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

package wxbot

import (
	"fmt"
	"github.com/mdp/qrterminal"
	"github.com/songtianyi/rrframework/config"
	"github.com/songtianyi/rrframework/logs"
	"github.com/songtianyi/rrframework/storage"
	"github.com/songtianyi/wechat-go/wxweb"
	"os"
	"time"
)

func AutoLogin(session *WxWebSession) {
	logs.Debug("%v", session.WxWebDefaultCommon)
	uuid, err := wxweb.JsLogin(session.WxWebDefaultCommon)
	if err != nil {
		panic(err)
	}
	logs.Debug(uuid)
	qrterminal.Generate("https://login.weixin.qq.com/l/"+uuid, qrterminal.L, os.Stdout)

	qrcb, err := wxweb.QrCode(session.WxWebDefaultCommon, uuid)
	if err != nil {
		panic(err)
	}
	ls := rrstorage.CreateLocalDiskStorage("./")
	if err := ls.Save(qrcb, "qrcode.jpg"); err != nil {
		panic(err)
	}

	redirectUri := ""
loop1:
	for {
		select {
		case <-time.After(5 * time.Second):
			redirectUri, err = wxweb.Login(session.WxWebDefaultCommon, uuid, "0")
			if err != nil {
				logs.Error(err)
			} else {
				break loop1
			}
		}
	}
	logs.Debug(redirectUri)

	if session.Cookies, err = wxweb.WebNewLoginPage(session.WxWebDefaultCommon, session.WxWebXcg, redirectUri); err != nil {
		panic(err)
	}

	jb, err := wxweb.WebWxInit(session.WxWebDefaultCommon, session.WxWebXcg)
	if err != nil {
		panic(err)
	}

	jc, err := rrconfig.LoadJsonConfigFromBytes(jb)
	if err != nil {
		panic(err)
	}

	session.SynKeyList, err = wxweb.GetSyncKeyListFromJc(jc)
	if err != nil {
		panic(err)
	}
	session.Bot, _ = wxweb.GetUserInfoFromJc(jc)
	logs.Debug(session.Bot)
	ret, err := wxweb.WebWxStatusNotify(session.WxWebDefaultCommon, session.WxWebXcg, session.Bot)
	if err != nil {
		panic(err)
	}
	if ret != 0 {
		panic(fmt.Errorf("WebWxStatusNotify fail, %d", ret))
	}

	cb, err := wxweb.WebWxGetContact(session.WxWebDefaultCommon, session.WxWebXcg, session.Cookies)
	if err != nil {
		panic(err)
	}
	session.Cm, err = CreateContactManagerFromBytes(cb)
	if err != nil {
		panic(err)
	}

	_, err = wxweb.WebWxBatchGetContact(session.WxWebDefaultCommon, session.WxWebXcg, session.Cookies, session.Cm.GetGroupContact())
	if err != nil {
		panic(err)
	}
}
