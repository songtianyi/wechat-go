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

package verify // auto verify friend request

import (
	"fmt"
	"github.com/songtianyi/rrframework/logs" // 导入日志包
	"github.com/songtianyi/wechat-go/wxweb"  // 导入协议包
)

// Register plugin
// 必须有的插件注册函数
// 指定session, 可以对不同用户注册不同插件
func Register(session *wxweb.Session) {
	// 将插件注册到session
	// 第一个参数: 指定消息类型, 所有该类型的消息都会被转发到此插件
	// 第二个参数: 指定消息处理函数, 消息会进入此函数
	// 第三个参数: 自定义插件名，不能重名，switcher插件会用到此名称
	session.HandlerRegister.Add(wxweb.MSG_FV, wxweb.Handler(verify), "verify")

	if err := session.HandlerRegister.EnableByName("verify"); err != nil {
		logs.Error(err)
	}
}

// 消息处理函数
func verify(session *wxweb.Session, msg *wxweb.ReceivedMessage) {

	logs.Info(msg.Content)

	master := session.Cm.GetContactByPYQuanPin("SONGTIANYI")

	if err := session.AcceptFriend("", []*wxweb.VerifyUser{{Value: msg.RecommendInfo.UserName, VerifyUserTicket: msg.RecommendInfo.Ticket}}); err != nil {
		errMsg := fmt.Sprintf("accept %s's friend request error, %s", msg.RecommendInfo.NickName, err.Error())
		logs.Error(errMsg)
		if master != nil {
			session.SendText(errMsg, session.Bot.UserName, master.UserName)
		}
		return
	}

	// 回复消息
	// 第一个参数: 回复的内容
	// 第二个参数: 机器人ID
	// 第三个参数: 联系人/群组/特殊账号ID
	if master != nil {
		session.SendText(fmt.Sprintf("%s accepted", msg.RecommendInfo.NickName), session.Bot.UserName, master.UserName)
	}

}
