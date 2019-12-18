package main

import (
	"errors"
	//"github.com/didip/tollbooth"
	"github.com/gin-gonic/gin"
	limiter "github.com/julianshen/gin-limiter"
	"time"

	// "golang.org/x/time/rate"
	"net/http"
)

func HelloHandler(w http.ResponseWriter, req *http.Request) {
	w.Write([]byte("Hello, World!"))
}

func JWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next() // 跳过jwt验证
	}
}
func main() {
	// 第一个参数是r Limit。代表每秒可以向Token桶中产生多少token。Limit实际上是float64的别名。
	// 第二个参数是b int。b代表Token桶的容量大小。
	// _ := rate.NewLimiter(10, 1) // 对于以上例子来说，其构造出的限流器含义为，其令牌桶大小为1, 以每秒10个Token的速率向桶中放置Token。

	// Create a request limiter per handler.
	r := gin.Default()
	// r.GET("/", tollbooth.LimitFuncHandler(tollbooth.NewLimiter(0.5, nil), HelloHandler))
	//Allow only 10 requests per minute per API-Key
	lm := limiter.NewRateLimiter(time.Second, 1, func(ctx *gin.Context) (string, error) {
		key := ctx.Request.Header.Get("X-API-KEY")
		if key != "" {
			return key, nil
		}
		return "", errors.New("API key is missing")
	})

	//Apply only to /ping
	r.GET("/ping", lm.Middleware(), func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	//Allow only 5 requests per second per user
	lm2 := limiter.NewRateLimiter(time.Second, 1, func(ctx *gin.Context) (string, error) {
		key := ctx.Request.Header.Get("X-USER-TOKEN")
		if key != "" {
			return key, nil
		}
		return "", errors.New("User is not authorized")
	})

	//Apply to a group
	x := r.Group("/v2")
	x.Use(lm2.Middleware())
	{
		x.GET("/ping", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"message": "pong",
			})
		})
		x.GET("/another_ping", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"message": "pong pong",
			})
		})
	}

	r.Run(":8087") // listen and server on 0.0.0.0:8080
}
