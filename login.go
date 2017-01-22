package wxbot

import (
	"github.com/songtianyi/wechat-go/wxweb"
	"github.com/songtianyi/rrframework/logs"
)

var (
	WxWebDefaultCommon *wxweb.Common
	wxWebXcg *wxweb.XmlConfig
	cookies []*http.Cookie
	SynList *wxweb.SyncKeyList
)

func int() {
	WxWebDefaultCommon = &wxweb.Common {
		AppId: "wx782c26e4c19acffb",
		LoginUrl: "https://login.weixin.qq.com",
		Lang: "zn_CN",
		DeviceID: "e" + GetRandomStringFromNum(15),
		UserAgent: "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_11_3) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/48.0.2564.109 Safari/537.36",
		CgiUrl: "https://wx.qq.com/cgi-bin/mmwebwx-bin",
		SyncSrvs: []string{
			"webpush.wx.qq.com",
			"webpush.weixin.qq.com",
			"webpush.wechat.com",
			"webpush1.wechat.com",
			"webpush2.wechat.com",
		},
	}

	WxWebXcg = &wxweb.XmlConfig{}
}

func AutoLogin() error {
	uuid, err := wxweb.JsLogin(WxWebDefaultCommon)
	if err != nil {
		panic(err)
	}

	qrcb, err := QrCode(WxWebDefaultCommon, uuid)

	se := rrstorage.CreateUfileStorage("j+4uUJbKZVVa39dGyIi7CxcFbcJ+F8I2Np7pGuvQksNbL2Bu", "6234fc01f795f7e4be705ec0e7ae9d898fcf35c6", "public-songtianyi", 2)
	if err := se.Save(qrcb, uuid+".jpg"); err != nil {
		panic(err)
	}

	redirectUri := ""
loop1:
	for {
		select {
		case <-time.After(5 * time.Second):
			redirectUri, err := wxweb.Login(WxWebDefaultCommon, uuid, 0)
			if err != nil {
				logs.Error(err)
			}
			break loop1
		}
	}

	if err := wxweb.WebNewLoginPage(wxWebXcg, cookies, redirectUri); err != nil {
		panic(err)
	}

	jb, err := wxweb.WebWxInit(wxWebXcg)
	if err != nil {
		panic(err)
	}

	jc, err := rrconfig.LoadJsonConfigFromBytes(jb)
	if err != nil {
		panic(err)
	}
	SynList, err := wxweb.GetSyncKeyListFromJc(jc)
	if err != nil {
		panic(err)
	}
}
