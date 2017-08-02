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

var (
	// SpecialContact: special contacts map
	SpecialContact = map[string]bool{
		"filehelper":            true,
		"newsapp":               true,
		"fmessage":              true,
		"weibo":                 true,
		"qqmail":                true,
		"tmessage":              true,
		"qmessage":              true,
		"qqsync":                true,
		"floatbottle":           true,
		"lbsapp":                true,
		"shakeapp":              true,
		"medianote":             true,
		"qqfriend":              true,
		"readerapp":             true,
		"blogapp":               true,
		"facebookapp":           true,
		"masssendapp":           true,
		"meishiapp":             true,
		"feedsapp":              true,
		"voip":                  true,
		"blogappweixin":         true,
		"weixin":                true,
		"brandsessionholder":    true,
		"weixinreminder":        true,
		"officialaccounts":      true,
		"wxitil":                true,
		"userexperience_alarm":  true,
		"notification_messages": true,
	}
)

// ContactManager: contact manager
type ContactManager struct {
	cl []*User //contact list
}

// create
// CreateContactManagerFromBytes: create contact maanger from bytes
func CreateContactManagerFromBytes(cb []byte) (*ContactManager, error) {
	var cr WxWebGetContactResponse
	if err := json.Unmarshal(cb, &cr); err != nil {
		return nil, err
	}
	cm := &ContactManager{
		cl: cr.MemberList,
	}
	return cm, nil
}

// update
// AddContactFromBytes: upate contact manager from wxwebgetcontact response
func (s *ContactManager) AddUserFromBytes(cb []byte) error {
	var cr WxWebGetContactResponse
	if err := json.Unmarshal(cb, &cr); err != nil {
		return err
	}
	s.cl = append(s.cl, cr.MemberList...)
	return nil
}

// AddContactFromUser: add a new user to contact manager
func (s *ContactManager) AddUser(user *User) {
	if user == nil {
		return
	}
	s.cl = append(s.cl, user)
}

// get
// GetContactByUserName: get contact by UserName
func (s *ContactManager) GetContactByUserName(un string) *User {
	for _, v := range s.cl {
		if v.UserName == un {
			return v
		}
	}
	return nil
}

// GetGroupContacts: get all group contacts
func (s *ContactManager) GetGroupContacts() []*User {
	clarray := make([]*User, 0)
	for _, v := range s.cl {
		if strings.Contains(v.UserName, "@@") {
			clarray = append(clarray, v)
		}
	}
	return clarray
}

// GetStrangers: not group contact and not StarFriend and not special contact
func (s *ContactManager) GetStrangers() []*User {
	clarray := make([]*User, 0)
	for _, v := range s.cl {
		if !strings.Contains(v.UserName, "@@") &&
			v.StarFriend == 0 &&
			v.VerifyFlag&8 == 0 &&
			!SpecialContact[v.UserName] {
			clarray = append(clarray, v)
		}
	}
	return clarray
}

// GetContactByName: get contacts by User.NickName | User.RemarkName | User.DisplayName
func (s *ContactManager) GetContactsByName(sig string) []*User {
	clarray := make([]*User, 0)
	for _, v := range s.cl {
		if v.NickName == sig || v.RemarkName == sig || v.DisplayName == sig {
			clarray = append(clarray, v)
		}
	}
	return clarray
}

// GetContactByPYQuanPin: get contact by User.PYQuanPin | User.RemarkPYQuanPin
func (s *ContactManager) GetContactByPYQuanPin(sig string) *User {
	for _, v := range s.cl {
		if v.PYQuanPin == sig || v.RemarkPYQuanPin == sig {
			return v
		}
	}
	return nil
}

// GetAll: get all contacts
func (s *ContactManager) GetAll() []*User {
	return s.cl
}
