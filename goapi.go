package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/apex/gateway"
)

var (
	count = 0
	port  = flag.Int("port", -1, "specify a port")
)

func main() {
	flag.Parse()
	count++
	http.HandleFunc("/api/code", code)
	http.HandleFunc("/api/feed", feed)
	http.HandleFunc("/api/root", root)
	listener := gateway.ListenAndServe
	portStr := "n/a"
	if *port != -1 {
		portStr = fmt.Sprintf(":%d", *port)
		listener = http.ListenAndServe
		http.Handle("/", http.FileServer(http.Dir("./public")))
	}
	log.Fatal(listener(portStr, nil))
}

func root(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	rsp, err := http.Get("https://www.baidu.com/")
	if err != nil {
		w.Write([]byte(fmt.Sprintf(`{"code":200,"msg":"root","err":"%v"}`, err.Error())))
		return
	}
	defer rsp.Body.Close()
	bodys, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		w.Write([]byte(fmt.Sprintf(`{"code":200,"msg":"root","err":"%v"}`, err.Error())))
		return
	}
	w.Write([]byte(fmt.Sprintf(`{"code":200,"msg":"root","path":"%v","body":"%v"}`, string(bodys))))
}

func code(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	w.Write([]byte(fmt.Sprintf(`{"code":200,"msg":"code","path":"%v"}`, r.URL.String())))
}

func feed(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	w.Write([]byte(fmt.Sprintf(`{"code":200,"msg":"feed","path":"%v"}`, r.URL.String())))
}
