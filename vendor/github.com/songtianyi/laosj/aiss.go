// Copyright 2016 laosj Author @songtianyi. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"fmt"
	"github.com/songtianyi/laosj/downloader"
	"github.com/songtianyi/rrframework/config"
	"github.com/songtianyi/rrframework/connector/redis"
	"github.com/songtianyi/rrframework/storage"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
)

var (
	kmlock sync.Mutex
)

type ReqBody struct {
	page   int
	userId int
}

func getSuiteList(page int) ([]byte, error) {
	uri := "http://api.pmkoo.cn/aiss/suite/suiteList.do"
	para := "page=" + strconv.FormatInt(int64(page), 10) + "&userId=153044"
	client := &http.Client{}
	req, err := http.NewRequest("POST", uri, strings.NewReader(para))
	req.Header.Add("Host", "api.pmkoo.cn")
	req.Header.Add("Accept", "*/*")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Accept-Language", "zh-Haq=1, en-CN;q=0.9")
	req.Header.Add("User-Agent", "aiss/1.0 (iPhone; iOS 10.2; Scale/2.00)")
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	return body, nil
}

func main() {
	//
	oss := "http://com-pmkoo-img.oss-cn-beijing.aliyuncs.com/picture/"
	sema := make(chan struct{}, 10)
	page := 1
	ok := true

	err, rc := rrredis.GetRedisClient("127.0.0.1:6379")
	if err != nil {
		panic(err)
	}
	km := make(map[string]bool)
	for ok {
		select {
		case sema <- struct{}{}:
			go func(pg int) {
				defer func() { <-sema }() // release
				b, err := getSuiteList(pg)
				if err != nil {
					fmt.Println(err)
					ok = false
					return
				}
				jc, _ := rrconfig.LoadJsonConfigFromBytes(b)
				ics, err := jc.GetInterfaceSlice("data.list")
				if err != nil {
					fmt.Println(err)
					ok = false
					return
				}
				for _, v := range ics {
					vm := v.(map[string]interface{})
					vsource := vm["source"].(map[string]interface{})
					catlog := vsource["catalog"].(string)
					pictureCount := int(vm["pictureCount"].(float64))
					issue := int(vm["issue"].(float64))
					for j := 0; j < pictureCount; j++ {
						uri := oss + catlog + "/"
						uri += strconv.FormatInt(int64(issue), 10) + "/"
						uri += strconv.FormatInt(int64(j), 10) + ".jpg"
						key := "DATA:IMAGE:" + catlog + ":" + strconv.FormatInt(int64(issue), 10)
						kmlock.Lock()
						km[key] = true
						kmlock.Unlock()
						if _, err := rc.RPush(key, uri); err != nil {
							fmt.Println(err)
						}
					}
				}
			}(page)
			page += 1
		}
	}
	var wg sync.WaitGroup
	for k := range km {
		wg.Add(1)
		go func(k string) {
			defer wg.Done()
			_ = os.MkdirAll("/data/sexx/pmkoo/"+k, os.ModeDir)
			d := &downloader.Downloader{
				ConcurrencyLimit: 10,
				UrlChannelFactor: 10,
				RedisConnStr:     "127.0.0.1:6379",
				SourceQueue:      k,
				Store:            rrstorage.CreateLocalDiskStorage("/data/sexx/pmkoo/" + k),
			}
			go d.Start()
			d.WaitCloser()
		}(k)
	}
	wg.Wait()

}
