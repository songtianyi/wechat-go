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

package faceplusplus

import (
	"bytes"
	"github.com/songtianyi/rrframework/config"
	"github.com/songtianyi/rrframework/logs"
	"github.com/songtianyi/wechat-go/wxweb"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"strconv"
	"strings"
)

const (
	apiKey    = "7qKNKrhL3wPTjaL5frh4CjYgZ0DjtH1q"
	apiSecret = "DJmEVbYsEX-vrgyn_xAKJ9yxTxelFsBV"
	apiUrl    = "https://api-cn.faceplusplus.com/facepp/v3/"
)

// Register 注册函数
func Register(session *wxweb.Session) {
	session.HandlerRegister.Add(wxweb.MSG_IMG, wxweb.Handler(faceDetectHandle), "faceplusplus")
	if err := session.HandlerRegister.EnableByName("faceplusplus"); err != nil {
		logs.Error(err)
	}
}

func faceDetectHandle(session *wxweb.Session, msg *wxweb.ReceivedMessage) {
	contact := session.Cm.GetContactByUserName(msg.FromUserName)
	// contact filter
	if contact == nil {
		logs.Error("no this contact, ignore", msg.FromUserName)
		return
	}

	b, err := session.GetImg(msg.MsgId)
	if err != nil {
		logs.Error(err)
		return
	}
	res, err := detect(msg.MsgId+".jpg", b)
	if err != nil {
		logs.Error(err)
		return
	}

	jc, _ := rrconfig.LoadJsonConfigFromBytes(res)
	//du, _ := jc.Dump()
	//logs.Debug(du)
	ages, err := jc.GetSliceInt("faces.attributes.age.value")
	if err != nil {
		logs.Error(err)
		return
	}
	genders, _ := jc.GetSliceString("faces.attributes.gender.value")
	str := ""
	for i, v := range ages {
		if strings.ToLower(genders[i]) == "female" {
			v = v * 8 / 10
		}
		str += genders[i] + "," + strconv.Itoa(v) + "\n"
	}
	if session.Bot.UserName == msg.FromUserName {
		session.SendText(str, session.Bot.UserName, msg.ToUserName)
	} else {
		session.SendText(str, session.Bot.UserName, msg.FromUserName)
	}

}

func detect(filename string, content []byte) ([]byte, error) {
	var b bytes.Buffer
	w := multipart.NewWriter(&b)
	fw, _ := w.CreateFormFile("image_file", filename)
	if _, err := io.Copy(fw, bytes.NewReader(content)); err != nil {
		return nil, err
	}
	fw, _ = w.CreateFormField("api_key")
	_, _ = fw.Write([]byte(apiKey))
	fw, _ = w.CreateFormField("api_secret")
	_, _ = fw.Write([]byte(apiSecret))
	fw, _ = w.CreateFormField("return_landmark")
	_, _ = fw.Write([]byte("0"))
	fw, _ = w.CreateFormField("return_attributes")
	_, _ = fw.Write([]byte("gender,age"))

	w.Close()
	client := &http.Client{}
	req, err := http.NewRequest("POST", apiUrl+"detect", &b)
	req.Header.Add("Content-Type", w.FormDataContentType())
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	body, _ := ioutil.ReadAll(resp.Body)
	return body, nil
}
