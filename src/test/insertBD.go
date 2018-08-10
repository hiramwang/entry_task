package main

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"strconv"
	"tcpServer"
	"time"
)

func main() {
	for i := 0; i < 5; i++ {
		go insert()
	}
	select {}
}

func insert() {
	for true {
		UnixStr := strconv.FormatInt(time.Now().UnixNano(), 10)
		err := tcpServer.MC.InsertUser(&tcpServer.Col{
			Account:  "u_" + UnixStr,
			Password: "gotest",
			Nickname: "gotest",
		})
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println(UnixStr)
		}
	}
}
