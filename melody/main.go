package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gopkg.in/olahol/melody.v1"
)

type Server struct {
	*gin.Engine
	*melody.Melody
}

func main() {
	//r := gin.Default()
	//m := melody.New()
	//
	//r.GET("/", func(c *gin.Context) {
	//	http.ServeFile(c.Writer, c.Request, "index.html")
	//})
	//
	//r.GET("/ws", func(c *gin.Context) {
	//	m.HandleRequest(c.Writer, c.Request)
	//})
	//
	//m.HandleMessage(func(s *melody.Session, msg []byte) {
	//	m.Broadcast(msg)
	//})
	//
	//r.Run(":5000")
	srv := Server{
		Engine: gin.Default(),
		Melody: melody.New(),
	}
	srv.GET("/", func(c *gin.Context) {
		http.ServeFile(c.Writer, c.Request, "index.html")
	})
	srv.GET("/ws", func(c *gin.Context) {
		srv.HandleRequest(c.Writer, c.Request)
	})
	srv.HandleMessage(func(s *melody.Session, msg []byte) {
		// srv.Broadcast(msg)
		srv.BroadcastFilter(msg, func(s *melody.Session) bool {
			if msg[0] == 82 {
				return true
			}
			return false
		})
	})

	srv.Run(":5000")
}
