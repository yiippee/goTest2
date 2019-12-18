package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/coreos/etcd/clientv3"
)

func main() {
	cfg := clientv3.Config{
		Endpoints: []string{
			"http://172.20.200.17:9002",
			"http://172.20.200.17:9004",
			"http://172.20.200.17:9006",
		},
		DialTimeout: 5 * time.Second,
	}

	client, err := clientv3.New(cfg)
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	resp, err := client.Put(ctx, "/test/sample_key", "sample_value2")
	cancel()
	if err != nil {
		// handle error!
		switch err {
		case context.Canceled:
			log.Fatalf("ctx is canceled by another routine: %v", err)
		case context.DeadlineExceeded:
			log.Fatalf("ctx is attached with a deadline is exceeded: %v", err)
		default:
			log.Fatalf("bad cluster endpoints, which are not etcd servers: %v", err)
		}
	}
	fmt.Println(resp.Header.Revision)
	// use the response
	r, err := client.Get(context.TODO(), "/test/sample_key", clientv3.WithRev(3670))
	fmt.Println(r)
}
