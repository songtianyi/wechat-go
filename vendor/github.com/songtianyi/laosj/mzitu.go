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
	"regexp"
	"strconv"
	"sync"
)

func main() {
	d := &downloader.Downloader{
		ConcurrencyLimit: 10,
		UrlChannelFactor: 10,
		RedisConnStr:     "127.0.0.1:6379",
		SourceQueue:      "DATA:IMAGE:MZITU:XINGGAN",
		Store:            rrstorage.CreateLocalDiskStorage("/data/sexx/taiwan/"),
	}
	go func() {
		d.Start()
	}()

	// step1: find total index pages
	s, err := spider.CreateSpiderFromUrl("http://www.mzitu.com/taiwan")
	if err != nil {
		logs.Error(err)
		return
	}
	rs, _ := s.GetText("div.main>div.main-content>div.postlist>nav.navigation.pagination>div.nav-links>a.page-numbers")
	max := spider.FindMaxFromSliceString(1, rs)

	// step2: for every index page, find every post entrance
	var wg sync.WaitGroup
	var mu sync.Mutex
	step2 := make([]string, 0)
	for i := 1; i <= max; i++ {
		wg.Add(1)
		go func(ix int) {
			defer wg.Done()
			ns, err := spider.CreateSpiderFromUrl(s.Url + "/page/" + strconv.Itoa(ix))
			if err != nil {
				logs.Error(err)
				return
			}
			t, _ := ns.GetHtml("div.main>div.main-content>div.postlist>ul>li")
			mu.Lock()
			step2 = append(step2, t...)
			mu.Unlock()
		}(i)
	}
	wg.Wait()
	// parse url
	for i, v := range step2 {
		re := regexp.MustCompile("href=\"(\\S+)\"")
		m := re.FindStringSubmatch(v)
		if len(m) < 2 {
			continue
		}
		step2[i] = m[1]
	}

	for _, v := range step2 {
		// step3: step in entrance, find max pagenum
		ns1, err := spider.CreateSpiderFromUrl(v)
		if err != nil {
			logs.Error(err)
			return
		}
		t1, _ := ns1.GetText("div.main>div.content>div.pagenavi>a")
		maxx := spider.FindMaxFromSliceString(1, t1)
		// step4: for every page
		for j := 1; j <= maxx; j++ {

			// step5: find img in this page
			ns2, err := spider.CreateSpiderFromUrl(v + "/" + strconv.Itoa(j))
			if err != nil {
				logs.Error(err)
				return
			}
			t2, err := ns2.GetHtml("div.main>div.content>div.main-image>p>a")
			if len(t2) < 1 {
				// ignore this page
				continue
			}
			sub := regexp.MustCompile("src=\"(\\S+)\"").FindStringSubmatch(t2[0])
			if len(sub) != 2 {
				// ignore this page
				continue
			}
			err, rc := rrredis.GetRedisClient(d.RedisConnStr)
			if err != nil {
				logs.Error(err)
				return
			}
			key := d.SourceQueue
			if _, err := rc.RPush(key, sub[1]); err != nil {
				logs.Error(err)
				return
			}
		}
	}
	d.WaitCloser()
}
