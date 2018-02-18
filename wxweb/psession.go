package wxweb

import (
	"encoding/json"
	"github.com/songtianyi/rrframework/logs"
	"io/ioutil"
	"net/http"
	"os"
)

type PSession struct {
	WxName      string
	WxWebCommon Common        `json:"common"`
	WxWebXcg    XmlConfig     `json:"config"`
	Cookies     []http.Cookie `json:"cookies,omitempty"`
	Bot         User          `json:"bot,omitempty"`
	QrcodePath  string        `json:"qrcode,omitempty"`
	QrcodeUUID  string        `json:"uuid,omitempty"`
	CreateTime  int64         `json:"time,omitempty"`
}

func WriteSessionData(multiSession map[string]*Session, path string) {
	dataStruct := make(map[string]PSession)

	for k, v := range multiSession {

		pSession := PSession{
			QrcodePath: v.QrcodePath, //qrcode path
			QrcodeUUID: v.QrcodeUUID, //uuid
			CreateTime: v.CreateTime,
		}
		if v.WxWebCommon != nil {
			pSession.WxWebCommon = *v.WxWebCommon
		}
		if v.WxWebXcg != nil {
			pSession.WxWebXcg = *v.WxWebXcg
		}

		if v.Cookies != nil {
			pSession.Cookies = make([]http.Cookie, 0)

			for _, vv := range v.Cookies {
				if vv.Name != "" {
					pSession.Cookies = append(pSession.Cookies, *vv)
				}
			}
		}
		if v.Bot != nil {
			pSession.WxName = v.Bot.NickName
			pSession.Bot = *v.Bot
		}

		dataStruct[k] = pSession
	}
	data, err := json.Marshal(dataStruct)
	if err != nil {
		logs.Error(err)
	} else {
		ioutil.WriteFile(path, data, 0666)
	}
}
func ReadSessionData(path string) map[string]*Session {
	_, err := os.Stat(path)
	dataStruct := make(map[string]PSession)
	multiSession := make(map[string]*Session)
	if err == nil {
		data, error := ioutil.ReadFile(path)
		if error == nil && len(data) > 0 {
			error = json.Unmarshal(data, &dataStruct)
			if error != nil {
				logs.Error(error)
			} else {
				for uuid, session := range dataStruct {
					var cookies = make([]*http.Cookie, 0, len(session.Cookies))

					for k, vv := range session.Cookies {
						if vv.Name != "" {
							cookies = append(cookies, &session.Cookies[k])
						}
					}
					var xml XmlConfig
					xml = session.WxWebXcg
					common := session.WxWebCommon
					wechatSession := &Session{
						WxWebCommon:     &common,
						WxWebXcg:        &xml,
						Cookies:         cookies,
						Bot:             &session.Bot,
						HandlerRegister: CreateHandlerRegister(),
						QrcodePath:      session.QrcodePath, //qrcode path
						QrcodeUUID:      session.QrcodeUUID, //uuid
						CreateTime:      session.CreateTime,
					}
					multiSession[uuid] = wechatSession
				}
			}
		}
	}
	return multiSession
}
