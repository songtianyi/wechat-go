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

package wxbot

import (
	"github.com/songtianyi/rrframework/logs"
	"github.com/songtianyi/rrframework/config"
	"github.com/songtianyi/wechat-go/wxweb"
)

func Run() {
	msg := make(chan []byte, 10000)
	// syncheck
	go producer(msg)
	for {
		select {
		case m := <-msg:
			go consumer(m)
		}
	}

}

func producer(msg chan []byte) {
	for {
		for _, v := range WxWebDefaultCommon.SyncSrvs {
			ret, sel, err := wxweb.SyncCheck(WxWebDefaultCommon, WxWebXcg, Cookies, v, SynKeyList)
			logs.Debug(v, ret, sel)
			if err != nil {
				logs.Error(err)
				continue
			}
			if ret == 0 {
				// check success
				if sel == 2 {
					// new message
					err := wxweb.WebWxSync(WxWebDefaultCommon, WxWebXcg, Cookies, msg, SynKeyList)
					if err != nil {
						logs.Error(err)
					}
				}
				break
			}
		}
	}

}

func consumer(msg []byte) {
	// analize message
	jc, _ := rrconfig.LoadJsonConfigFromBytes(msg)
	msgCount, _ := jc.GetInt("AddMsgCount")
	if msgCount < 1 {
		return
	}
	msgis, _ := jc.GetInterfaceSlice("AddMsgList")
	for _, v := range msgis {
		msgi := v.(map[string]interface{})
		msgType := int(msgi["MsgType"].(float64))
		if msgType == 51 {
			continue
		}
		if msgType == 1 {
			_, handles := HandlerRegister.Get(msgType)
			for _, v := range handles {
				v.Run(msgi)
			}
		}
	}
}
