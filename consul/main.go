//package main
//
//import "github.com/hashicorp/consul/api"
//import "fmt"
//
//func main() {
//	// Get a new client
//	client, err := api.NewClient(api.DefaultConfig())
//	if err != nil {
//		panic(err)
//	}
//
//	// Get a handle to the KV API
//	kv := client.KV()
//
//	// PUT a new KV pair
//	p := &api.KVPair{Key: "REDIS_MAXCLIENTS", Value: []byte("1000")}
//	_, err = kv.Put(p, nil)
//	if err != nil {
//		panic(err)
//	}
//
//	// Lookup the pair
//	pair, _, err := kv.Get("REDIS_MAXCLIENTS", nil)
//	if err != nil {
//		panic(err)
//	}
//	fmt.Printf("KV: %v %s\n", pair.Key, pair.Value)
//}

package main

import (
	"fmt"
	"github.com/hashicorp/consul/api/watch"
	"net/http"
	"time"

	consulApi "github.com/hashicorp/consul/api"
)

// 使用consul源码中的watch包监听服务变化
func main() {
	var (
		err    error
		params map[string]interface{}
		plan   *watch.Plan
		ch     chan int
	)
	ch = make(chan int, 1)

	params = make(map[string]interface{})
	params["type"] = "service"
	params["service"] = "test2"
	params["passingonly"] = false
	params["tag"] = "SERVER2"
	plan, err = watch.Parse(params)
	if err != nil {
		panic(err)
	}
	plan.Handler = func(index uint64, result interface{}) {
		if entries, ok := result.([]*consulApi.ServiceEntry); ok {
			fmt.Printf("serviceEntries:%#v \n", entries)
			// your code
			ch <- 1
		}
	}
	go func() {
		// your consul agent addr
		if err = plan.Run("127.0.0.1:8500"); err != nil {
			panic(err)
		}
	}()
	go http.ListenAndServe(":8080", nil)

	time.Sleep(3 * time.Second)

	go register()
	for {
		<-ch
		fmt.Printf("get change\n")
	}
}

func register() {
	var (
		err    error
		client *consulApi.Client
	)
	client, err = consulApi.NewClient(&consulApi.Config{Address: "127.0.0.1:8500"})
	if err != nil {
		panic(err)
	}
	err = client.Agent().ServiceRegister(&consulApi.AgentServiceRegistration{
		ID:   "",
		Name: "test2",
		Tags: []string{"SERVER2"},
		Port: 8080,
		Check: &consulApi.AgentServiceCheck{
			HTTP:     "https://localhost:5000/health",
			Interval: "10s",
			Timeout:  "15s",
		},
	})
	if err != nil {
		panic(err)
	}
}
