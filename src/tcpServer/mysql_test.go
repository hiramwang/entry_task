package tcpServer

import (
	"strconv"
	"testing"
	"time"
)

func TestInsertUser(t *testing.T) {
	UnixStr := strconv.FormatInt(time.Now().Unix(), 10)
	err := MC.InsertUser(&Col{
		Account:  "gotest_" + UnixStr,
		Password: "gotest",
		Nickname: "gotest",
	})
	if err != nil {
		t.Error(err)
	} else {
		t.Log("Mysql func InsertUser Succ")
	}
}

func TestGetUser(t *testing.T) {
	res := MC.GetUser("hiram")
	if res == nil {
		t.Error("No this user")
	} else {
		t.Log("Mysql func GetUser Succ")
	}
}

func TestUpdateNickname(t *testing.T) {
	UnixStr := strconv.FormatInt(time.Now().Unix(), 10)
	err := MC.UpdateNickname("update"+UnixStr, "hiram")
	if err != nil {
		t.Error(err)
	} else {
		t.Log("Mysql func Update Succ")
	}
}
