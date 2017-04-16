package replier
import (
    "github.com/songtianyi/wechat-go/wxweb"
)

func Register(session *wxweb.Session) {
	session.HandlerRegister.Add(wxweb.MSG_TEXT, wxweb.Handler(AutoReply))

}
func AutoReply(session *wxweb.Session, msg *wxweb.ReceivedMessage) {
	session.SendText("暂时不在，稍后回复", session.Bot.UserName, msg.FromUserName)
}

