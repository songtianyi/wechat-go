package wxbot
import (
	"github.com/songtianyi/wechat-go/wxweb"
	"github.com/songtianyi/rrframework/logs"
	"time"
)

func Run() {
	// syncheck
	msg := make(chan []byte, 10000)
	for {
		select {
		case <-time.After(300 * time.Millisecond):
			go producer(msg)
		case m := <-msg:
			go consumer(m)
		}
	}

}

func producer(msg chan []byte) {
	for _, v := range WxWebDefaultCommon.SyncSrvs {
		ret, _, err := wxweb.SyncCheck(WxWebDefaultCommon, WxWebXcg, v, SynKeyList, Cookies)
		if err != nil {
			logs.Error(err)
		}
		if ret == 0 {
			// check success
			err := wxweb.WebWxSync(WxWebDefaultCommon, WxWebXcg, Cookies, msg, SynKeyList)
			if err != nil {
				logs.Error(err)
			}
			break
		}
	}

}

func consumer(msg []byte) {
}
