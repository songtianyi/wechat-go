package replier

import (
	"github.com/songtianyi/rrframework/logs"
	"github.com/songtianyi/wechat-go/wxweb"
)

func Register(session *wxweb.Session) {
	session.HandlerRegister.Add(wxweb.MSG_TEXT, wxweb.Handler(AutoReply), "text_repiler")
	if err := session.HandlerRegister.Add(wxweb.MSG_IMG, wxweb.Handler(AutoReply), "img_repiler"); err != nil {
		logs.Error(err)
	}

}
func AutoReply(session *wxweb.Session, msg *wxweb.ReceivedMessage) {
	if !msg.IsGroup {
		session.SendText("暂时不在，稍后回复", session.Bot.UserName, msg.FromUserName)
	}
}
