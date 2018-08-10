package main

import (
	"httpServer"
	"math/rand"
	"runtime"
	"tcpServer"
	"time"
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	rand.Seed(time.Now().UnixNano())

	//if err := conf.InitSetting(); err != nil {
	//		panic(err)
	//}
}

func main() {

	//go api.FrontServer()
	httpEtc := &httpServer.HttpEtc{
		FrontListen:  "127.0.0.1:8888",
		ReadTimeout:  100,
		WriteTimeout: 100,
	}

	go httpServer.FrontServer(httpEtc)
	go tcpServer.RpcServer()

	select {}
}
