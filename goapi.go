package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/elazarl/goproxy"
)

var (
	glog = log.New(os.Stderr, "loger:", log.Lshortfile)
)

func main() {
	lamdaServerPort := os.Getenv("_LAMBDA_SERVER_PORT")
	proxy := goproxy.NewProxyHttpServer()
	proxy.Verbose = true
	proxy.NonproxyHandler = http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				glog.Println(err)
			}
		}()
		stime := time.Now()
		var envStr = strings.Join(os.Environ(), "\n")
		cost := time.Since(stime)
		http.Error(w, fmt.Sprintf("%v \nThis is a proxy server. Does not respond to non-proxy requests.\n reqeust cost=%v", envStr, cost), 500)
	})
	proxy.OnRequest().DoFunc(
		func(r *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
			glog.Println(ctx.Req.Proto, ctx.Req.Method, ctx.Req.URL.String())
			return r, nil
		})
	glog.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", lamdaServerPort), proxy))
}

func net(w http.ResponseWriter, r *http.Request) {
	glog.Println(r.Method, r.URL.String())
	w.Header().Set("content-type", "application/json")
	rsp, err := http.Get("https://www.baidu.com/")
	if err != nil {
		glog.Println(err)
		return
	}
	defer rsp.Body.Close()
	bodys, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		glog.Println(err)
		return
	}
	w.Write([]byte(fmt.Sprintf(`{"code":200,"msg":"net","path":"%v","body":"%v"}`, string(bodys))))
}

func readf(w http.ResponseWriter, r *http.Request) {
	glog.Println(r.Method, r.URL.String())
	w.Header().Set("content-type", "application/json")
	dir, err := os.Getwd()
	if err != nil {
		glog.Println(err)
		return
	}
	glog.Println(dir)
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		glog.Println(err)
		return
	}
	// 获取文件，并输出它们的名字
	for _, file := range files {
		glog.Println(file.Name())
	}
	fs, err := ioutil.ReadFile("config.txt")
	if err != nil {
		glog.Println(err)
		return
	}
	w.Write([]byte(fmt.Sprintf(`{"code":200,"msg":"code","content":"%v"}`, string(fs))))
}

func writef(w http.ResponseWriter, r *http.Request) {
	glog.Println(r.Method, r.URL.String())
	w.Header().Set("content-type", "application/json")
	fs, err := os.Create("create_new_file.txt")
	if err != nil {
		glog.Println(err)
		return
	}
	n, err := fs.Write([]byte("this is new text"))
	if err != nil {
		glog.Println(err)
		return
	}
	if n == 0 {
		glog.Println("write content is zero")
		return
	}
	w.Write([]byte(fmt.Sprintf(`{"code":200,"msg":"feed","path":"%v","wirten":"%v"}`, fs.Name(), n)))
}
