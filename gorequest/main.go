package main

import (
	"fmt"

	"github.com/parnurzeal/gorequest"
)

func main() {
	request := gorequest.New().Proxy("http://156.251.125.135:3976")
	resp, body, errs := request.Get("http://httpbin.org/get").End()
	//// To reuse same client with no_proxy, use empty string:
	//resp, body, errs = request.Proxy("").Get("http://example-no-proxy.com").End()
	fmt.Println(resp, body, errs)
}
