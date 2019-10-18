package main

import (
	"fmt"
	dis "goTest/gateway/discovery"
	"log"
	"math/rand"
	"net/http"
	"net/http/httputil"
	"net/url"
)

var robin = NewWeightedRR(RR_NGINX)

type handle2 struct {
	m *dis.Master
}

func (this *handle2) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// 选择后端服务器策略
	if len(this.m.Keys) == 0 {
		w.Write([]byte("no service!"))
		return
	}
	key := this.m.Keys[rand.Int()%len(this.m.Keys)]
	addr := this.m.Nodes[key].Info.IP

	remote, err := url.Parse("http://" + addr)
	if err != nil {
		panic(err)
	}
	proxy := httputil.NewSingleHostReverseProxy(remote)
	//如果代理出错，则转向其他后端服务，并检查这个出错服务是否正常，如果不正常则踢出iplist
	//proxy.ErrorHandler = func(writer http.ResponseWriter, request *http.Request, e error) {
	//	robin.Del(addr)
	//}
	proxy.ServeHTTP(w, r)
}

func startServer2(m *dis.Master) {
	// 被代理的服务器host和port
	h := &handle2{m}

	//w := 1
	//for _, e := range h.addrs {
	//	robin.Add(e, w)
	//	w++
	//}

	err := http.ListenAndServe(":12000", h)
	if err != nil {
		log.Fatalln("ListenAndServe: ", err)
	}
}

func main() {
	m, err := dis.NewMaster([]string{
		"http://127.0.0.1:2379",
	}, "services/")

	if err != nil {
		log.Fatal(err)
	}

	// service
	serviceName := "service-test"
	// 注册的信息
	serviceInfo := dis.ServiceInfo{
		IP: "127.0.0.1:12001",
	}

	s, err := dis.NewService(serviceName, serviceInfo, []string{
		"http://127.0.0.1:2379",
	})

	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("name:%s, ip:%s\n", s.Name, s.Info.IP)

	go s.Start()

	//
	//for {
	//	for k, v := range  m.Nodes {
	//		fmt.Printf("node:%s, ip=%s\n", k, v.Info.IP)
	//	}
	//	fmt.Printf("nodes num = %d\n",len(m.Nodes))
	//	time.Sleep(time.Second * 5)
	//}

	startServer2(m)
}
