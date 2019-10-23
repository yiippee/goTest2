package discovery

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/coreos/etcd/clientv3"
	"log"
	"sync"
	"time"
)

type Master struct {
	Path   string
	Keys   []string
	Nodes  map[string]*Node
	Client *clientv3.Client
	sync.Mutex
}

type ServiceState int

const (
	ON_LINE ServiceState = iota
	OFF_LINE
	UPGRADING // 升级维护
)

//node is a client
type Node struct {
	State ServiceState
	Key   string
	Info  ServiceInfo
}

func NewMaster(endpoints []string, watchPath string) (*Master, error) {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   endpoints,
		DialTimeout: time.Second,
	})

	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	master := &Master{
		Path:   watchPath,
		Nodes:  make(map[string]*Node),
		Client: cli,
	}

	go master.WatchNodes()
	return master, err
}

func (m *Master) AddNode(key string, info *ServiceInfo) {
	m.Lock()
	defer m.Unlock()
	node := &Node{
		State: ON_LINE,
		Key:   key,
		Info:  *info,
	}

	m.Nodes[node.Key] = node
	m.Keys = append(m.Keys, key)
}

func (m *Master) DelNode(key string) {
	m.Lock()
	defer m.Unlock()

	delete(m.Nodes, key)

	for k, v := range m.Keys {
		if v == key {
			m.Keys = append(m.Keys[:k], m.Keys[k+1:]...)
		}
	}
}

func GetServiceInfo(ev *clientv3.Event) *ServiceInfo {
	info := &ServiceInfo{}
	err := json.Unmarshal([]byte(ev.Kv.Value), info)
	if err != nil {
		log.Println(err)
	}
	return info
}

func (m *Master) WatchNodes() {
	rch := m.Client.Watch(context.Background(), m.Path, clientv3.WithPrefix())
	for wresp := range rch {
		for _, ev := range wresp.Events {
			switch ev.Type {
			case clientv3.EventTypePut:
				fmt.Printf("[%s] %q : %q\n", ev.Type, ev.Kv.Key, ev.Kv.Value)
				info := GetServiceInfo(ev)
				m.AddNode(string(ev.Kv.Key), info)
			case clientv3.EventTypeDelete:
				fmt.Printf("[%s] %q : %q\n", ev.Type, ev.Kv.Key, ev.Kv.Value)
				m.DelNode(string(ev.Kv.Key))
			}
		}
	}
}
