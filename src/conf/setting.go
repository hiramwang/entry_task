package conf

import (
	"errors"
	"os"
	"path/filepath"
)

var (
	TemplatePath = "/src/templates/"
	PhotoPath    = "/photo"
	HttpAddr     = "127.0.0.1:8888"
	TcpAddr      = "127.0.0.1"
	TcpPort      = ":8880"
	RedisAddr    = "127.0.0.1:6379"
	RedisPass    = ""
	MysqlUser    = "hiram"
	MysqlPass    = "hiram"
	MysqlDB      = "mysql"
	MysqlTable   = "shopee_test"
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

var (
	Log *Logger
)

func InitSetting() error {
	rootPath, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		return errors.New("Get rootPath fail: " + err.Error())
	}
	rootPath = rootPath + "/../"

	TemplatePath = filepath.Join(rootPath, TemplatePath)
	PhotoPath = filepath.Join(rootPath, PhotoPath)

	if _, err = os.Stat(PhotoPath); os.IsNotExist(err) {
		err = os.Mkdir(PhotoPath, 0777)
		if err != nil {
			return errors.New("Creat log path fail: " + err.Error())
		}
	}

	logDir := rootPath + "/log"
	logBackDir := rootPath + "/logbak"
	if _, err = os.Stat(logDir); os.IsNotExist(err) {
		err = os.Mkdir(logDir, 0777)
		if err != nil {
			return errors.New("Creat log path fail: " + err.Error())
		}
	}

	if _, err = os.Stat(logBackDir); os.IsNotExist(err) {
		err = os.Mkdir(logBackDir, 0777)
		if err != nil {
			return errors.New("Creat log backup path fail: " + err.Error())
		}
	}

	filename := filepath.Join(logDir, "EntryTask")
	Log, err = NewLogger(filename, "EntryTask", logBackDir)
	if err != nil {
		return err
	}
	return nil
}
