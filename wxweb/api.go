package wxweb

import (
	"net/http"
	"strconv"
	"time"
	"net/url"
	"io/ioutil"
	"strings"
	"fmt"
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


func Login(common *Common, uuid string, tip int) (string, error){
	km := url.Values{}
	km.Add("tip", strconv.FormatInt(tip, 10))
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
	if strings.Contains(strb, "windows.code=200") &&
		strings.Contains(strb, "windows.redirect_uri") {
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
		return fmt.Error("xc.Ret != 0, %s", xc)
	}
	ce = resp.Cookies()
	return nil
}

func WebWxInit(ce *XmlConfig) ([]byte, error){
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
	return body
}

func SyncCheck(ce *XmlConfig, server, synckey string, cookies []*http.Cookie) (int, int, error) {
	km := url.Values{}
	km.Add("r", strconv.FormatInt(time.Now().Unix()*1000, 10))
	km.Add("sid", ce.Wxsid)
	km.Add("uin", ce.Wxuin)
	km.Add("skey", ce.Skey)
	km.Add("deviceid", common.DeviceID)
	km.Add("synckey", syncKey)
	km.Add("_", strconv.FormatInt(time.Now().Unix()*1000, 10))
	uri := "https://" + server + "/cgi-bin/mmwebwx-bin/synccheck?" + km.Encode()

	js := InitReqBody{
		BaseRequest: &BaseRequest{
			ce.Wxuin,
			ce.Wxsid,
			ce.Skey,
			DeviceID,
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

func WebWxSync(ce *XmlConfig, msg chan []byte, skl *SyncKeyList) ([]byte){
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
			DeviceID,
		},
		SyncKey: skl,
		rr: ^int(time.Now().Unix()) + 1,
	}

	b, err := json.Marshal(js)
	if err != nil {
		panic(err)
	}
	jar, _ := cookiejar.New(nil)
	u, _ := url.Parse(uri)
	jar.SetCookies(u, Cookies)
	client := &http.Client{Jar: jar}
	req, err := http.NewRequest("POST", uri, bytes.NewReader(b))
	req.Header.Add("Content-Type", "application/json; charset=UTF-8")
	req.Header.Add("User-Agent", common.UserAgent)

	resp, err := client.Do(req)
	if err != nil {
		return nil
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)

	jc, err := rrconfig.LoadJsonConfigFromBytes(body)
	if err != nil {
		panic(err)
	}
	retcode, err := jc.GetInt("BaseResponse.Ret")
	if err != nil {
		panic(err)
	}
	if retcode != 0 {
		logs.Error("sync message fail")
		return syncKey
	}
	messanger <- body

	synckeylist = synckeylist[:0]
	is, err := jc.GetInterfaceSlice("SyncKey.List") //[]interface{}
	if err != nil {
		panic(err)
	}
	iss := make([]string, 0)
	for _, v := range is {
		// interface{}
		vm := v.(map[string]interface{})
		iss = append(iss, strconv.FormatFloat(vm["Key"].(float64), 'f', -1, 64)+"_"+strconv.FormatFloat(vm["Val"].(float64), 'f', -1, 64))
		sk := SyncKey{
			Key: int(vm["Key"].(float64)),
			Val: int(vm["Val"].(float64)),
		}
		synckeylist = append(synckeylist, sk)
	}
	syncKey = strings.Join(iss, "|")
	return syncKey
}

func WebWxStatusNotify(ce Cookie, userInfo *User) string {
	km := url.Values{}
	km.Add("pass_ticket", ce.PassTicket)
	km.Add("lang", "zh_CN")
	uri := "https://wx.qq.com/cgi-bin/mmwebwx-bin/webwxstatusnotify?" + km.Encode()

	js := InitReqBody{
		BaseRequest: &BaseRequest{
			ce.Wxuin,
			ce.Wxsid,
			ce.Skey,
			DeviceID,
		},
		Code:         3,
		FromUserName: userInfo.UserName,
		ToUserName:   userInfo.UserName,
		ClientMsgId:  int(time.Now().Unix()),
	}

	b, _ := json.Marshal(js)
	client := &http.Client{}
	req, err := http.NewRequest("POST", uri, bytes.NewReader(b))
	req.Header.Add("Content-Type", "application/json; charset=UTF-8")
	req.Header.Add("User-Agent", common.UserAgent)

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	return string(body)
}

func WebWxGetContact(common *Common, ce Cookie) ([]byte, error) {
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
			DeviceID,
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
	jar.SetCookies(u, Cookies)
	client := &http.Client{Jar: jar}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	return body, nil
}

func WebWxSendMsg(common *Common, ce *XmlConfig, 
	from string, to string, msg string, msgId int) {

	km := url.Values{}
	km.Add("pass_ticket", ce.PassTicket)

	uri := common.CgiUrl + "/webwxsendmsg?" + km.Encode()

	js := InitReqBody{
		BaseRequest: &BaseRequest{
			ce.Wxuin,
			ce.Wxsid,
			ce.Skey,
			DeviceID,
		},
		Msg: &TextMessage{
			Type:         1,
			Content:      msg,
			FromUserName: from,
			ToUserName:   to,
			LocalID:      int(time.Now().Unix() * 1e4),
			ClientMsgId:  msgId,
		},
	}

	b, _ := json.Marshal(js)
	req, err := http.NewRequest("POST", uri, bytes.NewReader(b))
	req.Header.Add("Content-Type", "application/json; charset=UTF-8")
	req.Header.Add("User-Agent", common.UserAgent)

	jar, _ := cookiejar.New(nil)
	u, _ := url.Parse(uri)
	jar.SetCookies(u, Cookies)
	client := &http.Client{Jar: jar}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	logs.Debug("sending response", string(body))
}
