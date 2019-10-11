package main

import (
	"fmt"
	"log"
	"net"
	"testing"

	"github.com/funny/link"
	"github.com/funny/link/codec"
)

type AddReq struct {
	A, B int
}

type AddRsp struct {
	C int
}

type Server struct{}

func Test(t *testing.T) {
	json := codec.Json()
	json.Register(AddReq{})
	json.Register(AddRsp{})

	listen, err := net.Listen("tcp", "")
	checkErr(err)
	server := link.NewServer(listen, json, 1024, new(Server))
	go server.Serve()
	addr := server.Listener().Addr()

	clientSession, err := link.Dial(addr.Network(), addr.String(), json, 1024)
	checkErr(err)
	go clientSessionLoop(clientSession)

	clientSession2, err := link.Dial(addr.Network(), addr.String(), json, 1024)
	checkErr(err)
	go clientSessionLoop(clientSession2)
	select {

	}
}

func (srv *Server) HandleSession(session *link.Session) {
	for {
		req, err := session.Receive()
		checkErr(err)

		err = session.Send(&AddRsp{
			req.(*AddReq).A + req.(*AddReq).B,
		})

		//
		//s := session.Manager().GetSession(3)
		//if s == nil {
		//	continue
		//}
		//err = s.Send(&AddRsp{
		//	req.(*AddReq).A + req.(*AddReq).B,
		//})

		checkErr(err)
	}
}

func clientSessionLoop(session *link.Session) {
	fmt.Println("session id: ", session.ID())
	for i := 0; i < 10; i++ {
		//err := session.Send(&AddReq{
		//	i, i,
		//})
		//checkErr(err)
		//log.Printf("Send: %d + %d", i, i)

		rsp, err := session.Receive()
		checkErr(err)
		log.Printf("Receive: %d", rsp.(*AddRsp).C)
	}
}

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
