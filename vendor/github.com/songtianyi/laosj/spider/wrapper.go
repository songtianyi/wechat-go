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

package spider

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"net/http"
	"sync"
)

// Spider
type Spider struct {
	Url string // page that spider would deal with
	doc *goquery.Document
}

// Start spider
func CreateSpiderFromUrl(url string) (*Spider, error) {
	doc, err := goquery.NewDocument(url)
	if err != nil {
		return nil, fmt.Errorf("url %s, error %s", url, err)
	}
	return &Spider{Url: url, doc: doc}, nil
}

func CreateSpiderFromResponse(r *http.Response) (*Spider, error) {
	doc, err := goquery.NewDocumentFromResponse(r)
	if err != nil {
		return nil, fmt.Errorf("error %s", err)
	}
	return &Spider{doc: doc}, nil
}

func (s *Spider) GetHtml(rule string) ([]string, error) {
	var (
		res = make([]string, 0) //for leaf
		wg  sync.WaitGroup
		mu  sync.Mutex
	)

	s.doc.Find(rule).Each(func(ix int, sl *goquery.Selection) {
		wg.Add(1)
		go func() {
			defer wg.Done()
			content, _ := sl.Html()
			mu.Lock()
			res = append(res, content)
			mu.Unlock()

		}()
	})
	wg.Wait()
	return res, nil
}

func (s *Spider) GetText(rule string) ([]string, error) {
	var (
		res = make([]string, 0) //for leaf
		wg  sync.WaitGroup
		mu  sync.Mutex
	)

	s.doc.Find(rule).Each(func(ix int, sl *goquery.Selection) {
		wg.Add(1)
		go func() {
			defer wg.Done()
			mu.Lock()
			res = append(res, sl.Text())
			mu.Unlock()
		}()
	})
	wg.Wait()
	return res, nil
}

func (s *Spider) GetAttr(rule, attr string) ([]string, error) {
	var (
		res = make([]string, 0) //for leaf
		wg  sync.WaitGroup
		mu  sync.Mutex
	)

	s.doc.Find(rule).Each(func(ix int, sl *goquery.Selection) {
		wg.Add(1)
		go func() {
			defer wg.Done()
			attr, ok := sl.Attr(attr)
			if ok {
				mu.Lock()
				res = append(res, attr)
				mu.Unlock()
			}
		}()
	})
	wg.Wait()
	return res, nil
}
