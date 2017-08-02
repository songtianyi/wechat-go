package joker

import (
	"github.com/songtianyi/rrframework/config"
	"github.com/songtianyi/rrframework/logs"
	"github.com/songtianyi/wechat-go/wxweb"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

var (
	jokeQueue     = make(chan string, 100)
	lastFetchTime = time.Now().Unix() - 2*24*3600
)

// Register plugin
func Register(session *wxweb.Session) {
	session.HandlerRegister.Add(wxweb.MSG_TEXT, wxweb.Handler(listenCmd), "joker")
	if err := session.HandlerRegister.EnableByName("joker"); err != nil {
		logs.Error(err)
		return
	}
	go fetchJokes()

}

func fetchJokes() {
	for {
		km := url.Values{}
		km.Add("sort", "desc")
		km.Add("page", "1")
		km.Add("pagesize", "20")
		km.Add("time", strconv.FormatInt(lastFetchTime, 10))
		km.Add("key", "6f2e982b7b9f86591c063d2db0fb20eb")
		uri := "http://japi.juhe.cn/joke/content/list.from?" + km.Encode()
		resp, err := http.Get(uri)
		if err != nil {
			logs.Error(err)
			continue
		}
		body, _ := ioutil.ReadAll(resp.Body)
		jc, err := rrconfig.LoadJsonConfigFromBytes(body)
		if err != nil {
			logs.Error(err)
			continue
		}
		jokes, _ := jc.GetSliceString("result.data.content")
		times, _ := jc.GetSliceInt64("result.data.unixtime")
		for i, v := range jokes {
			jokeQueue <- v
			if times[i] > lastFetchTime {
				lastFetchTime = times[i]
			}
		}
	}
}

func listenCmd(session *wxweb.Session, msg *wxweb.ReceivedMessage) {
	if contact := session.Cm.GetContactByUserName(msg.FromUserName); contact == nil {
		logs.Error("ignore the messages from", msg.FromUserName)
		return
	}
	if !strings.Contains(msg.Content, "笑话") {
		return
	}
	logs.Debug("jokeQueue length:", len(jokeQueue))
	select {
	case <-time.After(time.Second * 5):
		return
	case joke := <-jokeQueue:
		session.SendText(joke, session.Bot.UserName, wxweb.RealTargetUserName(session, msg))
	}
}
