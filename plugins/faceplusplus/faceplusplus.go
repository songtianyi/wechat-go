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
)

const (
	apiKey    = "7qKNKrhL3wPTjaL5frh4CjYgZ0DjtH1q"
	apiSecret = "DJmEVbYsEX-vrgyn_xAKJ9yxTxelFsBV"
	apiUrl    = "https://api-cn.faceplusplus.com/facepp/v3/"
)

func Register(session *wxweb.Session) {
	session.HandlerRegister.Add(3, wxweb.Handler(FaceDetectHandle))

}

func FaceDetectHandle(session *wxweb.Session, msg *wxweb.ReceivedMessage) {
	contact := session.Cm.GetContactByUserName(msg.FromUserName)
	// contact filter
	if contact == nil {
		logs.Error("no this contact", msg.FromUserName)
		return
	}

	b, err := session.GetImg(msg.MsgId)
	if err != nil {
		logs.Error(err)
		return
	}
	res, err := Detect(msg.MsgId+".jpg", b)
	if err != nil {
		logs.Error(err)
		return
	}
	logs.Debug(string(res))
	jc, _ := rrconfig.LoadJsonConfigFromBytes(res)
	ages, err := jc.GetSliceInt("faces.attributes.age.value")
	if err != nil {
		logs.Error(err)
		return
	}
	genders, _ := jc.GetSliceString("faces.attributes.gender.value")
	str := ""
	for i, v := range ages {
		str += genders[i] + "," + strconv.Itoa(v) + "\n"
	}
	if !msg.IsGroup {
		session.SendText(str, session.Bot.UserName, msg.FromUserName)
	}

}

func Detect(filename string, content []byte) ([]byte, error) {
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
