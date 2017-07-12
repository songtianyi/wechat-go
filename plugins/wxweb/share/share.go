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

package share // 以插件名命令包名

import (
	"github.com/songtianyi/rrframework/logs"
	"github.com/songtianyi/wechat-go/wxweb" // 导入协议包
	"strings"
)

// Register plugin
// 必须有的插件注册函数
// 指定session, 可以对不同用户注册不同插件
func Register(session *wxweb.Session) {
	// 将插件注册到session
	// 第一个参数: 指定消息类型, 所有该类型的消息都会被转发到此插件
	// 第二个参数: 指定消息处理函数, 消息会进入此函数
	// 第三个参数: 自定义插件名，不能重名，switcher插件会用到此名称
	session.HandlerRegister.Add(wxweb.MSG_TEXT, wxweb.Handler(share), "share")

	if err := session.HandlerRegister.EnableByName("share"); err != nil {
		logs.Error(err)
	}
}

// 消息处理函数
func share(session *wxweb.Session, msg *wxweb.ReceivedMessage) {

	// 取出收到的内容
	// 取text
	if strings.Contains(msg.Content, "纸牌屋") {
		text := "https://pan.baidu.com/s/1sl4S0nr#list/path=%2F"
		session.SendText("纸牌屋第五季在线观看 无毒无广\n"+text, session.Bot.UserName, wxweb.RealTargetUserName(session, msg))
	}

	// for issue#13 debug
	//who := session.Cm.GetContactByUserName(msg.FromUserName)
	//logs.Debug("who send", who)
	//if msg.IsGroup {
	//	mm, err := wxweb.CreateMemberManagerFromGroupContact(session, who)
	//	if err != nil {
	//		logs.Error(err)
	//		return
	//	}
	//	info := mm.GetContactByUserName(msg.Who)
	//	logs.Debug(info)
	//}
}
