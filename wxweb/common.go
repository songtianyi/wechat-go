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

const (
	// msg types
	MSG_TEXT        = 1     // text message
	MSG_IMG         = 3     // image message
	MSG_VOICE       = 34    // voice message
	MSG_FV          = 37    // friend verification message
	MSG_PF          = 40    // POSSIBLEFRIEND_MSG
	MSG_SCC         = 42    // shared contact card
	MSG_VIDEO       = 43    // video message
	MSG_EMOTION     = 47    // gif
	MSG_LOCATION    = 48    // location message
	MSG_LINK        = 49    // shared link message
	MSG_VOIP        = 50    // VOIPMSG
	MSG_INIT        = 51    // wechat init message
	MSG_VOIPNOTIFY  = 52    // VOIPNOTIFY
	MSG_VOIPINVITE  = 53    // VOIPINVITE
	MSG_SHORT_VIDEO = 62    // short video message
	MSG_SYSNOTICE   = 9999  // SYSNOTICE
	MSG_SYS         = 10000 // system message
	MSG_WITHDRAW    = 10002 // withdraw notification message

)

// Common: session config
type Common struct {
	AppId       string
	LoginUrl    string
	Lang        string
	DeviceID    string
	UserAgent   string
	CgiUrl      string
	CgiDomain   string
	SyncSrv     string
	UploadUrl   string
	MediaCount  uint32
	RedirectUri string
}

type UrlGroup struct {
	IndexUrl  string
	UploadUrl string
	SyncUrl   string
}

// InitReqBody: common http request body struct
type InitReqBody struct {
	BaseRequest        *BaseRequest
	Msg                interface{}
	SyncKey            *SyncKeyList
	rr                 int
	Code               int
	FromUserName       string
	ToUserName         string
	ClientMsgId        int
	ClientMediaId      int
	TotalLen           int
	StartPos           int
	DataLen            int
	MediaType          int
	Scene              int
	Count              int
	List               []*User
	Opcode             int
	SceneList          []int
	SceneListCount     int
	VerifyContent      string
	VerifyUserList     []*VerifyUser
	VerifyUserListSize int
	skey               string
	MemberCount        int
	MemberList         []*User
	Topic              string
}

// RevokeReqBody: revoke message api http request body
type RevokeReqBody struct {
	BaseRequest *BaseRequest
	ClientMsgId string
	SvrMsgId    string
	ToUserName  string
}

// LogoutReqBody: logout api http request body
type LogoutReqBody struct {
	sid string
	uin string
}

// BaseRequest: http request body BaseRequest
type BaseRequest struct {
	Uin      string
	Sid      string
	Skey     string
	DeviceID string
}

// XmlConfig: web api xml response struct
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

// SyncKey: struct for synccheck
type SyncKey struct {
	Key int
	Val int
}

// SyncKeyList: list of synckey
type SyncKeyList struct {
	Count int
	List  []SyncKey
}

// s.String output synckey list in string
func (s *SyncKeyList) String() string {
	strs := make([]string, 0)
	for _, v := range s.List {
		strs = append(strs, strconv.Itoa(v.Key)+"_"+strconv.Itoa(v.Val))
	}
	return strings.Join(strs, "|")
}

// User: contact struct
type User struct {
	Uin               int
	UserName          string
	NickName          string
	HeadImgUrl        string
	ContactFlag       int
	MemberCount       int
	MemberList        []*User
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
	Statues           int
	AttrStatus        uint32
	Province          string
	City              string
	Alias             string
	VerifyFlag        int
	OwnerUin          int
	WebWxPluginSwitch int
	HeadImgFlag       int
	SnsFlag           int
	UniFriend         int
	DisplayName       string
	ChatRoomId        int
	KeyWord           string
	EncryChatRoomId   string
	IsOwner           int
	MemberStatus      int
}

// TextMessage: text message struct
type TextMessage struct {
	Type         int
	Content      string
	FromUserName string
	ToUserName   string
	LocalID      int
	ClientMsgId  int
}

// MediaMessage
type MediaMessage struct {
	Type         int
	Content      string
	FromUserName string
	ToUserName   string
	LocalID      int
	ClientMsgId  int
	MediaId      string
}

// EmotionMessage: gif/emoji message struct
type EmotionMessage struct {
	ClientMsgId  int
	EmojiFlag    int
	FromUserName string
	LocalID      int
	MediaId      string
	ToUserName   string
	Type         int
}

// BaseResponse: web api http response body BaseResponse struct
type BaseResponse struct {
	Ret    int
	ErrMsg string
}

// WxWebGetContactResponse: get contact response struct
type WxWebGetContactResponse struct {
	BaseResponse *BaseResponse
	MemberCount  int
	MemberList   []*User
	Seq          int
}

// WxWebBatchGetContactResponse: batch get contact response struct
type WxWebBatchGetContactResponse struct {
	BaseResponse *BaseResponse
	Count        int
	ContactList  []*User
}

// VerifyUser: verify user request body struct
type VerifyUser struct {
	Value            string
	VerifyUserTicket string
}

type RecommendInfo struct {
	Ticket     string
	UserName   string
	NickName   string
	Content    string
	Alias      string
	AttrStatus uint32
	City       string
	OpCode     int
	Province   string
	QQNum      int
	Scene      int
	Sex        int
	Signature  string
	VerifyFlag int
}

// ReceivedMessage: for received message
type ReceivedMessage struct {
	IsGroup       bool
	MsgId         string
	Content       string
	FromUserName  string
	ToUserName    string
	Who           string
	MsgType       int
	SubType       int
	OriginContent string
	At            string
	Url           string

	RecommendInfo *RecommendInfo
}
