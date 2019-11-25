package main

import (
	"fmt"
	"github.com/gorilla/websocket"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:    4096,
	WriteBufferSize:   4096,
	EnableCompression: true,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func ChatServer(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade:", err)
		return
	}
	defer conn.Close()

	//var test []byte
	//var payload []byte
	for {
		//t, b, err := conn.ReadMessage()
		//if err != nil {
		//	log.Println(err)
		//}

		_, r, err := conn.NextReader()
		if err != nil {
			if err != io.EOF {
				log.Println("NextReader:", err)
			}
			return
		}

		//fmt.Println("Payload: ", len(b), t)
		//
		//test = append(test, b...)
		//fmt.Println("Test: ", len(test))

		fo, err := os.Create(fmt.Sprintf("E:\\lzb\\golang\\src\\goTest2\\websocket\\image\\./%d.png", time.Now().UnixNano()))
		io.Copy(fo, r)
		//check(err)
		//_, err = fo.Write(test)
		//check(err)
		fo.Close()
	}
	log.Print("DONE")

}

func main() {
	fmt.Println("Starting... ")

	http.HandleFunc("/ws", ChatServer)
	go client()
	err := http.ListenAndServe(":3000", nil)

	if err != nil {
		log.Fatal("ListenAndServe ", err)
	}
}

func client() {
	time.Sleep(3 * time.Second)
	u := url.URL{Scheme: "ws", Host: "127.0.0.1:3000", Path: "/ws"}
	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal(err)
	}

	file, err := os.Open("E:\\lzb\\golang\\src\\goTest2\\websocket\\image\\123.png")
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()
	buff, err := ioutil.ReadAll(file)

	err = conn.WriteMessage(2, buff)
	if err != nil {
		log.Println("WriteMessage:", err)
	}
}
