package wxweb

import (
	"net/http"
	"strconv"
	"time"
	"net/url"
	"io/ioutil"
	"strings"
	"fmt"
	"encoding/xml"
	"encoding/json"
	"bytes"
	"net/http/cookiejar"
	"regexp"
	"github.com/songtianyi/rrframework/config"
	"github.com/songtianyi/rrframework/logs"
)


func JsLogin(common *Common) (string, error) {
	km := url.Values{}
	km.Add("appid", common.AppId)
	km.Add("fun", "new")
	km.Add("lang", common.Lang)
	km.Add("_", strconv.FormatInt(time.Now().Unix(), 10))
	uri := common.LoginUrl + "/jslogin?" + km.Encode()
	resp, err := http.Get(uri)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	ss := strings.Split(string(body), "\"")
	if len(ss) < 2 {
		return "", fmt.Errorf("jslogin response invalid, %s", string(body))
	}
	return ss[1], nil
}

func QrCode(common *Common, uuid string) ([]byte, error) {
	km := url.Values{}
	km.Add("t", "webwx")
	km.Add("_", strconv.FormatInt(time.Now().Unix(), 10))
	uri := common.LoginUrl + "/qrcode/" + uuid + "?" + km.Encode()
	resp, err := http.Post(uri, "application/octet-stream", nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	return body, nil
}


func Login(common *Common, uuid, tip string) (string, error){
	km := url.Values{}
	km.Add("tip", tip)
	km.Add("uuid", uuid)
	km.Add("_", strconv.FormatInt(time.Now().Unix(), 10))
	uri := common.LoginUrl + "/cgi-bin/mmwebwx-bin/login?" + km.Encode()
	resp, err := http.Get(uri)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	strb := string(body)
	logs.Debug(strb)
	if strings.Contains(strb, "window.code=200") &&
		strings.Contains(strb, "window.redirect_uri") {
		ss := strings.Split(strb, "\"")
		if len(ss) < 2 {
			return "", fmt.Errorf("parse redirect_uri fail, %s", strb)
		}
		return ss[1], nil
	}else {
		return "", fmt.Errorf("invalid response, %s", strb)
	}
}

func WebNewLoginPage(xc *XmlConfig, ce []*http.Cookie, uri string) error {
	resp, err := http.Get(uri)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	if err := xml.Unmarshal(body, xc); err != nil {
		return err
	}
	if xc.Ret != 0 {
		return fmt.Errorf("xc.Ret != 0, %s", xc)
	}
	ce = resp.Cookies()
	return nil
}

func WebWxInit(common *Common, ce *XmlConfig) ([]byte, error){
	km := url.Values{}
	km.Add("pass_ticket", ce.PassTicket)
	km.Add("skey", ce.Skey)
	km.Add("r", strconv.FormatInt(time.Now().Unix(), 10))

	uri := common.CgiUrl + "/webwxinit?" + km.Encode()

	js := InitReqBody{
		BaseRequest: &BaseRequest{
			ce.Wxuin,
			ce.Wxsid,
			ce.Skey,
			common.DeviceID,
		},
	}

	b, _ := json.Marshal(js)
	client := &http.Client{}
	req, err := http.NewRequest("POST", uri, bytes.NewReader(b))
	req.Header.Add("Content-Type", "application/json; charset=UTF-8")
	req.Header.Add("User-Agent", common.UserAgent)

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	return body, nil
}

func SyncCheck(common *Common, ce *XmlConfig,
	server string, skl *SyncKeyList, cookies []*http.Cookie) (int, int, error) {
	km := url.Values{}
	km.Add("r", strconv.FormatInt(time.Now().Unix()*1000, 10))
	km.Add("sid", ce.Wxsid)
	km.Add("uin", ce.Wxuin)
	km.Add("skey", ce.Skey)
	km.Add("deviceid", common.DeviceID)
	km.Add("synckey", skl.String())
	km.Add("_", strconv.FormatInt(time.Now().Unix()*1000, 10))
	uri := "https://" + server + "/cgi-bin/mmwebwx-bin/synccheck?" + km.Encode()

	js := InitReqBody{
		BaseRequest: &BaseRequest{
			ce.Wxuin,
			ce.Wxsid,
			ce.Skey,
			common.DeviceID,
		},
	}

	b, _ := json.Marshal(js)
	jar, _ := cookiejar.New(nil)
	u, _ := url.Parse(uri)
	jar.SetCookies(u, cookies)
	client := &http.Client{Jar: jar}
	req, err := http.NewRequest("GET", uri, bytes.NewReader(b))
	if err != nil {
		return 0, 0, err
	}

	req.Header.Add("Content-Type", "application/json; charset=UTF-8")
	req.Header.Add("User-Agent", common.UserAgent)

	resp, err := client.Do(req)
	if err != nil {
		return 0, 0, err
	}

	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	strb := string(body)
	reg := regexp.MustCompile("window.synccheck={retcode:\"(\\d+)\",selector:\"(\\d+)\"}")
	sub := reg.FindStringSubmatch(strb)
	retcode, _ := strconv.Atoi(sub[1])
	selector, _ := strconv.Atoi(sub[2])
	return retcode, selector, nil
}

func WebWxSync(common *Common,
	ce *XmlConfig,
	cookies []*http.Cookie,
	msg chan []byte, skl *SyncKeyList) error {

	km := url.Values{}
	km.Add("skey", ce.Skey)
	km.Add("sid", ce.Wxsid)
	km.Add("lang", common.Lang)
	km.Add("pass_ticket", ce.PassTicket)

	uri := common.CgiUrl + "/webwxsync?" + km.Encode()

	js := InitReqBody{
		BaseRequest: &BaseRequest{
			ce.Wxuin,
			ce.Wxsid,
			ce.Skey,
			common.DeviceID,
		},
		SyncKey: skl,
		rr: ^int(time.Now().Unix()) + 1,
	}

	b, _ := json.Marshal(js)
	jar, _ := cookiejar.New(nil)
	u, _ := url.Parse(uri)
	jar.SetCookies(u, cookies)
	client := &http.Client{Jar: jar}
	req, err := http.NewRequest("POST", uri, bytes.NewReader(b))
	req.Header.Add("Content-Type", "application/json; charset=UTF-8")
	req.Header.Add("User-Agent", common.UserAgent)

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)

	jc, err := rrconfig.LoadJsonConfigFromBytes(body)
	if err != nil {
		return err
	}
	retcode, _ := jc.GetInt("BaseResponse.Ret")
	if retcode != 0 {
		return fmt.Errorf("BaseResponse.Ret %d", retcode)
	}

	msg <- body

	skl.List = skl.List[:0]
	skl1, _ := GetSyncKeyListFromJc(jc)
	skl.Count = skl1.Count
	skl.List = append(skl.List, skl1.List...)
	return nil
}

func WebWxStatusNotify(common *Common, ce *XmlConfig, bot *User) (int, error) {
	km := url.Values{}
	km.Add("pass_ticket", ce.PassTicket)
	km.Add("lang", common.Lang)
	uri := common.CgiUrl + "/webwxstatusnotify?" + km.Encode()

	js := InitReqBody{
		BaseRequest: &BaseRequest{
			ce.Wxuin,
			ce.Wxsid,
			ce.Skey,
			common.DeviceID,
		},
		Code: 3,
		FromUserName: bot.UserName,
		ToUserName:   bot.UserName,
		ClientMsgId:  int(time.Now().Unix()),
	}

	b, _ := json.Marshal(js)
	client := &http.Client{}
	req, err := http.NewRequest("POST", uri, bytes.NewReader(b))
	req.Header.Add("Content-Type", "application/json; charset=UTF-8")
	req.Header.Add("User-Agent", common.UserAgent)

	resp, err := client.Do(req)
	if err != nil {
		return -1, err
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	jc, _ := rrconfig.LoadJsonConfigFromBytes(body)
	ret, _ := jc.GetInt("BaseResponse.Ret")
	return ret, nil
}

func WebWxGetContact(common *Common, ce *XmlConfig, cookies []*http.Cookie) ([]byte, error) {
	km := url.Values{}
	km.Add("r", strconv.FormatInt(time.Now().Unix(), 10))
	km.Add("seq", "0")
	km.Add("skey", ce.Skey)
	uri := common.CgiUrl + "/webwxgetcontact?" + km.Encode()

	js := InitReqBody{
		BaseRequest: &BaseRequest{
			ce.Wxuin,
			ce.Wxsid,
			ce.Skey,
			common.DeviceID,
		},
	}

	b, _ := json.Marshal(js)
	req, err := http.NewRequest("POST", uri, bytes.NewReader(b))
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/json; charset=UTF-8")
	req.Header.Add("User-Agent", common.UserAgent)

	jar, _ := cookiejar.New(nil)
	u, _ := url.Parse(uri)
	jar.SetCookies(u, cookies)
	client := &http.Client{Jar: jar}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	return body, nil
}

func WebWxSendMsg(common *Common, ce *XmlConfig, cookies []*http.Cookie,
	from, to string, msg interface{}) (int, error) {

	km := url.Values{}
	km.Add("pass_ticket", ce.PassTicket)

	uri := common.CgiUrl + "/webwxsendmsg?" + km.Encode()

	js := InitReqBody{
		BaseRequest: &BaseRequest{
			ce.Wxuin,
			ce.Wxsid,
			ce.Skey,
			common.DeviceID,
		},
		Msg: msg,
	}

	b, _ := json.Marshal(js)
	req, err := http.NewRequest("POST", uri, bytes.NewReader(b))
	if err != nil {
		return -1, err
	}
	req.Header.Add("Content-Type", "application/json; charset=UTF-8")
	req.Header.Add("User-Agent", common.UserAgent)

	jar, _ := cookiejar.New(nil)
	u, _ := url.Parse(uri)
	jar.SetCookies(u, cookies)
	client := &http.Client{Jar: jar}
	resp, err := client.Do(req)
	if err != nil {
		return -1, err
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	jc, _ := rrconfig.LoadJsonConfigFromBytes(body)
	ret, _ := jc.GetInt("BaseResponse.Ret")
	return ret, nil
}
