package wxbot
import (
	"github.com/songtianyi/wechat-go/wxweb"
)

func Run() {
	// syncheck
	msg := make(chan []byte, 10000)
	for {
		select {
		case <-time.After(300 * time.MilliSecond):
			go producer(msg)
		case m := <-msg:
			go consumer(m)
		}
	}

}

func producer(msg chan []byte) {
	for _, v := range wxweb.DefaultCommon.SyncSrvs {
		if wxweb.SyncCheck(v, XmlCfg, SynKeyList) {
			// check success
			syncKey = wxweb.WxWebSync(XmlCfg, syncKey, msg, SynKeyList)
			break
		}
	}

}

func consumer() {
}
