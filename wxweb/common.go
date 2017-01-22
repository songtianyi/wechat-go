package wxweb

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
	Msg          *TextMessage
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
	for ix, v := range s.List {
		strs = append(strs, strconv.FormatInt(v.Key, 10) + "_" + strconv.FormatInt(v.Val, 10)
	}
	return strings.Join(strs, "|")
}

type User struct {
}
