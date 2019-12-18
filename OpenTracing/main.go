package main

import (
	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"log"
	"math/rand"
	"net/http"
)

func MyMiddleware(c *gin.Context) {
	spanCtx, _ := opentracing.GlobalTracer().Extract(opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(r.Header))
	sp := opentracing.GlobalTracer().StartSpan(c.Request.URL.Path, opentracing.ChildOf(spanCtx))
	defer sp.Finish()

	if err := opentracing.GlobalTracer().Inject(
		sp.Context(),
		opentracing.HTTPHeaders,
		opentracing.HTTPHeadersCarrier(c.Request.Header)); err != nil {
		log.Println(err)
	}

	sct := &status_code.StatusCodeTracker{ResponseWriter: w, Status: http.StatusOK}
	h.ServeHTTP(sct.WrappedResponseWriter(), r)

	ext.HTTPMethod.Set(sp, c.Request.Method)
	ext.HTTPUrl.Set(sp, c.Request.URL.EscapedPath())
	ext.HTTPStatusCode.Set(sp, uint16(sct.Status))
	if sct.Status >= http.StatusInternalServerError {
		ext.Error.Set(sp, true)
	} else if rand.Intn(100) > sf {
		ext.SamplingPriority.Set(sp, 0)
	}

	c.Next()
}

func main() {
	router := gin.Default()
	router.Use(MyMiddleware)

	router.GET("ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong"})
	})

	router.Run(":8085")
}
