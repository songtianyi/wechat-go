package laosj

import (
	"github.com/songtianyi/laosj/spider"
	"github.com/songtianyi/rrframework/logs"
	"github.com/songtianyi/wechat-go/wxweb"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strings"
	"time"
)

func Register(session *wxweb.Session) {
	session.HandlerRegister.Add(wxweb.MSG_TEXT, wxweb.Handler(ListenCmd), "laosj")
}

func ListenCmd(session *wxweb.Session, msg *wxweb.ReceivedMessage) {
	// contact filter
	contact := session.Cm.GetContactByUserName(msg.FromUserName)
	if contact == nil {
		logs.Error("no this contact, ignore", msg.FromUserName)
		return
	}
	if !strings.Contains(msg.Content, "美女") &&
		!strings.Contains(strings.ToLower(msg.Content), "sexy") {
		return
	}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	uri := "http://www.mzitu.com/zipai/"
	s, err := spider.CreateSpiderFromUrl(uri)
	if err != nil {
		logs.Error(err)
		return
	}
	srcs, _ := s.GetAttr("div.main>div.main-content>div.postlist>div>ul>li>div.comment-body>p>img", "src")
	if len(srcs) < 1 {
		logs.Error("cannot get mzitu images")
		return
	}
	img := srcs[r.Intn(len(srcs))]
	resp, err := http.Get(img)
	if err != nil {
		logs.Error(err)
		return
	}
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logs.Error(err)
		return
	}

	session.SendImgFromBytes(b, img, session.Bot.UserName, wxweb.RealTargetUserName(session, msg))
}
