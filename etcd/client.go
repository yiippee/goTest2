package main

import (
	"context"
	"log"
	"time"

	"go.etcd.io/etcd/clientv3"
)

func main() {
	cfg := clientv3.Config{
		Endpoints:               []string{"http://127.0.0.1:2379"},
		Transport:               clientv3.DefaultTransport,
		// set timeout per request to fail fast when the target endpoint is unavailable
		HeaderTimeoutPerRequest: time.Second,
	}
	c, err := clientv3.New(cfg)
	if err != nil {
		log.Fatal(err)
	}

	kapi := clientv3.NewKeysAPI(c)
	log.Print("Setting '/foo' key with 'bar' value")
	resp, err := kapi.Set(context.Background(), "/foo", "bar", nil)
	if err != nil {
		log.Fatal(err)
	} else {
		log.Printf("Set is done. Metadata is %q\n", resp)
	}

	resp, err = kapi.Get(context.Background(), "/foo", nil)
	if err != nil {
		log.Fatal(err)
	} else {
		log.Printf("Get is done. Metadata is %q\n", resp)
		log.Printf("%q key has %q value\n", resp.Node.Key, resp.Node.Value)
	}
}
