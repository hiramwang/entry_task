package tcpServer

import (
	"strconv"
	"strings"
	"testing"
	"time"
)

var l Logic

var reply = ""

func TestLoginCheck(t *testing.T) {
	err := l.LoginCheck(&Col{
		Account:  "hiram1",
		Password: "hiram",
	}, &reply)
	if err != nil {
		t.Error(err)
	} else {
		t.Log("Login check succ")
	}
}

func TestRegister(t *testing.T) {
	err := l.Register(&Col{
		Account:  "u_" + strconv.FormatInt(time.Now().UnixNano(), 10),
		Password: "hiram",
		Nickname: "hiram",
	}, &reply)
	if err != nil {
		t.Error(err)
	} else {
		t.Log("register succ")
	}
}

func TestChangeNickname(t *testing.T) {
	err := l.UpdateNickname(&Col{
		Account:  "hiram1",
		Nickname: strings.TrimPrefix(strconv.FormatInt(time.Now().UnixNano(), 10), "1534"),
	}, &reply)
	if err != nil {
		t.Error(err)
	} else {
		t.Log("update name succ")
	}
}
