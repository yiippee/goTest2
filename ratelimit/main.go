package main

import (
	"github.com/didip/tollbooth"
	// "golang.org/x/time/rate"
	"net/http"
)

func HelloHandler(w http.ResponseWriter, req *http.Request) {
	w.Write([]byte("Hello, World!"))
}

func main() {
	// 第一个参数是r Limit。代表每秒可以向Token桶中产生多少token。Limit实际上是float64的别名。
	// 第二个参数是b int。b代表Token桶的容量大小。
	// limiter := rate.NewLimiter(10, 1) // 对于以上例子来说，其构造出的限流器含义为，其令牌桶大小为1, 以每秒10个Token的速率向桶中放置Token。

	// Create a request limiter per handler.
	http.Handle("/", tollbooth.LimitFuncHandler(tollbooth.NewLimiter(1, nil), HelloHandler))
	http.ListenAndServe(":12345", nil)
}
