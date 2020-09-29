// proxyServer.go

package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

//将request转发给 http://127.0.0.1:2003
func helloHandler(w http.ResponseWriter, r *http.Request) {

	trueServer := "http://127.0.0.1:9090"

	url, err := url.Parse(trueServer)
	if err != nil {
		log.Println(err)
		return
	}

	proxy := httputil.NewSingleHostReverseProxy(url)
	proxy.ServeHTTP(w, r)

}

func main() {
	http.HandleFunc("/", helloHandler)
	log.Fatal(http.ListenAndServe(":9080", nil))
}
