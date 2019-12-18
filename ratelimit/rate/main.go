package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
	"net/http"
	"time"
)

var m map[string]*rate.Limiter = make(map[string]*rate.Limiter)

func main() {
	r := gin.Default()

	//Apply only to /ping
	r.GET("/ping", SetUp(1), func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	r.Run(":8087") // listen and server on 0.0.0.0:8080
}

func SetUp(maxBurstSize int) gin.HandlerFunc {
	return func(c *gin.Context) {
		val := c.Request.Header.Get("Token")
		var limiter *rate.Limiter
		if val != "" {
			var ok bool
			limiter, ok = m[val]
			if !ok {
				limiter = rate.NewLimiter(rate.Every(time.Second*1), maxBurstSize)
				m[val] = limiter
			}
		} else {
			c.JSON(http.StatusOK, gin.H{"msg": "no token"})
			c.Abort() // 打断调用链
			return
		}

		if limiter.Allow() {
			c.Next()
			return
		}
		fmt.Println("Too many requests")
		c.Writer.WriteString("Too many requests")
		c.Abort()
		return
	}
}
