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
	//"github.com/songtianyi/laosj/downloader"
	"flag"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/songtianyi/laosj/spider"
	"github.com/songtianyi/rrframework/logs"
	"github.com/songtianyi/rrframework/storage"
	"github.com/songtianyi/rrframework/utils"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
	"sync"
)

var (
	cmd = flag.String("cmd", "", "choose what you want\n"+
		"-cmd list, get all av star videos list\n"+
		"-cmd get -art GS-064, get specified video torrent file\n")
	artOp   = flag.String("art", "", "specify one art")
	withImg = flag.Int("image", 0, "-image {1|0} whether show video image in list")
)

const (
	JAV_PREFIX = "http://www.javlibrary.com/cn/"
)

func parseStar(item string, artist string) {
	m := regexp.MustCompile("href=\"(\\S+)\"").FindStringSubmatch(item)
	if len(m) < 2 {
		logs.Error("Find href error, %s", "len(m) < 2")
		return
	}
	url := m[1]
	jvname := strings.Split(url, "=")[1]
	s, err := spider.CreateSpiderFromUrl(JAV_PREFIX + url + "&mode=2")
	if err != nil {
		logs.Error(err)
		return
	}
	arts, _ := s.GetText("div.videothumblist>div.videos>div.video>a>div.id")
	artsrcs, _ := s.GetHtml("div.videothumblist>div.videos>div.video>a")

	if len(artsrcs) != len(arts) {
		logs.Error(s.Url, len(artsrcs), len(arts))
		logs.Debug(artsrcs)
		logs.Debug(arts)
		panic(fmt.Errorf(""))
		return
	}
	for i, artsrc := range artsrcs {
		m1 := regexp.MustCompile("src=\"(\\S+)\"").FindStringSubmatch(artsrc)
		if len(m1) < 2 {
			logs.Error("Find href fail, %s", artsrc)
			return
		}
		coverurl := m1[1]
		if *withImg == 0 {
			logs.Info(jvname, arts[i], artist, m1[1])
		} else {
			surl := strings.Split(coverurl, "/")
			filename := surl[len(surl)-1]
			logs.Info(jvname, arts[i], artist, filename)
		}
	}
}

func runList() {
	// for every prefix=?
	for i := 0; i < 26; i++ {
		// get page num
		s, err := spider.CreateSpiderFromUrl(JAV_PREFIX + "star_list.php?prefix=" + string(rune(i+65)))
		if err != nil {
			logs.Error(err)
			continue
		}
		rs, _ := s.GetText("div.page_selector>a.page")
		max := spider.FindMaxFromSliceString(0, rs)
		// for every page
		var wg sync.WaitGroup
		for j := 1; j <= max; j++ {
			// find stars in this page
			wg.Add(1)
			go func(k int) {
				defer wg.Done()
				s1, err := spider.CreateSpiderFromUrl(s.Url + "&page=" + string(rune(k)))
				if err != nil {
					logs.Error(err)
					return
				}
				rs1, _ := s1.GetHtml("div.searchitem")
				rs2, _ := s1.GetText("div.searchitem>a")
				if len(rs1) != len(rs2) {
					logs.Error("assert fail")
					return
				}
				for n := range rs1 {
					parseStar(rs1[n], rs2[n])
				}
			}(j)
		}
		wg.Wait()
	}
}

func runGet(art string) {
	store := rrstorage.CreateLocalDiskStorage(".")

	// get download url
	url := "https://sukebei.nyaa.se/?page=search&cats=8_30&filter=0&sort=4&term=" + art
	rul := "div.content>table.tlist>tbody>tr.tlistrow.trusted>td.tlistdownload"
	doc, err := goquery.NewDocument(url)
	if err != nil {
		logs.Error(err)
		return
	}
	doc.Find(rul).Each(func(ix int, sl *goquery.Selection) {
		uri, _ := sl.Find("a").Attr("href")
		uri = "https:" + uri

		logs.Info("downloading", uri)
		resp, err := http.Get(uri)
		if err != nil {
			logs.Error(err)
			return
		}
		b, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			logs.Error(err)
			return
		}
		if err := store.Save(b, art+"-"+rrutils.NewV4().String()+".torrent"); err != nil {
			logs.Error(err)
			return
		}
	})
}

func main() {
	flag.Parse()
	if *cmd == "list" {
		runList()
	} else if *cmd == "get" {
		if *artOp == "" {
			logs.Error("wrong usage\n./jav -cmd get -art IENE-706")
			return
		}
		runGet(*artOp)
	} else {
	}
}
