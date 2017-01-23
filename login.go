package wxbot

import (
	"github.com/songtianyi/wechat-go/wxweb"
	"github.com/songtianyi/rrframework/logs"
	"github.com/songtianyi/rrframework/config"
	"fmt"
	"time"
	"net/http"
	"github.com/mdp/qrterminal"
	"os"
)

var (
	WxWebDefaultCommon *wxweb.Common
	WxWebXcg *wxweb.XmlConfig
	Cookies []*http.Cookie
	SynKeyList *wxweb.SyncKeyList
	Bot *wxweb.User
)

func init() {
	WxWebDefaultCommon = &wxweb.Common {
		AppId: "wx782c26e4c19acffb",
		LoginUrl: "https://login.weixin.qq.com",
		Lang: "zh_CN",
		DeviceID: "e" + wxweb.GetRandomStringFromNum(15),
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

func AutoLogin() {
	logs.Debug("%v", WxWebDefaultCommon)
	uuid, err := wxweb.JsLogin(WxWebDefaultCommon)
	if err != nil {
		panic(err)
	}
	logs.Debug(uuid)
	qrterminal.Generate("https://login.weixin.qq.com/l/" + uuid, qrterminal.L, os.Stdout)

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
			}else {
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
}
