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
		case <-time.After(1 * time.Second):
			go producer(msg)
		case m := <-msg:
			go consumer(m)
		}
	}

}

func producer(msg chan []byte) {
	for _, v := range WxWebDefaultCommon.SyncSrvs {
		ret, sel, err := wxweb.SyncCheck(WxWebDefaultCommon, WxWebXcg, Cookies, v, SynKeyList)
		logs.Debug(v, ret, sel)
		if err != nil {
			logs.Error(err)
		}
		if ret == 0 {
			// check success
			if sel == 2 {
				// new message
				err := wxweb.WebWxSync(WxWebDefaultCommon, WxWebXcg, Cookies, msg, SynKeyList)
				if err != nil {
					logs.Error(err)
				}
			}
			break
		}
	}

}

func consumer(msg []byte) {
	logs.Debug("received", string(msg))
}
