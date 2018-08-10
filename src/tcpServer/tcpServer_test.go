package tcpServer

import (
	"testing"
)

func TestLoginCheck(t *testing.T) {
	var l Logic
	reply := ""
	err := l.LoginCheck(&Col{
		Account:  "hiram1",
		Password: "hiram",
	}, &reply)
	if err != nil {
		t.Error(err)
	} else {
		t.Log("Login check fail")
	}
}
