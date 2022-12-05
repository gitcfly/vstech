package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/apex/gateway"
)

var (
	glog                   = log.New(os.Stderr, "loger:", log.Lshortfile)
	port                   = flag.Int("port", -1, "specify a port")
	_LAMBDA_SERVER_PORT    = ""
	AWS_LAMBDA_RUNTIME_API = ""
)

func main() {
	flag.Parse()
	http.HandleFunc("/api/readf", readf)
	http.HandleFunc("/api/writef", writef)
	http.HandleFunc("/api/net", net)
	_LAMBDA_SERVER_PORT = os.Getenv("_LAMBDA_SERVER_PORT")
	AWS_LAMBDA_RUNTIME_API = os.Getenv("AWS_LAMBDA_RUNTIME_API")
	if *port == -1 {
		glog.Fatal(gateway.ListenAndServe("", nil))
		return
	}
	http.Handle("/", http.FileServer(http.Dir("./view")))
	portStr := fmt.Sprintf(":%d", *port)
	glog.Fatal(http.ListenAndServe(portStr, nil))
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
	glog.Println(_LAMBDA_SERVER_PORT)
	glog.Println(AWS_LAMBDA_RUNTIME_API)
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
