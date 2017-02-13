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

package wxweb

import (
	"encoding/xml"
	"strconv"
	"strings"
)

type Common struct {
	AppId     string
	LoginUrl  string
	Lang      string
	DeviceID  string
	UserAgent string
	CgiUrl    string
	SyncSrvs  []string
	UploadUrl string
	MediaCount uint32
}

type InitReqBody struct {
	BaseRequest  *BaseRequest
	Msg          interface{}
	SyncKey      *SyncKeyList
	rr           int
	Code         int
	FromUserName string
	ToUserName   string
	ClientMsgId  int
	ClientMediaId int
	TotalLen int
	StartPos int
	DataLen int
	MediaType int
	Scene int
}

type BaseRequest struct {
	Uin      string
	Sid      string
	Skey     string
	DeviceID string
}

type XmlConfig struct {
	XMLName     xml.Name `xml:"error"`
	Ret         int      `xml:"ret"`
	Message     string   `xml:"message"`
	Skey        string   `xml:"skey"`
	Wxsid       string   `xml:"wxsid"`
	Wxuin       string   `xml:"wxuin"`
	PassTicket  string   `xml:"pass_ticket"`
	IsGrayscale int      `xml:"isgrayscale"`
}

type SyncKey struct {
	Key int
	Val int
}

type SyncKeyList struct {
	Count int
	List  []SyncKey
}

func (s *SyncKeyList) String() string {
	strs := make([]string, 0)
	for _, v := range s.List {
		strs = append(strs, strconv.Itoa(v.Key)+"_"+strconv.Itoa(v.Val))
	}
	return strings.Join(strs, "|")
}

type User struct {
	Uin               int
	UserName          string
	NickName          string
	HeadImgUrl        string
	RemarkName        string
	PYInitial         string
	PYQuanPin         string
	RemarkPYInitial   string
	RemarkPYQuanPin   string
	HideInputBarFlag  int
	StarFriend        int
	Sex               int
	Signature         string
	AppAccountFlag    int
	VerifyFlag        int
	ContactFlag       int
	WebWxPluginSwitch int
	HeadImgFlag       int
	SnsFlag           int
}

type TextMessage struct {
	Type         int
	Content      string
	FromUserName string
	ToUserName   string
	LocalID      int
	ClientMsgId  int
}

type MediaMessage struct {
	Type         int
	Content      string
	FromUserName string
	ToUserName   string
	LocalID      int
	ClientMsgId  int
	MediaId string
}
