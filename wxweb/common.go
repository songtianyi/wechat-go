package wxweb

import (
	"encoding/xml"
	"strconv"
	"strings"
)

type Common struct {
	AppId string
	LoginUrl string
	Lang string
	DeviceID string
	UserAgent string
	CgiUrl string
	SyncSrvs []string
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
		strs = append(strs, strconv.Itoa(v.Key) + "_" + strconv.Itoa(v.Val))
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
