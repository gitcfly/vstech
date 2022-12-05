package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/apex/gateway"
)

var (
	port = flag.Int("port", -1, "specify a port")
)

func main() {
	flag.Parse()
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
	w.Write([]byte(fmt.Sprintf(`{"code":200,"msg":"root","path":"%v"}`, r.URL.String())))
}

func code(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	w.Write([]byte(fmt.Sprintf(`{"code":200,"msg":"code","path":"%v"}`, r.URL.String())))
}

func feed(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	w.Write([]byte(fmt.Sprintf(`{"code":200,"msg":"feed","path":"%v"}`, r.URL.String())))
}

//
//func feed(w http.ResponseWriter, r *http.Request) {
//	url := "https://www.zhihu.com/api/v3/feed/topstory/hot-lists/total?limit=50&desktop=true"
//	method := "GET"
//
//	client := &http.Client{}
//	req, err := http.NewRequest(method, url, nil)
//
//	if err != nil {
//		fmt.Println(err)
//		return
//	}
//
//	res, err := client.Do(req)
//	if err != nil {
//		fmt.Println(err)
//		return
//	}
//	defer res.Body.Close()
//
//	body, err := ioutil.ReadAll(res.Body)
//	if err != nil {
//		fmt.Println(err)
//		return
//	}
//
//	w.Header().Set("content-type", "application/json")
//	w.Write(body)
//}
