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
	"encoding/json"
	"strings"
)

type ContactManager struct {
	cl []*User //contact list
}

func CreateContactManagerFromBytes(cb []byte) (*ContactManager, error) {
	var cr ContactResponse
	if err := json.Unmarshal(cb, &cr); err != nil {
		return nil, err
	}
	cm := &ContactManager{
		cl: cr.MemberList,
	}
	return cm, nil
}

func (s *ContactManager) AddConactFromBytes(cb []byte) error {
	var cr ContactResponse
	if err := json.Unmarshal(cb, &cr); err != nil {
		return err
	}
	s.cl = append(s.cl, cr.MemberList...)
	return nil
}

func (s *ContactManager) GetContactByUserName(un string) *User {
	for _, v := range s.cl {
		if v.UserName == un {
			return v
		}
	}
	return nil
}

func (s *ContactManager) GetGroupContact() []*User {
	clarray := make([]*User, 0)
	for _, v := range s.cl {
		if strings.Contains(v.UserName, "@@") {
			clarray = append(clarray, v)
		}
	}
	return clarray
}

func (s *ContactManager) GetContactByName(sig string) []*User {
	clarray := make([]*User, 0)
	for _, v := range s.cl {
		if v.NickName == sig || v.RemarkName == sig {
			clarray = append(clarray, v)
		}
	}
	return clarray
}

func (s *ContactManager) GetContactByQuanPin(sig string) *User {
	for _, v := range s.cl {
		if v.PYQuanPin == sig || v.RemarkPYQuanPin == sig {
			return v
		}
	}
	return nil
}

func (s *ContactManager) GetAll() []*User {
	return s.cl
}
