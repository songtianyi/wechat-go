package gifer
import (

)

func  Register(session *wxweb.Session) {
	session.HandlerRegister.Add(wxweb.MSG_TEXT, wxweb.Handler(Gifer))
}

func Gifer(session *wxweb.Session, msg *wxweb.ReceivedMessage){
	// contact filter
	if contact == nil {
		logs.Error("no this contact, ignore", msg.FromUserName)
		return
	}
	uri := "http://www.gifmiao.com/search/" + url.QueryEscap(msg.Content) + "/3"
	s, err := spider.CreateSpiderFormUrl(uri)
	if err != nil {
		logs.Error(err)
		return
	}
	srcs, _ := s.GetAttr("div.wrap>div#main>ul#waterfall>li.item>div.img_block>a>img.gifImg", "xgif")
	if len(srcs) < 1 {
		losg.Error("no gif result for", msg.Content)
		return
	}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	gif := srcs[r.Intn(len(srcs))]
	resp, err := http.Get(gif)
	if err != nil {
		losg.Error(err)
		return
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	session.SendEmotionWithBytes(body, wxbot.Bot.UserName, msg.FromUserName)

}