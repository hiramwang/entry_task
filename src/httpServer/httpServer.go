package httpServer

import (
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
	"time"
)

var (
	TemplatePath = "./templates/"
	RpcClient    *rpc.Client
	PhotoPath    = "./photo"
)

type HttpEtc struct {
	FrontListen  string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

type HttpResponse struct {
	Code   string
	Reason string
}

type Dispatcher struct{}

func FrontServer(etc *HttpEtc) {
	s := &http.Server{
		Addr: etc.FrontListen,
	}

	http.HandleFunc("/", BaseHandler)
	http.HandleFunc("/index", IndexHandler)

	err := s.ListenAndServe()
	if err != nil {
		fmt.Println("  panic  ")
		panic(err)
	}

}

func BaseHandler(w http.ResponseWriter, req *http.Request) {
	return

	err := req.ParseForm()
	if err != nil {
		_, _ = w.Write(GetErrRes("4000"))
		return
	}

	if RpcClient == nil && DialRpc() != nil {
		// rpc error
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
			fmt.Println(err)
			return
		}
		defer file.Close()
		CheckDirExist()
		f, err := os.OpenFile(strings.Join([]string{PhotoPath, account}, "/"), os.O_WRONLY|os.O_CREATE, 0666)
		io.Copy(f, file)
	}
}

func IndexHandler(w http.ResponseWriter, req *http.Request) {
	t, err := template.ParseFiles(TemplatePath + "index.html")
	if t == nil {
		panic(err)
	}
	t.Execute(w, "Hello")
}

func (h *Dispatcher) LoginHandler(w http.ResponseWriter, col *tcpServer.Col) {

	if col.Password == "" {
		_, _ = w.Write(GetErrRes("4002"))
		return
	}

	var reply string
	err := RpcClient.Call("Logic.LoginCheck", col, &reply)
	if err != nil {
		_, _ = w.Write(GetErrRes(err.Error()))
		return
	}

	t, _ := template.ParseFiles(TemplatePath + "login.html")
	if t == nil {
		_, _ = w.Write(GetErrRes("5002"))
	} else {
		t.Execute(w, col.Account)
	}
}

func (h *Dispatcher) RegisterHandler(w http.ResponseWriter, col *tcpServer.Col) {

	if col.Password == "" {
		_, _ = w.Write(GetErrRes("4002"))
	}

	var reply string
	err := RpcClient.Call("Logic.Register", col, &reply)
	if err != nil {
		_, _ = w.Write(GetErrRes(err.Error()))
	}
}

func (h *Dispatcher) ChangeNameHandler(w http.ResponseWriter, col *tcpServer.Col) {

	if col.Nickname == "" {
		_, _ = w.Write(GetErrRes("4003"))
	}

	var reply string
	err := RpcClient.Call("Logic.UpdateNickname", col, &reply)
	if err != nil {
		_, _ = w.Write(GetErrRes(err.Error()))
	}
}

func DialRpc() error {
	var err error
	RpcClient, err = rpc.DialHTTP("tcp", "127.0.0.1:8880")
	return err
}

func GetErrRes(code string) []byte {
	res := &HttpResponse{
		Code:   code,
		Reason: tcpServer.ErrType[code],
	}
	s, _ := json.Marshal(res)
	return s
}

func CheckDirExist() {
	if _, err := os.Stat(PhotoPath); os.IsNotExist(err) {
		os.Mkdir(PhotoPath, os.ModePerm)
	}
}
