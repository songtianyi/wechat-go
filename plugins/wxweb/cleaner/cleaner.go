package cleaner

import (
	"fmt"
	"github.com/songtianyi/rrframework/logs"
	"github.com/songtianyi/wechat-go/wxweb"
	"math"
	"strings"
)

func Register(session *wxweb.Session) {
	session.HandlerRegister.Add(wxweb.MSG_TEXT, wxweb.Handler(ListenCmd), "cleaner")
}

func ListenCmd(session *wxweb.Session, msg *wxweb.ReceivedMessage) {
	// from myself
	if msg.FromUserName != session.Bot.UserName {
		return
	}
	// command filter
	if !strings.Contains(strings.ToLower(msg.Content), "run cleaner") {
		return
	}
	// do clean
	users := session.Cm.GetStrangers()
	for _, v := range users {
		fmt.Println(v.NickName)
	}
	max := 25
	for i, _ := range users {
		if i%max == 0 {
			b, err := wxweb.WebWxCreateChatroom(session.WxWebCommon, session.WxWebXcg, session.Cookies, users[i:int(math.Max(float64(i+max), float64(len(users))))], "test")
			if err != nil {
				logs.Error(err)
				return
			}
			logs.Debug(string(b.([]byte)))
		}
	}
}
