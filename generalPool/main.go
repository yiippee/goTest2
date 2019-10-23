package main

import (
	"fmt"
	"github.com/silenceper/pool"
	"net/url"
	"time"

	"github.com/gorilla/websocket"
)

func main() {
	//factory 创建连接的方法
	factory := func() (interface{}, error) {
		// return net.Dial("tcp", "127.0.0.1:4000")
		u := url.URL{Scheme: "ws", Host: "127.0.0.1:8080", Path: "/r"}
		conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
		if err != nil {
			return nil, err
		}
		return conn, nil
	}

	//close 关闭连接的方法
	close := func(v interface{}) error {
		return v.(*websocket.Conn).Close()
	}

	//ping 检测连接的方法
	//ping := func(v interface{}) error { return nil }

	//创建一个连接池： 初始化5，最大连接30
	poolConfig := &pool.Config{
		InitialCap: 5,
		MaxCap:     30,
		Factory:    factory,
		Close:      close,
		//Ping:       ping,
		//连接最大空闲时间，超过该时间的连接 将会关闭，可避免空闲时连接EOF，自动失效的问题
		IdleTimeout: 15 * time.Second,
	}
	p, err := pool.NewChannelPool(poolConfig)
	if err != nil {
		fmt.Println("err=", err)
	}

	//从连接池中取得一个连接
	v, err := p.Get()

	//do something
	conn := v.(*websocket.Conn)
	err = conn.WriteMessage(1, []byte("123"))
	if err != nil {
		fmt.Println("err: ", err)
	}
	//将连接放回连接池中
	p.Put(v)

	//释放连接池中的所有连接
	p.Release()

	//查看当前连接中的数量
	current := p.Len()

	fmt.Println(current)
}
