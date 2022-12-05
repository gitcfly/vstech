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
	port = flag.Int("port", -1, "specify a port")
)

func main() {
	flag.Parse()
	http.HandleFunc("/api/readf", readf)
	http.HandleFunc("/api/writef", writef)
	http.HandleFunc("/api/net", net)
	if *port == -1 {
		log.Fatal(gateway.ListenAndServe("", nil))
		return
	}
	http.Handle("/", http.FileServer(http.Dir("./view")))
	portStr := fmt.Sprintf(":%d", *port)
	log.Fatal(http.ListenAndServe(portStr, nil))
}

func net(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	rsp, err := http.Get("https://www.baidu.com/")
	if err != nil {
		w.Write([]byte(fmt.Sprintf(`{"code":200,"msg":"net","err":"%v"}`, err.Error())))
		return
	}
	defer rsp.Body.Close()
	bodys, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		w.Write([]byte(fmt.Sprintf(`{"code":200,"msg":"net","err":"%v"}`, err.Error())))
		return
	}
	w.Write([]byte(fmt.Sprintf(`{"code":200,"msg":"net","path":"%v","body":"%v"}`, string(bodys))))
}

func readf(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	dir, err := os.Getwd()
	if err != nil {
		log.Println(err)
		return
	}
	log.Println(dir)

	files, err := ioutil.ReadDir(dir)
	if err != nil {
		log.Println(err)
		return
	}
	// 获取文件，并输出它们的名字
	for _, file := range files {
		log.Println(file.Name())
	}

	fs, err := ioutil.ReadFile("config/config.txt")
	if err != nil {
		w.Write([]byte(fmt.Sprintf(`{"code":200,"msg":"readf","err":"%v"}`, err)))
		return
	}
	w.Write([]byte(fmt.Sprintf(`{"code":200,"msg":"code","content":"%v"}`, string(fs))))
}

func writef(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	fs, err := os.Create("create_new_file.txt")
	if err != nil {
		w.Write([]byte(fmt.Sprintf(`{"code":200,"msg":"writef","err":"%v"}`, err)))
		return
	}
	n, err := fs.Write([]byte("this is new text"))
	if err != nil {
		w.Write([]byte(fmt.Sprintf(`{"code":200,"msg":"writef","err":"%v"}`, err)))
		return
	}
	if n == 0 {
		w.Write([]byte(fmt.Sprintf(`{"code":200,"msg":"writef","err":"write content is zero"}`)))
		return
	}
	w.Write([]byte(fmt.Sprintf(`{"code":200,"msg":"feed","path":"%v","wirten":"%v"}`, n)))
}
