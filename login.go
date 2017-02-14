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
	"github.com/songtianyi/wechat-go/wxweb"
	"net/http"
	"os"
	"time"
)

var (
	WxWebDefaultCommon *wxweb.Common
	WxWebXcg           *wxweb.XmlConfig
	Cookies            []*http.Cookie
	SynKeyList         *wxweb.SyncKeyList
	Bot                *wxweb.User
	Cm *ContactManager
)

func init() {
	WxWebDefaultCommon = &wxweb.Common{
		AppId:     "wx782c26e4c19acffb",
		LoginUrl:  "https://login.weixin.qq.com",
		Lang:      "zh_CN",
		DeviceID:  "e" + wxweb.GetRandomStringFromNum(15),
		UserAgent: "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_11_3) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/48.0.2564.109 Safari/537.36",
		CgiUrl:    "https://wx.qq.com/cgi-bin/mmwebwx-bin",
		SyncSrvs: []string{
			"webpush.wx.qq.com",
			"webpush.weixin.qq.com",
			"webpush.wechat.com",
			"webpush1.wechat.com",
			"webpush2.wechat.com",
		},
		UploadUrl: "https://file.wx.qq.com/cgi-bin/mmwebwx-bin/webwxuploadmedia?f=json",
		MediaCount: 0,
	}
	WxWebXcg = &wxweb.XmlConfig{}
}

func AutoLogin() {
	logs.Debug("%v", WxWebDefaultCommon)
	uuid, err := wxweb.JsLogin(WxWebDefaultCommon)
	if err != nil {
		panic(err)
	}
	logs.Debug(uuid)
	qrterminal.Generate("https://login.weixin.qq.com/l/"+uuid, qrterminal.L, os.Stdout)

	//qrcb, err := wxweb.QrCode(WxWebDefaultCommon, uuid)
	//if err != nil {
	//	panic(err)
	//}

	redirectUri := ""
loop1:
	for {
		select {
		case <-time.After(5 * time.Second):
			redirectUri, err = wxweb.Login(WxWebDefaultCommon, uuid, "0")
			if err != nil {
				logs.Error(err)
			} else {
				break loop1
			}
		}
	}
	logs.Debug(redirectUri)

	if Cookies, err = wxweb.WebNewLoginPage(WxWebDefaultCommon, WxWebXcg, redirectUri); err != nil {
		panic(err)
	}

	jb, err := wxweb.WebWxInit(WxWebDefaultCommon, WxWebXcg)
	if err != nil {
		panic(err)
	}

	jc, err := rrconfig.LoadJsonConfigFromBytes(jb)
	if err != nil {
		panic(err)
	}

	SynKeyList, err = wxweb.GetSyncKeyListFromJc(jc)
	if err != nil {
		panic(err)
	}
	Bot, _ = wxweb.GetUserInfoFromJc(jc)
	logs.Debug(Bot)
	ret, err := wxweb.WebWxStatusNotify(WxWebDefaultCommon, WxWebXcg, Bot)
	if err != nil {
		panic(err)
	}
	if ret != 0 {
		panic(fmt.Errorf("WebWxStatusNotify fail, %d", ret))
	}

	cb, err := wxweb.WebWxGetContact(WxWebDefaultCommon, WxWebXcg, Cookies)
	if err != nil {
		panic(err)
	}
	Cm, err = LoadContactFromBytes(cb)
	if err != nil {
		panic(err)
	}

	_, err = wxweb.WebWxBatchGetContact(WxWebDefaultCommon, WxWebXcg, Cookies, Cm.GetGroupContact())
	if err != nil {
		panic(err)
	}
}
