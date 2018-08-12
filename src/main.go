package main

import (
	"conf"
	"httpServer"
	"math/rand"
	"runtime"
	"tcpServer"
	"time"
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	rand.Seed(time.Now().UnixNano())

	if err := conf.InitSetting(); err != nil {
		panic(err)
	}
}

func main() {

	go httpServer.FrontServer()
	go tcpServer.RpcServer()

	select {}
}
