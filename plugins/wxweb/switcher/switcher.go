package switcher

import (
	"github.com/songtianyi/rrframework/logs"
	"github.com/songtianyi/wechat-go/wxweb"
	"strings"
)

func Register(session *wxweb.Session) {
	session.HandlerRegister.Add(wxweb.MSG_TEXT, wxweb.Handler(Switcher), "switcher")
}

func Switcher(session *wxweb.Session, msg *wxweb.ReceivedMessage) {
	// contact filter
	contact := session.Cm.GetContactByUserName(msg.FromUserName)
	if contact == nil {
		logs.Error("no this contact, ignore", msg.FromUserName)
		return
	}

	if strings.ToLower(msg.Content) == "dump" {
		session.SendText(session.HandlerRegister.Dump(), session.Bot.UserName, wxweb.RealTargetUserName(session, msg))
		return
	}

	if !strings.Contains(strings.ToLower(msg.Content), "enable") &&
		!strings.Contains(strings.ToLower(msg.Content), "disable") {
		return
	}

	ss := strings.Split(msg.Content, " ")
	if len(ss) < 2 {
		return
	}
	if strings.ToLower(ss[1]) == "switcher" {
		session.SendText("hehe, you think too much", session.Bot.UserName, wxweb.RealTargetUserName(session, msg))
		return
	}
	if strings.ToLower(ss[0]) == "enable" {
		session.HandlerRegister.EnableByName(ss[1])
	} else if strings.ToLower(ss[0]) == "disable" {
		session.HandlerRegister.DisableByName(ss[1])
	}
	if session.Bot.UserName == msg.FromUserName {
		session.SendText(msg.Content+" [DONE]", session.Bot.UserName, msg.ToUserName)
	} else {
		session.SendText(msg.Content+" [DONE]", session.Bot.UserName, msg.FromUserName)
	}
}
