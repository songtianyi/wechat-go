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

package wxweb

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/mdp/qrterminal"
	"github.com/songtianyi/rrframework/config"
	"github.com/songtianyi/rrframework/logs"
	"github.com/songtianyi/rrframework/storage"
)

const (
	// WEB_MODE: in this mode CreateSession will return a QRCode image url
	WEB_MODE = iota + 1
	// MINAL_MODE:  CreateSession will output qrcode in terminal
	TERMINAL_MODE

	BACKGROUND_MODE
	TERMINAL_MODE_GOLAND
)

var (
	// DefaultCommon: default session config
	DefaultCommon = &Common{
		AppId:      "wx782c26e4c19acffb",
		LoginUrl:   "https://login.weixin.qq.com",
		Lang:       "zh_CN",
		DeviceID:   "e" + GetRandomStringFromNum(15),
		UserAgent:  "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_11_3) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/48.0.2564.109 Safari/537.36",
		SyncSrv:    "webpush.wx.qq.com",
		UploadUrl:  "https://file.wx.qq.com/cgi-bin/mmwebwx-bin/webwxuploadmedia?f=json",
		MediaCount: 0,
	}
	URLPool = []UrlGroup{
		{"wx2.qq.com", "file.wx2.qq.com", "webpush.wx2.qq.com"},
		{"wx8.qq.com", "file.wx8.qq.com", "webpush.wx8.qq.com"},
		{"qq.com", "file.wx.qq.com", "webpush.wx.qq.com"},
		{"web2.wechat.com", "file.web2.wechat.com", "webpush.web2.wechat.com"},
		{"wechat.com", "file.web.wechat.com", "webpush.web.wechat.com"},
	}
	LogoutSign = make(map[int]int)
)

func init() {
	LogoutSign[1100] = 1
	LogoutSign[1101] = 1
	LogoutSign[1102] = 1
	LogoutSign[1205] = 1
}

// Session: wechat bot session
type Session struct {
	WxWebCommon     *Common
	WxWebXcg        *XmlConfig
	muCookie        sync.RWMutex
	Cookies         []*http.Cookie
	SynKeyList      *SyncKeyList
	Bot             *User
	Cm              *ContactManager
	QrcodePath      string //qrcode path
	QrcodeUUID      string //uuid
	HandlerRegister *HandlerRegister
	CreateTime      int64
	LastMsgID       string
	Api             *ApiV2
	OnLoginAvatar   func(string) error
	AfterLogin      func() error
}

// CreateSession: create wechat bot session
// if common is nil, session will be created with default config
// if handlerRegister is nil,  session will create a new HandlerRegister
func CreateSession(common *Common, handlerRegister *HandlerRegister, qrmode int) (*Session, error) {
	if common == nil {
		common = DefaultCommon
	}

	wxWebXcg := &XmlConfig{}

	// get qrcode
	api := NewApiV2()
	uuid, err := api.JsLogin(common)
	if err != nil {
		return nil, err
	}
	logs.Info(uuid)
	session := &Session{
		WxWebCommon: common,
		WxWebXcg:    wxWebXcg,
		QrcodeUUID:  uuid,
		Api:         api,
		CreateTime:  time.Now().Unix(),
		OnLoginAvatar: func(string) error {
			return nil
		},
		AfterLogin: func() error {
			return nil
		},
	}

	if handlerRegister != nil {
		session.HandlerRegister = handlerRegister
	} else {
		session.HandlerRegister = CreateHandlerRegister()
	}
	if qrmode == BACKGROUND_MODE {
		logs.Info("https://login.weixin.qq.com/l/" + uuid)
	} else if qrmode == TERMINAL_MODE {
		qrterminal.Generate("https://login.weixin.qq.com/l/"+uuid, qrterminal.L, os.Stdout)

	} else if qrmode == TERMINAL_MODE_GOLAND {

		config := qrterminal.Config{
			Level:     qrterminal.L,
			Writer:    os.Stdout,
			BlackChar: qrterminal.WHITE,
			WhiteChar: qrterminal.BLACK,
			QuietZone: 1,
		}

		qrterminal.GenerateWithConfig("https://login.weixin.qq.com/l/"+uuid, config)
	} else if qrmode == WEB_MODE {
		qrcb, err := api.QrCode(common, uuid)
		if err != nil {
			return nil, err
		}
		ls := rrstorage.CreateLocalDiskStorage("../web/public/qrcode/")
		if err := ls.Save(qrcb, uuid+".jpg"); err != nil {
			return nil, err
		}
		session.QrcodePath = "../web/public/qrcode/" + uuid + ".jpg"
		logs.Info("QrcodePath: %s", session.QrcodePath)
	}
	return session, nil
}
func CreateWebSessionWithPath(common *Common, handlerRegister *HandlerRegister, qrcode_path string) (*Session, error) {
	if common == nil {
		common = DefaultCommon
	}

	wxWebXcg := &XmlConfig{}
	api := NewApiV2()

	// get qrcode
	uuid, err := JsLogin(common)
	if err != nil {
		return nil, err
	}
	logs.Info(uuid)
	session := &Session{
		WxWebCommon: common,
		WxWebXcg:    wxWebXcg,
		QrcodeUUID:  uuid,
		CreateTime:  time.Now().Unix(),
		Api:         api,
	}

	if handlerRegister != nil {
		session.HandlerRegister = handlerRegister
	} else {
		session.HandlerRegister = CreateHandlerRegister()
	}

	qrcb, err := api.QrCode(common, uuid)
	if err != nil {
		return nil, err
	}
	ls := rrstorage.CreateLocalDiskStorage(qrcode_path)
	if err := ls.Save(qrcb, uuid+".jpg"); err != nil {
		return nil, err
	}
	session.QrcodePath = uuid + ".jpg"
	logs.Info("QrcodePath: %s", session.QrcodePath)
	return session, nil
}

func (s *Session) analizeVersion(uri string) {
	u, _ := url.Parse(uri)

	// version may change
	s.WxWebCommon.CgiDomain = u.Scheme + "://" + u.Host
	s.WxWebCommon.CgiUrl = s.WxWebCommon.CgiDomain + "/cgi-bin/mmwebwx-bin"

	for _, urlGroup := range URLPool {
		if strings.Contains(u.Host, urlGroup.IndexUrl) {
			s.WxWebCommon.SyncSrv = urlGroup.SyncUrl
			s.WxWebCommon.UploadUrl = fmt.Sprintf("https://%s/cgi-bin/mmwebwx-bin/webwxuploadmedia?f=json", urlGroup.UploadUrl)
			return
		}
	}
}

func (s *Session) scanWaiter(onAvatar func(string) error) error {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()
	for range ticker.C {
		redirectUri, err := s.Api.Login(s.WxWebCommon, s.QrcodeUUID, "0")
		if err != nil {
			logs.Info(err)
			switch {
			case strings.Contains(err.Error(), "window.code=400"),
				strings.Contains(err.Error(), "window.code=500"),
				strings.Contains(err.Error(), "window.code=0"):
				return err
			case strings.Contains(err.Error(), "window.code=201"):
				avatar, err := GetLoginAvatar(err.Error())
				if err != nil {
					return err
				}
				if err := onAvatar(avatar); err != nil {
					return err
				}
			}
		} else {
			s.WxWebCommon.RedirectUri = redirectUri
			s.analizeVersion(s.WxWebCommon.RedirectUri)
			break
		}
	}
	return nil
}
func (s *Session) SetCookies(cookies []*http.Cookie) {
	s.muCookie.Lock()
	defer s.muCookie.Unlock()
	s.Cookies = cookies
}

// LoginAndServe: login wechat web and enter message receiving loop
func (s *Session) LoginAndServe(useCache bool) error {

	var (
		err error
	)

	if !useCache {
		if s.GetCookies() != nil {
			// confirmWaiter
		}

		if err := s.scanWaiter(s.OnLoginAvatar); err != nil {
			return err
		}
		var cookies []*http.Cookie
		// update cookies
		if cookies, err = s.Api.WebNewLoginPage(s.WxWebCommon, s.WxWebXcg, s.WxWebCommon.RedirectUri); err != nil {
			return err
		}
		s.SetCookies(cookies)

	}

	jb, err := s.Api.WebWxInit(s.WxWebCommon, s.WxWebXcg)
	if err != nil {
		return err
	}

	jc, err := rrconfig.LoadJsonConfigFromBytes(jb)
	if err != nil {
		return err
	}

	s.SynKeyList, err = GetSyncKeyListFromJc(jc)
	if err != nil {
		return err
	}
	s.Bot, _ = GetUserInfoFromJc(jc)
	logs.Info(s.Bot)
	ret, err := s.Api.WebWxStatusNotify(s.WxWebCommon, s.WxWebXcg, s.Bot)
	if err != nil {
		return err
	}
	if ret != 0 {
		return fmt.Errorf("WebWxStatusNotify fail, %d", ret)
	}

	cb, err := s.Api.WebWxGetContact(s.WxWebCommon, s.WxWebXcg, s.GetCookies())
	if err != nil {
		return err
	}

	s.Cm, err = CreateContactManagerFromBytes(cb)
	if err != nil {
		return err
	}

	// for v2
	s.Cm.AddUser(s.Bot)
	s.AfterLogin()

	if err := s.serve(); err != nil {
		return err
	}
	return nil
}

func (s *Session) SetAfterLogin(f func() error) {
	s.AfterLogin = f
}

func (s *Session) SetOnLoginAvatar(f func(avatar string) error) {
	s.OnLoginAvatar = f
}

func (s *Session) serve() error {
	msg := make(chan []byte, 1000)
	// syncheck
	errChan := make(chan error)
	go s.producer(msg, errChan)
	for {
		select {
		case m := <-msg:
			go s.consumer(m)
		case err := <-errChan:
			// TODO maybe not all consumed messages ended
			return err
		}
	}
}
func (s *Session) producer(msg chan []byte, errChan chan error) {
	logs.Info("entering synccheck loop")
loop1:
	for {
		var ret, sel int
		var err error
		for i := 0; i <= 10; i++ {
			ret, sel, err = s.Api.SyncCheck(s.WxWebCommon, s.WxWebXcg, s.GetCookies(), s.WxWebCommon.SyncSrv, s.SynKeyList)
			if err != nil {
				if i >= 10 {
					logs.Error("Err SyncCheck  %s try %d", err.Error(), i)
				} else {
					logs.Info("SyncCheck uin %d tiem %d", s.Bot.Uin, i)
				}
			} else {
				break
			}
		}

		logs.Info(s.WxWebCommon.SyncSrv, ret, sel, s.Bot.Uin) //检查状态返回的值

		if ret == 0 { //0 正常
			// check success
			// new message
			for i := 0; i <= 10; i++ {
				cookies, err := s.Api.WebWxSync(s.WxWebCommon, s.WxWebXcg, s.GetCookies(), msg, s.SynKeyList)
				if err != nil {
					if i >= 10 {
						logs.Error("Err WebWxSync try  %s try %d", err.Error(), i)
					} else {
						logs.Info("WebWxSync uin %d tiem %d", s.Bot.Uin, i)
					}
				} else {
					if cookies != nil {
						s.SetCookies(cookies)
					}
					break
				}
			}

		} else if s.isLogout(ret) { //1100 失败/登出微信

			errChan <- fmt.Errorf("api blocked, ret:%d", ret)
			break loop1
		} else {
			errChan <- fmt.Errorf("unhandled exception ret %d", ret)
			break loop1
		}
	}

}
func (s *Session) isLogout(code int) bool {
	_, has := LogoutSign[code]
	return has
}

func (s *Session) consumer(msg []byte) {
	// analize message
	jc, _ := rrconfig.LoadJsonConfigFromBytes(msg)
	msgCount, _ := jc.GetInt("AddMsgCount")
	if msgCount < 1 {
		// no msg details
		return
	}
	msgis, _ := jc.GetInterfaceSlice("AddMsgList")
	for _, v := range msgis {
		rmsg := s.analize(v.(map[string]interface{}))
		err, handles := s.HandlerRegister.Get(rmsg.MsgType)
		if err != nil {
			logs.Warn(err)
			continue
		}
		for _, v := range handles {
			go v.Run(s, rmsg)
		}
	}
}

func (s *Session) analize(msg map[string]interface{}) *ReceivedMessage {
	rmsg := &ReceivedMessage{
		MsgId:         msg["MsgId"].(string),
		OriginContent: msg["Content"].(string),
		FromUserName:  msg["FromUserName"].(string),
		ToUserName:    msg["ToUserName"].(string),
		MsgType:       int(msg["MsgType"].(float64)),
		SubType:       int(msg["SubMsgType"].(float64)),
		Url:           msg["Url"].(string),
	}

	// friend verify message
	if rmsg.MsgType == MSG_FV {
		riif := msg["RecommendInfo"].(map[string]interface{})
		rmsg.RecommendInfo = &RecommendInfo{
			Ticket:     riif["Ticket"].(string),
			UserName:   riif["UserName"].(string),
			NickName:   riif["NickName"].(string),
			Content:    riif["Content"].(string),
			Sex:        int(riif["Sex"].(float64)),
			Alias:      riif["Alias"].(string),
			AttrStatus: uint32(riif["AttrStatus"].(float64)),
			City:       riif["City"].(string),
			OpCode:     int(riif["OpCode"].(float64)),
			Province:   riif["Province"].(string),
			QQNum:      int(riif["QQNum"].(float64)),
			Scene:      int(riif["Scene"].(float64)),
			Signature:  riif["Signature"].(string),
			VerifyFlag: int(riif["VerifyFlag"].(float64)),
		}
	}

	if strings.Contains(rmsg.FromUserName, "@@") ||
		strings.Contains(rmsg.ToUserName, "@@") {
		rmsg.IsGroup = true
		// group message
		ss := strings.Split(rmsg.OriginContent, ":<br/>")
		if len(ss) > 1 {
			rmsg.Who = ss[0]
			rmsg.Content = ss[1]
		} else {
			rmsg.Who = s.Bot.UserName
			rmsg.Content = rmsg.OriginContent
		}
	} else {
		// none group message
		rmsg.Who = rmsg.FromUserName
		rmsg.Content = rmsg.OriginContent
	}

	if rmsg.MsgType == MSG_TEXT &&
		len(rmsg.Content) > 1 &&
		strings.HasPrefix(rmsg.Content, "@") {
		// @someone
		ss := strings.Split(rmsg.Content, "\u2005")
		if len(ss) == 2 {
			rmsg.At = ss[0] + "\u2005"
			rmsg.Content = ss[1]
		}
	}
	return rmsg
}

// message funcs
func (s *Session) After(duration time.Duration) *Session {
	select {
	case <-time.After(duration):
		return s
	}
}

func (s *Session) At(d time.Time) *Session {
	return s.After(d.Sub(time.Now()))
}

// SendText: send text msg type 1
func (s *Session) SendText(msg, from, to string) (string, string, error) {
	b, err := s.Api.WebWxSendMsg(s.WxWebCommon, s.WxWebXcg, s.GetCookies(), from, to, msg)
	if err != nil {
		return "", "", err
	}
	jc, _ := rrconfig.LoadJsonConfigFromBytes(b)
	ret, _ := jc.GetInt("BaseResponse.Ret")
	if ret != 0 {
		errMsg, _ := jc.GetString("BaseResponse.ErrMsg")
		return "", "", fmt.Errorf("WebWxSendMsg Ret=%d, ErrMsg=%s", ret, errMsg)
	}
	msgID, _ := jc.GetString("MsgID")
	localID, _ := jc.GetString("LocalID")
	return msgID, localID, nil
}

// SendImg: send img, upload then send
func (s *Session) SendImg(path, from, to string) {
	ss := strings.Split(path, "/")
	b, err := ioutil.ReadFile(path)
	if err != nil {
		logs.Error(err)
		return
	}
	mediaId, err := s.Api.WebWxUploadMedia(s.WxWebCommon, s.WxWebXcg, s.GetCookies(), ss[len(ss)-1], b)
	if err != nil {
		logs.Error(err)
		return
	}
	ret, err := s.Api.WebWxSendMsgImg(s.WxWebCommon, s.WxWebXcg, s.GetCookies(), from, to, mediaId)
	if err != nil || ret != 0 {
		logs.Error(ret, err)
		return
	}
}

// SendImgFromBytes: send image from mem
func (s *Session) SendImgFromBytes(b []byte, path, from, to string) {
	ss := strings.Split(path, "/")
	mediaId, err := s.Api.WebWxUploadMedia(s.WxWebCommon, s.WxWebXcg, s.GetCookies(), ss[len(ss)-1], b)
	if err != nil {
		logs.Error(err)
		return
	}
	ret, err := s.Api.WebWxSendMsgImg(s.WxWebCommon, s.WxWebXcg, s.GetCookies(), from, to, mediaId)
	if err != nil || ret != 0 {
		logs.Error(ret, err)
		return
	}
}

// GetImg: get img by MsgId
func (s *Session) GetImg(msgId string) ([]byte, error) {
	return s.Api.WebWxGetMsgImg(s.WxWebCommon, s.WxWebXcg, s.GetCookies(), msgId)
}

// SendEmotionFromPath: send gif, upload then send
func (s *Session) SendEmotionFromPath(path, from, to string) {
	ss := strings.Split(path, "/")
	b, err := ioutil.ReadFile(path)
	if err != nil {
		logs.Error(err)
		return
	}
	mediaId, err := s.Api.WebWxUploadMedia(s.WxWebCommon, s.WxWebXcg, s.GetCookies(), ss[len(ss)-1], b)
	if err != nil {
		logs.Error(err)
		return
	}
	ret, err := s.Api.WebWxSendEmoticon(s.WxWebCommon, s.WxWebXcg, s.GetCookies(), from, to, mediaId)
	if err != nil || ret != 0 {
		logs.Error(ret, err)
	}
}

// SendEmotionFromBytes: send gif/emoji from mem
func (s *Session) SendEmotionFromBytes(b []byte, from, to string) {
	mediaId, err := s.Api.WebWxUploadMedia(s.WxWebCommon, s.WxWebXcg, s.GetCookies(), from+".gif", b)
	if err != nil {
		logs.Error(err)
		return
	}
	ret, err := s.Api.WebWxSendEmoticon(s.WxWebCommon, s.WxWebXcg, s.GetCookies(), from, to, mediaId)
	if err != nil || ret != 0 {
		logs.Error(ret, err)
	}
}

// RevokeMsg: revoke message
func (s *Session) RevokeMsg(clientMsgId, svrMsgId, toUserName string) {
	err := s.Api.WebWxRevokeMsg(s.WxWebCommon, s.WxWebXcg, s.GetCookies(), clientMsgId, svrMsgId, toUserName)
	if err != nil {
		logs.Error("revoke msg %s failed, %s", clientMsgId+":"+svrMsgId, err)
		return
	}
}
func (s *Session) GetCookies() []*http.Cookie {
	s.muCookie.RLock()
	defer s.muCookie.RUnlock()
	return s.Cookies
}

// user funcs
// Logout: logout web wechat
func (s *Session) Logout() error {
	return s.Api.WebWxLogout(s.WxWebCommon, s.WxWebXcg, s.GetCookies())
}

func (s *Session) AcceptFriend(verifyContent string, vul []*VerifyUser) error {
	b, err := s.Api.WebWxVerifyUser(s.WxWebCommon, s.WxWebXcg, s.GetCookies(), 3, verifyContent, vul)
	if err != nil {
		return err
	}
	jc, err := rrconfig.LoadJsonConfigFromBytes(b)
	if err != nil {
		return err
	}
	retcode, err := jc.GetInt("BaseResponse.Ret")
	if err != nil {
		return err
	}
	if retcode != 0 {
		return fmt.Errorf("BaseResponse.Ret %d", retcode)
	}
	return nil
}
