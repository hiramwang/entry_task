package tcpServer

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-redis/redis"
	"net"
	"net/http"
	"net/rpc"
	"time"
)

var ErrType map[string]string = map[string]string{
	"":     "default error",
	"4000": "Request parse fail",
	"4001": "Account does not exist",
	"4002": "Password invalid",
	"4003": "Nickname invalid",
	"5001": "rpc server error",
	"5002": "httpserver error",
}

type Logic int

var RC *redis.Client

const ExpireT = time.Second * 3600

func init() {
	RC = redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379",
		Password: "",
	})

	_, err := RC.Ping().Result()
	if err != nil {
		fmt.Println("redis error")
	} else {
		fmt.Println("redis ping succ")
	}
}

func (l *Logic) LoginCheck(args *Col, reply *string) error {
	UserData := &Col{}
	UserDataString, err := RC.Get(args.Account).Result()

	if err != nil && err.Error() != "redis: nil" {
		fmt.Println(err)
	}

	if UserDataString != "" {
		err = json.Unmarshal([]byte(UserDataString), UserData)
		if err != nil {
			fmt.Println(err)
		}
	} else {
		fmt.Println("use mysql to check login: " + args.Account)
		UserData = MC.GetUser(args.Account)
		*reply = "loginRpcReply"
		if UserData == nil {
			return errors.New("4001")
		} else {
			dataBt, err := json.Marshal(UserData)
			err = RC.Set(UserData.Account, string(dataBt), ExpireT).Err()
			if err != nil {
				fmt.Println(err)
			}
		}
	}

	return nil
	if UserData.Password != args.Password {
		return errors.New("4002")
	}

	return nil
}

func (l *Logic) Register(args *Col, reply *string) error {
	err := MC.InsertUser(args)
	if err != nil {
		return errors.New("")
	}
	return nil
}

func (l *Logic) UpdateNickname(args *Col, reply *string) error {
	err := MC.UpdateNickname(args.Nickname, args.Account)
	if err != nil {
		return errors.New("")
	}
	return nil
}

func (l *Logic) UpLoadPic(args *Col, reply *string) error {
	return nil
}

func RpcServer() {
	fmt.Println("tcpServer test start")
	logic := new(Logic)
	rpc.Register(logic)
	rpc.HandleHTTP()
	l, err := net.Listen("tcp", ":8880")
	if err != nil {
		panic(err)
	}
	go http.Serve(l, nil)
}
