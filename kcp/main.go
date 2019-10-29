package main

import (
	"bytes"
	"fmt"
	_ "io"
	"net"

	"github.com/itfantasy/gonode/utils/timer"
	"github.com/xtaci/kcp-go"
)

func main() {
	testKcp()
	select {}
}

func testKcp() {
	go testKcpServer()
	testKcpClient()
}

func testKcpServer() {
	listen, err := kcp.Listen("0.0.0.0:10086")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer listen.Close()
	for {
		conn, err := listen.Accept()
		if err != nil {
			fmt.Println(err)
			return
		}
		go handleConnS(conn)
	}
}

func handleConnS(conn net.Conn) {
	for {
		fmt.Println("recv -----> ")
		datas := bytes.NewBuffer(nil)
		var buf [512]byte

		n, err := conn.Read(buf[0:])
		fmt.Println(n)
		datas.Write(buf[0:n])
		if err != nil {
			fmt.Println(err.Error())
		}

		fmt.Print("datas : ")
		fmt.Println(string(datas.Bytes()))
		conn.Write(datas.Bytes())
	}
}

func testKcpClient() {
	conn, err := kcp.Dial("127.0.0.1:10086")
	if err != nil {
		fmt.Println(err)
		return
	}
	for {
		fmt.Println("send ------> ")
		ret, err2 := conn.Write([]byte("hello kcp!!"))
		if err2 != nil {
			fmt.Println(err2)
		} else {
			fmt.Println(ret)
		}
		timer.Sleep(1000)
	}
}

func handleConnC(conn net.Conn) {
	for {
		fmt.Println("recv -----> ")
		datas := bytes.NewBuffer(nil)
		var buf [512]byte

		n, err := conn.Read(buf[0:])
		datas.Write(buf[0:n])
		if err != nil {
			fmt.Println(err.Error())
		}

		fmt.Print("datas : ")
		fmt.Println(datas.Bytes())
	}
}
