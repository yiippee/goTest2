package discovery

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/coreos/etcd/clientv3"
	"log"
	"net/http"
	"time"
)

//the detail of service
type ServiceInfo struct {
	IP string
}

type Service struct {
	Name    string
	Info    ServiceInfo
	stop    chan error
	leaseid clientv3.LeaseID
	client  *clientv3.Client
}

func NewService(name string, info ServiceInfo, endpoints []string) (*Service, error) {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   endpoints,
		DialTimeout: 2 * time.Second,
	})

	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	return &Service{
		Name:   name,
		Info:   info,
		stop:   make(chan error),
		client: cli,
	}, err
}

func (s *Service) Start() error {

	go func() {
		time.Sleep(time.Second * 10)
		// s.stop <- errors.New("time out.")
	}()
	// 开启自己的服务
	mux := http.NewServeMux()
	mux.HandleFunc("/v1/test", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello,world."))
	})
	mux.HandleFunc("/v1/listKeys", func(w http.ResponseWriter, r *http.Request) {
		config := clientv3.Config{
			Endpoints:   []string{"127.0.0.1:2379"},
			DialTimeout: time.Duration(1000) * time.Millisecond,
		}
		var client *clientv3.Client
		var err error
		if client, err = clientv3.New(config); err != nil {
			return
		}
		kv := clientv3.NewKV(client)
		if getResp, err := kv.Get(context.TODO(), "services/", clientv3.WithPrefix()); err != nil {
			fmt.Fprintf(w, "no keys")
			return
		} else {
			fmt.Println(getResp.Kvs)
			fmt.Fprintln(w, getResp.Kvs)
		}
	})
	go http.ListenAndServe(s.Info.IP, mux)

	ch, err := s.keepAlive()
	if err != nil {
		log.Fatal(err)
		return err
	}

	for {
		select {
		case err := <-s.stop:
			s.revoke()
			log.Println(err)
			return err
		case <-s.client.Ctx().Done():
			return errors.New("server closed")
		case _, ok := <-ch:
			if !ok {
				log.Println("keep alive channel closed")
				s.revoke()
				return nil
			} else {
				//log.Printf("Recv reply from service: %s, lease ID: %x ttl:%d", s.Name, s.leaseid, ka.TTL)
			}
		}
	}
}

func (s *Service) Stop() {
	s.stop <- nil
}

func (s *Service) keepAlive() (<-chan *clientv3.LeaseKeepAliveResponse, error) {

	info := &s.Info

	key := "services/" + s.Name
	value, _ := json.Marshal(info)

	// minimum lease TTL is 5-second
	resp, err := s.client.Grant(context.TODO(), 5)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	_, err = s.client.Put(context.TODO(), key, string(value), clientv3.WithLease(resp.ID))
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	s.leaseid = resp.ID

	return s.client.KeepAlive(context.TODO(), resp.ID)
}

func (s *Service) revoke() error {

	_, err := s.client.Revoke(context.TODO(), s.leaseid)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("servide:%s stop\n", s.Name)
	return err
}
