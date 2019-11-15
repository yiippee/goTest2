package main

import (
	"github.com/gin-gonic/gin"
	"github.com/unrolled/secure"
	"net/http"
)

func main() {
	router := gin.Default()
	router.Use(TlsHandler())
	router.GET("/test/222", func(c *gin.Context) {

		//time.Sleep(time.Millisecond * 1)

		c.String(http.StatusOK, "OK")
	})
	_ = router.RunTLS(":26379", "ssl.pem", "ssl.key")
}

func TlsHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		secureMiddleware := secure.New(secure.Options{
			SSLRedirect: true,
			SSLHost:     "localhost:8080",
		})
		err := secureMiddleware.Process(c.Writer, c.Request)

		// If there was an error, do not continue.
		if err != nil {
			return
		}

		c.Next()
	}
}