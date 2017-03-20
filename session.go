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
	"github.com/songtianyi/rrframework/logs"
	"github.com/songtianyi/wechat-go/wxweb"
	"io/ioutil"
	"net/http"
	"strings"
)

type WxWebSession struct {
	WxWebDefaultCommon *wxweb.Common
	WxWebXcg           *wxweb.XmlConfig
	Cookies            []*http.Cookie
	SynKeyList         *wxweb.SyncKeyList
	Bot                *wxweb.User
	Cm                 *ContactManager
}

func CreateWxWebSession(common *wxweb.Common) *WxWebSession {
	if common == nil {
		common = &wxweb.Common{
			AppId:     "wx782c26e4c19acffb",
			LoginUrl:  "https://login.weixin.qq.com",
			Lang:      "zh_CN",
			DeviceID:  "e" + wxweb.GetRandomStringFromNum(15),
			UserAgent: "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_11_3) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/48.0.2564.109 Safari/537.36",
			CgiUrl:    "https://wx.qq.com/cgi-bin/mmwebwx-bin",
			CgiDomain: "https://wx.qq.com",
			SyncSrvs: []string{
				"webpush.wx.qq.com",
				"webpush.weixin.qq.com",
				"webpush.wechat.com",
				"webpush1.wechat.com",
				"webpush2.wechat.com",
			},
			UploadUrl:   "https://file.wx.qq.com/cgi-bin/mmwebwx-bin/webwxuploadmedia?f=json",
			MediaCount:  0,
			RedirectUri: "https://wx.qq.com/cgi-bin/mmwebwx-bin/webwxnewloginpage",
		}
	}

	wxWebXcg = &wxweb.XmlConfig{}
	return &WxWebSession{
		WxWebDefaultCommon: common,
		WxWebXcg:           wxWebXcg,
	}
}

// send text msg type 1
func (s *WxWebSession) SendText(msg, from, to string) {
	ret, err := wxweb.WebWxSendTextMsg(s.WxWebDefaultCommon, s.WxWebXcg, s.Cookies, from, to, msg)
	if ret != 0 {
		logs.Error(ret, err)
		return
	}
}

// send img, upload then send
func (s *WxWebSession) SendImg(path, from, to string) {
	ss := strings.Split(path, "/")
	b, err := ioutil.ReadFile(path)
	if err != nil {
		logs.Error(err)
		return
	}
	mediaId, err := wxweb.WebWxUploadMedia(s.WxWebDefaultCommon, s.WxWebXcg, s.Cookies, ss[len(ss)-1], b)
	if err != nil {
		logs.Error(err)
		return
	}
	ret, err := wxweb.WebWxSendMsgImg(s.WxWebDefaultCommon, s.WxWebXcg, s.Cookies, from, to, mediaId)
	if err != nil || ret != 0 {
		logs.Error(ret, err)
		return
	}
}

// get img by MsgId
func (s *WxWebSession) GetImg(msgId string) ([]byte, error) {
	return wxweb.WebWxGetMsgImg(s.WxWebDefaultCommon, s.WxWebXcg, s.Cookies, msgId)
}

// send gif, upload then send
func (s *WxWebSession) SendEmotion(path, from, to string) {
	ss := strings.Split(path, "/")
	b, err := ioutil.ReadFile(path)
	if err != nil {
		logs.Error(err)
		return
	}
	mediaId, err := wxweb.WebWxUploadMedia(s.WxWebDefaultCommon, s.WxWebXcg, s.Cookies, ss[len(ss)-1], b)
	if err != nil {
		logs.Error(err)
		return
	}
	ret, err := wxweb.WebWxSendEmoticon(s.WxWebDefaultCommon, s.WxWebXcg, s.Cookies, from, to, mediaId)
	if err != nil || ret != 0 {
		logs.Error(ret, err)
	}
}
