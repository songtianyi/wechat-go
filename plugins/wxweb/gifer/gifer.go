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

func Register(session *wxweb.Session) {
	session.HandlerRegister.Add(wxweb.MSG_TEXT, wxweb.Handler(Gifer), "gifer")
}

func Gifer(session *wxweb.Session, msg *wxweb.ReceivedMessage) {
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
		session.SendEmotionWithBytes(body, session.Bot.UserName, msg.ToUserName)
	} else {
		session.SendEmotionWithBytes(body, session.Bot.UserName, msg.FromUserName)
	}
}
