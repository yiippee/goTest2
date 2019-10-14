package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

var robin = NewWeightedRR(RR_NGINX)

type handle2 struct {
	addrs []string
}

func (this *handle2) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// 选择后端服务器策略
	addr := robin.Next().(string)
	remote, err := url.Parse("http://" + addr)
	if err != nil {
		panic(err)
	}
	proxy := httputil.NewSingleHostReverseProxy(remote)
	// proxy.ErrorHandler =   如果代理出错，则转向其他后端服务，并检查这个出错服务是否正常，如果不正常则踢出iplist
	proxy.ServeHTTP(w, r)
}

func startServer2() {
	// 被代理的服务器host和port
	h := &handle2{}
	h.addrs = []string{"127.0.0.1:28081", "127.0.0.1:28082"}

	w := 1
	for _, e := range h.addrs {
		robin.Add(e, w)
		w++
	}
	err := http.ListenAndServe(":28080", h)
	if err != nil {
		log.Fatalln("ListenAndServe: ", err)
	}
}

func main() {
	startServer2()
}
