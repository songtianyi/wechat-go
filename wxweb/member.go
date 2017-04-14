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
	"fmt"
	"github.com/songtianyi/wechat-go/wxweb"
)

type MemberManager struct {
	Group *wxweb.User
}

func CreateMemberManagerFromGroupContact(session *WxWebSession, user *wxweb.User) (*MemberManager, error) {
	b, err := wxweb.WebWxBatchGetContact(session.WxWebCommon, session.WxWebXcg, session.Cookies, []*wxweb.User{user})
	if err != nil {
		return nil, err
	}
	return CreateMemberManagerFromBytes(b)
}

func CreateMemberManagerFromBytes(b []byte) (*MemberManager, error) {
	var gcr wxweb.GroupContactResponse
	if err := json.Unmarshal(b, &gcr); err != nil {
		return nil, err
	}
	if gcr.BaseResponse.Ret != 0 {
		return nil, fmt.Errorf("WebWxBatchGetContact ret=%d", gcr.BaseResponse.Ret)
	}

	if gcr.ContactList == nil || len(gcr.ContactList) < 1 {
		return nil, fmt.Errorf("ContactList empty")
	}

	mm := &MemberManager{
		Group: gcr.ContactList[0],
	}
	return mm, nil
}

func (s *MemberManager) Update() error {
	members := make([]*wxweb.User, len(s.Group.MemberList))
	for i, v := range s.Group.MemberList {
		members[i] = &wxweb.User{
			UserName:        v.UserName,
			EncryChatRoomId: s.Group.UserName,
		}
	}
	b, err := wxweb.WebWxBatchGetContact(WxWebCommon, WxWebXcg, Cookies, members)
	if err != nil {
		return err
	}

	var gcr wxweb.GroupContactResponse
	if err := json.Unmarshal(b, &gcr); err != nil {
		return err
	}
	s.Group.MemberList = gcr.ContactList
	return nil
}

func (s *MemberManager) GetHeadImgUrlByGender(sex int) []string {
	uris := make([]string, 0)
	for _, v := range s.Group.MemberList {
		if v.Sex == sex {
			uris = append(uris, v.HeadImgUrl)
		}
	}
	return uris
}

func (s *MemberManager) GetContactsByGender(sex int) []*wxweb.User {
	contacts := make([]*wxweb.User, 0)
	for _, v := range s.Group.MemberList {
		if v.Sex == sex {
			contacts = append(contacts, v)
		}
	}
	return contacts
}

func (s *MemberManager) GetContactByUserName(username string) *wxweb.User {
	for _, v := range s.Group.MemberList {
		if v.UserName == username {
			return v
		}
	}
	return nil
}
