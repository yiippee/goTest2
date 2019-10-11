package main

import (
"github.com/rakyll/statik/fs"
	"log"
	"net/http"

	_ "goTest/staticFile/statik" // TODO: Replace with the absolute import path
)

// 将资源文件全部打包进go文件中，打包成一个单一的可执行二进制，简单
func main() {
	statikFS, err := fs.New()
	if err != nil {
		log.Fatal(err)
	}

	http.Handle("/public/", http.StripPrefix("/public/", http.FileServer(statikFS)))
	_ = http.ListenAndServe(":8080", nil)
}
