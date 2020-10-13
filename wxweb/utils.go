/*
Copyright 2017 wechat-go Authors. All Rights Reserved.
MIT License

Copyright (c) 2017

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
*/

package wxweb

import (
	"errors"
	"math/rand"
	"reflect"
	"regexp"
	"time"

	"github.com/songtianyi/rrframework/config"
)

func GetRandomStringFromNum(length int) string {
	bytes := []byte("0123456789")
	result := []byte{}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < length; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	return string(result)
}

func GetSyncKeyListFromJc(jc *rrconfig.JsonConfig) (*SyncKeyList, error) {
	is, err := jc.GetInterfaceSlice("SyncKey.List") //[]interface{}
	if err != nil {
		return nil, err
	}
	synks := make([]SyncKey, 0)
	for _, v := range is {
		// interface{}
		vm := v.(map[string]interface{})
		sk := SyncKey{
			Key: int(vm["Key"].(float64)),
			Val: int(vm["Val"].(float64)),
		}
		synks = append(synks, sk)
	}
	return &SyncKeyList{
		Count: len(synks),
		List:  synks,
	}, nil
}

func GetUserInfoFromJc(jc *rrconfig.JsonConfig) (*User, error) {
	user, _ := jc.GetInterface("User")
	u := &User{}
	fields := reflect.ValueOf(u).Elem()
	for k, v := range user.(map[string]interface{}) {
		field := fields.FieldByName(k)
		if vv, ok := v.(float64); ok {
			field.Set(reflect.ValueOf(int(vv)))
		} else {
			field.Set(reflect.ValueOf(v))
		}
	}
	return u, nil
}

func RealTargetUserName(session *Session, msg *ReceivedMessage) string {
	if session.Bot.UserName == msg.FromUserName {
		return msg.ToUserName
	}
	return msg.FromUserName
}

func GetLoginAvatar(resp string) (string, error) {
	match := regexp.MustCompile(`window.userAvatar = '(.+)'`).
		FindStringSubmatch(resp)
	if len(match) < 2 {
		return "", errors.New("login avatar not found")
	}
	return match[1], nil
}
