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
	"github.com/songtianyi/laosj/downloader"
	"github.com/songtianyi/laosj/spider"
	"github.com/songtianyi/rrframework/connector/redis"
	"github.com/songtianyi/rrframework/logs"
	"github.com/songtianyi/rrframework/storage"
	"net/http"
	"time"
	"flag"
)

const (
	UserAgent = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_11_4) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/56.0.2924.87 Safari/537.36"
)

var (
	startPage = flag.String("s", "http://www.douban.com/group/haixiuzu/discussion", "douban group start page")
)

func main() {
	flag.Parse()
	url := *startPage
	d := &downloader.Downloader{
		ConcurrencyLimit: 3,
		UrlChannelFactor: 10,
		RedisConnStr:     "127.0.0.1:6379",
		SourceQueue:      "DATA:IMAGE:HAIXIUZU",
		Store:            rrstorage.CreateLocalDiskStorage("./sexx/haixiuzu/"),
	}
	err, rc := rrredis.GetRedisClient("127.0.0.1:6379")
	if err != nil {
		panic(err)
	}
	go func() {
		d.Start()
	}()
	refer := ""

	for {
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			logs.Error(err)
			break
		}
		req.Header.Add("User-Agent", UserAgent)
		req.Header.Add("Referer", refer)

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			logs.Error(err)
			break
		}
		if resp.StatusCode != 200 {
			logs.Debug(resp)
			break
		}
		s, err := spider.CreateSpiderFromResponse(resp)
		if err != nil {
			logs.Debug(err)
			break
		}
		rs, _ := s.GetAttr("div.grid-16-8.clearfix>div.article>div>table.olt>tbody>tr>td.title>a", "href")
		refer = url
		for _, v := range rs {
			req01, _ := http.NewRequest("GET", v, nil)
			req01.Header.Add("User-Agent", UserAgent)
			req01.Header.Add("Referer", refer)
			resp01, err := client.Do(req01)
			if err != nil {
				logs.Error(err)
				continue
			}
			if resp01.StatusCode != 200 {
				logs.Debug(resp01)
				continue
			}
			s01, err := spider.CreateSpiderFromResponse(resp01)
			if err != nil {
				logs.Error(err)
				continue
			}
			rs01, _ := s01.GetAttr("div.grid-16-8.clearfix>div.article>div.topic-content.clearfix>div.topic-doc>div#link-report>div.topic-content>div.topic-figure.cc>img", "src")
			for _, vv := range rs01 {
				if _, err := rc.RPush("DATA:IMAGE:HAIXIUZU", vv); err != nil {
					logs.Error(err)
				}
			}
			time.Sleep(5 * time.Second)
		}
		rs1, _ := s.GetAttr("div.grid-16-8.clearfix>div.article>div.paginator>span.next>a", "href")
		if len(rs1) != 1 {
			break
		}
		url = rs1[0]
		logs.Notice("redirect to", url)
		time.Sleep(5 * time.Second)
	}
}
