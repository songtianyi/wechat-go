package wxweb

import (
	"testing"
)

func TestGetLoginAvatar(t *testing.T) {
	expect := "data:img/jpg"
	resp := "window.code=201;window.userAvatar = 'data:img/jpg'"
	avatar, err := GetLoginAvatar(resp)
	if err != nil {
		t.Error(err)
		return
	}
	if avatar != expect {
		t.Errorf("got avatar failed, expect: %s, got: %s", expect, avatar)
		return
	}
}
