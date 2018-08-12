package httpServer

import (
	"conf"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"net/rpc"
	"os"
	"reflect"
	"strings"
	"tcpServer"
)

var (
	RpcClient *rpc.Client
)

type HttpResponse struct {
	Code   string
	Reason string
}

type Dispatcher struct{}

func FrontServer() {
	s := &http.Server{
		Addr: conf.HttpAddr,
	}

	http.HandleFunc("/", BaseHandler)
	http.HandleFunc("/index", IndexHandler)

	err := s.ListenAndServe()
	if err != nil {
		panic("Http server start fail: " + err.Error())
	}

}

func BaseHandler(w http.ResponseWriter, req *http.Request) {
	conf.Log.Trace("httpServer", "request", req.RequestURI)

	err := req.ParseForm()
	if err != nil {
		_, _ = w.Write(GetErrRes("4000"))
		return
	}

	if RpcClient == nil && DialRpc() != nil {
		conf.Log.Error("httpServer", "BaseHandler", "DialRpc fail")
		_, _ = w.Write(GetErrRes("5001"))
		return
	}

	fmt.Println(req.Form)

	account := req.FormValue("account")
	if account == "" {
		_, _ = w.Write(GetErrRes("4001"))
		return
	}

	dispatcher := reflect.ValueOf(&Dispatcher{})
	path := strings.Trim(req.URL.Path, "/")
	funcName := strings.Title(path) + "Handler"

	if req.Method == http.MethodGet {
		method := dispatcher.MethodByName(funcName)
		method.Call([]reflect.Value{
			reflect.ValueOf(w),
			reflect.ValueOf(&tcpServer.Col{
				Account:  account,
				Password: req.FormValue("password"),
				Nickname: req.FormValue("nickname"),
			})})
	} else {
		req.ParseMultipartForm(32 << 20)
		file, _, err := req.FormFile("photo")
		if err != nil {
			conf.Log.Error("httpServer", "BaseHandler", err.Error())
			return
		}
		defer file.Close()
		CheckDirExist()
		f, err := os.OpenFile(strings.Join([]string{conf.PhotoPath, account}, "/"), os.O_WRONLY|os.O_CREATE, 0666)
		io.Copy(f, file)
	}
}

func IndexHandler(w http.ResponseWriter, req *http.Request) {
	t, err := template.ParseFiles(conf.TemplatePath + "index.html")
	if t == nil {
		conf.Log.Error("httpServer", "IndexHandler", err.Error())
		return
	}
	t.Execute(w, "Hello")
}

func (h *Dispatcher) LoginHandler(w http.ResponseWriter, col *tcpServer.Col) {

	if col.Password == "" {
		_, _ = w.Write(GetErrRes("4002"))
		conf.Log.Trace("httpServer", "LoginHandler", conf.ErrType["4002"])
		return
	}

	var reply string
	err := RpcClient.Call("Logic.LoginCheck", col, &reply)
	if err != nil {
		_, _ = w.Write(GetErrRes(err.Error()))
		conf.Log.Error("httpServer", "LoginHandler", err.Error())
		return
	}

	t, _ := template.ParseFiles(conf.TemplatePath + "login.html")
	if t == nil {
		_, _ = w.Write(GetErrRes("5002"))
		conf.Log.Error("httpServer", "IndexHandler", err.Error())
	} else {
		t.Execute(w, col.Account)
	}
}

func (h *Dispatcher) RegisterHandler(w http.ResponseWriter, col *tcpServer.Col) {

	if col.Password == "" {
		_, _ = w.Write(GetErrRes("4002"))
		conf.Log.Trace("httpServer", "LoginHandler", conf.ErrType["4002"])
		return
	}

	var reply string
	err := RpcClient.Call("Logic.Register", col, &reply)
	if err != nil {
		_, _ = w.Write(GetErrRes(err.Error()))
		conf.Log.Error("httpServer", "LoginHandler", err.Error())
	}
}

func (h *Dispatcher) ChangeNameHandler(w http.ResponseWriter, col *tcpServer.Col) {

	if col.Nickname == "" {
		_, _ = w.Write(GetErrRes("4003"))
		conf.Log.Trace("httpServer", "ChangeNameHandler", conf.ErrType["4003"])
		return
	}

	var reply string
	err := RpcClient.Call("Logic.UpdateNickname", col, &reply)
	if err != nil {
		_, _ = w.Write(GetErrRes(err.Error()))
		conf.Log.Error("httpServer", "ChangeNameHandler", err.Error())
	}
}

func DialRpc() error {
	var err error
	RpcClient, err = rpc.DialHTTP("tcp", conf.TcpAddr+conf.TcpPort)
	return err
}

func GetErrRes(code string) []byte {
	res := &HttpResponse{
		Code:   code,
		Reason: conf.ErrType[code],
	}
	s, _ := json.Marshal(res)
	return s
}

func CheckDirExist() {
	if _, err := os.Stat(conf.PhotoPath); os.IsNotExist(err) {
		os.Mkdir(conf.PhotoPath, os.ModePerm)
	}
}
