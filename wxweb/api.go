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
	"io"
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"github.com/songtianyi/rrframework/config"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"
	"mime/multipart"
	"sync/atomic"
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

func Login(common *Common, uuid, tip string) (string, error) {
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
	if strings.Contains(strb, "window.code=200") &&
		strings.Contains(strb, "window.redirect_uri") {
		ss := strings.Split(strb, "\"")
		if len(ss) < 2 {
			return "", fmt.Errorf("parse redirect_uri fail, %s", strb)
		}
		return ss[1], nil
	} else {
		return "", fmt.Errorf("invalid response, %s", strb)
	}
}

func WebNewLoginPage(common *Common, xc *XmlConfig, uri string) ([]*http.Cookie, error) {
	parsed, _ := url.Parse(uri)
	km := parsed.Query()
	km.Add("fun", "new")
	uri = common.CgiUrl + "/webwxnewloginpage?" + km.Encode()

	resp, err := http.Get(uri)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	if err := xml.Unmarshal(body, xc); err != nil {
		return nil, err
	}
	if xc.Ret != 0 {
		return nil, fmt.Errorf("xc.Ret != 0, %s", xc)
	}
	return resp.Cookies(), nil
}

func WebWxInit(common *Common, ce *XmlConfig) ([]byte, error) {
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

func SyncCheck(common *Common, ce *XmlConfig, cookies []*http.Cookie,
	server string, skl *SyncKeyList) (int, int, error) {
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
		rr:      ^int(time.Now().Unix()) + 1,
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
		Code:         3,
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

func WebWxSendTextMsg(common *Common, ce *XmlConfig, cookies []*http.Cookie,
	from, to string, msg string) (int, error) {

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
		Msg: &TextMessage{
			Type: 1,
			Content: msg,
			FromUserName: from,
			ToUserName: to,
			LocalID: int(time.Now().Unix() * 1e4),
			ClientMsgId: int(time.Now().Unix() * 1e4),
		},
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

func WebWxUploadMedia(common *Common, ce *XmlConfig, cookies []*http.Cookie,
	filename string, content []byte) (string, error) {

	var b bytes.Buffer
	w := multipart.NewWriter(&b)
	fw, _ := w.CreateFormFile("filename", filename)
	if _, err := io.Copy(fw, bytes.NewReader(content)); err != nil {
		return "", err
	}

	fw, _ = w.CreateFormField("id")
	_, _ = fw.Write([]byte("WU_FILE_" + strconv.Itoa(int(common.MediaCount))))
	common.MediaCount = atomic.AddUint32(&common.MediaCount, 1)

	fw, _ = w.CreateFormField("name")
	_, _ = fw.Write([]byte(filename))

	fw, _ = w.CreateFormField("type")
	_, _ = fw.Write([]byte("image/jpeg"))

	fw, _ = w.CreateFormField("lastModifieDate")
	_, _ = fw.Write([]byte("Mon Feb 13 2017 17:27:23 GMT+0800 (CST)"))

	fw, _ = w.CreateFormField("size")
	_, _ = fw.Write([]byte(strconv.Itoa(len(content))))

	fw, _ = w.CreateFormField("mediatype")
	_, _ = fw.Write([]byte("pic"))


	js := InitReqBody{
		BaseRequest: &BaseRequest{
			ce.Wxuin,
			ce.Wxsid,
			ce.Skey,
			common.DeviceID,
		},
		ClientMediaId: int(time.Now().Unix() * 1e4),
		TotalLen: len(content),
		StartPos: 0,
		DataLen: len(content),
		MediaType: 4,
	}

	jb, _ := json.Marshal(js)

	fw, _ = w.CreateFormField("uploadmediarequest")
	_, _ = fw.Write(jb)

	fw, _ = w.CreateFormField("webwx_data_ticket")
	for _, v := range cookies {
		if strings.Contains(v.String(), "webwx_data_ticket") {
			_, _ = fw.Write([]byte(strings.Split(v.String(), "=")[1]))
			break
		}
	}

	fw, _ = w.CreateFormField("pass_ticket")
	_, _ = fw.Write([]byte(ce.PassTicket))
	w.Close()

	req, err := http.NewRequest("POST", common.UploadUrl, &b)
	if err != nil {
		return "", err
	}
	req.Header.Add("Content-Type",  w.FormDataContentType())
	req.Header.Add("User-Agent", common.UserAgent)

	jar, _ := cookiejar.New(nil)
	u, _ := url.Parse(common.UploadUrl)
	jar.SetCookies(u, cookies)
	client := &http.Client{Jar: jar}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	jc, _ := rrconfig.LoadJsonConfigFromBytes(body)
	ret, _ := jc.GetInt("BaseResponse.Ret")
	if ret != 0 {
		return "", fmt.Errorf("BaseResponse.Ret=%d", ret)
	}
	mediaId, _ := jc.GetString("MediaId")
	return mediaId, nil
}

func WebWxSendMsgImg(common *Common, ce *XmlConfig, cookies []*http.Cookie,
	from, to, media string) (int, error) {

	km := url.Values{}
	km.Add("pass_ticket", ce.PassTicket)
	km.Add("fun", "async")
	km.Add("f", "json")
	km.Add("lang", common.Lang)

	uri := common.CgiUrl + "/webwxsendmsgimg?" + km.Encode()

	js := InitReqBody{
		BaseRequest: &BaseRequest{
			ce.Wxuin,
			ce.Wxsid,
			ce.Skey,
			common.DeviceID,
		},
		Msg: &MediaMessage{
			Type: 3,
			Content: "",
			FromUserName: from,
			ToUserName: to,
			LocalID: int(time.Now().Unix() * 1e4),
			ClientMsgId: int(time.Now().Unix() * 1e4),
			MediaId: media,
		},
		Scene: 0,
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
