package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

func main() {
	var s sliceMap
	s.Add([]byte("name"), []byte("lzb"))
	s.Add([]byte("name"), []byte("zsm"))

	mux := http.NewServeMux()
	mux.HandleFunc("/test/111", Test)
	mux.HandleFunc("/test2/", Test2)
	mux.HandleFunc("/image", func(w http.ResponseWriter, r *http.Request) {
		file, err := os.Open("E:\\lzb\\golang\\src\\goTest2\\http\\123.png")
		if err != nil {
			fmt.Println(err)
		}
		defer file.Close()
		buff, err := ioutil.ReadAll(file)
		w.Write(buff)
	})

	mux2 := http.NewServeMux()
	mux2.HandleFunc("/test/222", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "%s", "hello,world 222.")
	})
	mux2.HandleFunc("/test/223", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "%s", "hello,world 223.")
	})
	mux2.HandleFunc("/test2/", Test2)

	go http.ListenAndServe(":16379", mux)
	http.ListenAndServe(":26379", mux2)
}

func Test(w http.ResponseWriter, r *http.Request) {
	time.Sleep(2 * time.Millisecond)
	fmt.Println("test...")
	fmt.Fprintf(w, "%s", "hello,world.")
}

func Test2(w http.ResponseWriter, r *http.Request) {
	time.Sleep(2 * time.Millisecond)
	fmt.Println("test...")
	fmt.Fprintf(w, "%s", "test2")
}

// slice 复用
type kv struct {
	key   []byte
	value []byte
}

type sliceMap []kv

func (sm *sliceMap) Add(k, v []byte) {
	kvs := *sm
	if cap(kvs) > len(kvs) {
		kvs = kvs[:len(kvs)+1]
	} else {
		kvs = append(kvs, kv{})
	}

	kv := &kvs[len(kvs)-1]
	kv.key = append(kv.key[:0], k...)
	kv.value = append(kv.value[:0], v...)

	*sm = kvs
}
