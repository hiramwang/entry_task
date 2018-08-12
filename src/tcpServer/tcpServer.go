package tcpServer

import (
	"conf"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-redis/redis"
	"net"
	"net/http"
	"net/rpc"
	"time"
)

type Logic int

var RC *redis.Client

const ExpireT = time.Second * 3600

func init() {
	RC = redis.NewClient(&redis.Options{
		Addr:     conf.RedisAddr,
		Password: conf.RedisPass,
	})

	_, err := RC.Ping().Result()
	if err != nil {
		panic("redis error")
	}
}

func (l *Logic) LoginCheck(args *Col, reply *string) error {
	UserData := &Col{}
	UserDataString, err := RC.Get(args.Account).Result()

	if err != nil && err.Error() != "redis: nil" {
		conf.Log.Error("TcpServer", "LoginCheck", "Redis error", err.Error())
	}

	if UserDataString != "" {
		err = json.Unmarshal([]byte(UserDataString), UserData)
		if err != nil {
			conf.Log.Error("TcpServer", "LoginCheck", "Json unmarshal error", err.Error())
		}
	} else {
		UserData = MC.GetUser(args.Account)
		*reply = "loginRpcReply"
		if UserData == nil {
			return errors.New("4001")
		} else {
			dataBt, err := json.Marshal(UserData)
			err = RC.Set(UserData.Account, string(dataBt), ExpireT).Err()
			if err != nil {
				conf.Log.Error("TcpServer", "LoginCheck", "Redis set error", err.Error())
			}
		}
	}

	if UserData.Password != args.Password {
		return errors.New("4002")
	}

	return nil
}

func (l *Logic) Register(args *Col, reply *string) error {
	err := MC.InsertUser(args)
	if err != nil {
		return errors.New("Insert mysql error: " + err.Error())
	}
	return nil
}

func (l *Logic) UpdateNickname(args *Col, reply *string) error {
	err := MC.UpdateNickname(args.Nickname, args.Account)
	if err != nil {
		return errors.New("Update mysql error: " + err.Error())
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
	l, err := net.Listen("tcp", conf.TcpPort)
	if err != nil {
		panic(err)
	}
	go http.Serve(l, nil)
}
