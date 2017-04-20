package demo // 以插件名命令包名

import (
	"github.com/songtianyi/rrframework/logs" // 导入日志包
	"github.com/songtianyi/wechat-go/wxweb"  // 导入协议包
)

// 必须有的插件注册函数
// 指定session, 可以对不同用户注册不同插件
func Register(session *wxweb.Session) {
	// 将插件注册到session
	// 第一个参数: 指定消息类型, 所有该类型的消息都会被转发到此插件
	// 第二个参数: 指定消息处理函数, 消息会进入此函数
	// 第三个参数: 自定义插件名，不能重名，switcher插件会用到此名称
	session.HandlerRegister.Add(wxweb.MSG_TEXT, wxweb.Handler(demo), "demo")

	// 可以两个消息类型使用同一个处理函数，也可以分开
	session.HandlerRegister.Add(wxweb.MSG_IMG, wxweb.Handler(demo), "demo")
}

// 消息处理函数
func demo(session *wxweb.Session, msg *wxweb.ReceivedMessage) {

	// 可选:避免此插件对所有群/联系人生效 可以用contact manager来过滤
	contact := session.Cm.GetContactByUserName(msg.FromUserName)
	if contact == nil {
		logs.Error("ignore the messages from", msg.FromUserName)
		return
	}

	// 可选: 过滤消息类型
	if msg.MsgType == wxweb.MSG_IMG {
		return
	}

	// 可选: 根绝wxweb.User数据结构结构中的数据来过滤
	if contact.PYQuanPin != "songtianyi" {
		// 根据用户昵称的拼音全拼来过滤
		return
	}

	// 可选:过滤和自己无关的群组消息
	if msg.IsGroup && msg.At != session.Bot.UserName {
		return
	}

	// 取出收到的内容
	// 取text
	logs.Info(msg.Content)
	//// 取img
	//if b, err := session.GetImg(msg.MsgId); err == nil {
	//	logs.Debug(string(b))
	//}

	// anything

	// 回复消息
	// 第一个参数: 回复的内容
	// 第二个参数: 机器人ID
	// 第三个参数: 联系人/群组/特殊账号ID
	session.SendText("plugin demo", session.Bot.UserName, wxweb.RealTargetUserName(session, msg))

}
