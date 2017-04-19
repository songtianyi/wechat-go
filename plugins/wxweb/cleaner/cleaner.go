package cleaner

import (
	"fmt"
	"github.com/songtianyi/wechat-go/wxweb"
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
}
