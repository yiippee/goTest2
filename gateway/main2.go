package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	dis "goTest2/gateway/discovery"
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
	// 获取所有的服务列表
	if r.URL.Path == "/list" {
		for _, v := range this.m.Nodes {
			fmt.Fprintf(w, "%v\n", *v)
		}

		return
	}
	// 选择后端服务器策略
	if len(this.m.Keys) == 0 {
		w.Write([]byte("no service!"))
		return
	}
	// todo 获取需要加锁
	key := this.m.Keys[rand.Int()%len(this.m.Keys)]
	addr := this.m.Nodes[key].Info.IP

	remote, err := url.Parse("http://" + addr)
	if err != nil {
		panic(err)
	}

	// NewSingleHostReverseProxy 需要池化吗？ 感觉不需要，NewSingleHostReverseProxy也没做多少事
	// 虽然每一次都需要重新创建。但是又感觉可以池化啊
	proxy := httputil.NewSingleHostReverseProxy(remote) // 这个也可以代理websocket？？？感觉是可以upgrade的
	//如果代理出错，则转向其他后端服务，并检查这个出错服务是否正常，如果不正常则踢出iplist
	//proxy.ErrorHandler = func(writer http.ResponseWriter, request *http.Request, e error) {
	//	robin.Del(addr)
	//}
	proxy.ServeHTTP(w, r)
}

func startServer2(m *dis.Master) {
	// 被代理的服务器host和port
	h := &handle2{m}

	r := gin.New()
	r.GET("/ws", func(c *gin.Context) {
		fmt.Fprintf(c.Writer, "%s", "hello world.")
	})

	v1 := r.Group("/v1") // 所有以v1开头的，全部路由到这里。可以动态增加吗？
	{
		v1.GET("/*key", func(c *gin.Context) {
			//v2 := r.Group("/v2")
			//{
			//	v2.GET("/*key", func(c *gin.Context) {
			//		fmt.Fprintf(c.Writer, "%s", "this is v2.")
			//	})
			//}
			//
			//fmt.Fprintf(c.Writer, "%s", "hello world group.")
			// 选择后端服务器策略

			remote, err := url.Parse("http://" + "127.0.0.1:12003")
			if err != nil {
				panic(err)
			}

			// NewSingleHostReverseProxy 需要池化吗？ 感觉不需要，NewSingleHostReverseProxy也没做多少事
			// 虽然每一次都需要重新创建。但是又感觉可以池化啊
			proxy := NewSingleHostReverseProxy(remote) // 这个也可以代理websocket？？？感觉是可以upgrade的
			//如果代理出错，则转向其他后端服务，并检查这个出错服务是否正常，如果不正常则踢出iplist
			//proxy.ErrorHandler = func(writer http.ResponseWriter, request *http.Request, e error) {
			//	robin.Del(addr)
			//}
			proxy.ServeHTTP(c.Writer, c.Request)
		})
	}

	r.Run(":5000")
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

func serviceStart(serviceName string, ip string) {
	// 注册的信息
	serviceInfo := dis.ServiceInfo{
		IP: ip,
	}

	s, err := dis.NewService(serviceName, serviceInfo, []string{
		"http://127.0.0.1:2379",
	})

	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("name:%s, ip:%s\n", s.Name, s.Info.IP)

	go s.Start()
}
func main() {
	m, err := dis.NewMaster([]string{
		"http://127.0.0.1:2379",
	}, "services/")

	if err != nil {
		log.Fatal(err)
	}

	// 启动具体的服务
	serviceStart("service1", "127.0.0.1:12003")
	//serviceStart("service2", "http://127.0.0.1:12002")
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
