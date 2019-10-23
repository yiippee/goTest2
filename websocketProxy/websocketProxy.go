package main

import (
	"errors"
	"fmt"
	"github.com/fatih/pool"
	"github.com/gorilla/websocket"
	"goTest/websocketProxy/discovery"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync"
)

const ETCD_ADDR = "http://127.0.0.1:2379" // etcd 服务地址
const SERVICES_DIR = "services/"          // 监听的目录
const LOCALHOST = ":6699"                 // 本地地址

var bufpool *sync.Pool
var connPool pool.Pool

var upgrader = websocket.Upgrader{
	ReadBufferSize:    4096,
	WriteBufferSize:   4096,
	EnableCompression: true,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type handler struct {
	master *discovery.Master
}

func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query()["token"]
	if token == nil {
		// 去鉴权
		fmt.Fprintf(w, "%s", "no auth,please login.")
		return
	}
	// 根据token路由
	//h.master.Lock()
	//l := len(h.master.Keys)
	//if l == 0 {
	//	// 没有机器啊
	//	fmt.Fprintf(w, "%s", "no service.")
	//	h.master.Unlock()
	//	return
	//}
	//key := h.master.Keys[int(crc32.ChecksumIEEE([]byte(token[0])))%l]
	//ip := h.master.Nodes[key].Info.IP
	//h.master.Unlock()
	//fmt.Println(ip)
	// 为什么一定要改成 http 呢？ 因为在握手阶段使用的就是http啊
	// 在具体的通信阶段才使用ws协议通信。这里还是需要http进行的
	// 具体通信过程，代理服务器自动做了。只要你在header里面添加upgrade字段就可以标识是ws了
	ip := "http://127.0.0.1:8080" // "ws://127.0.0.1:8080"

	backURL, _ := url.Parse(ip)
	rproxy := httputil.NewSingleHostReverseProxy(backURL)

	rproxy.ServeHTTP(w, r)
	// 所以就不需要以下那么多代码去实现转发了，不过下面的实现应该是可以的。只是现在golang官方标准库有实现
	//u := url.URL{Scheme: "ws", Host: ip, Path: r.URL.Path}
	//dst, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//
	//src, err := upgrader.Upgrade(w, r, nil)
	//
	//exitChan := make(chan struct{}, 1)
	//go func() {
	//	defer func() {
	//		exitChan <- struct{}{}
	//	}()
	//
	//	for {
	//		mt, r, err := src.NextReader()
	//		if err != nil {
	//			if err != io.EOF {
	//				log.Println("NextReader:", err)
	//			}
	//			return
	//		}
	//		if mt == websocket.TextMessage {
	//			r = &validator{r: r}
	//		}
	//		w, err := dst.NextWriter(mt)
	//		if err != nil {
	//			log.Println("NextWriter:", err)
	//			return
	//		}
	//		if mt == websocket.TextMessage {
	//			r = &validator{r: r}
	//		}
	//		_, err = io.Copy(w, r)
	//
	//		if err != nil {
	//			if err == errInvalidUTF8 {
	//				_ = src.WriteControl(websocket.CloseMessage,
	//					websocket.FormatCloseMessage(websocket.CloseInvalidFramePayloadData, ""),
	//					time.Time{})
	//			}
	//			log.Println("Copy:", err)
	//			return
	//		}
	//		err = w.Close()
	//		if err != nil {
	//			log.Println("Close:", err)
	//			return
	//		}
	//	}
	//}()
	//
	//go func() {
	//	for {
	//		defer func() {
	//			exitChan <- struct{}{}
	//		}()
	//		mt, r, err := dst.NextReader()
	//		if err != nil {
	//			if err != io.EOF {
	//				log.Println("NextReader:", err)
	//			}
	//			return
	//		}
	//		if mt == websocket.TextMessage {
	//			r = &validator{r: r}
	//		}
	//		w, err := src.NextWriter(mt)
	//		if err != nil {
	//			log.Println("NextWriter:", err)
	//			return
	//		}
	//		if mt == websocket.TextMessage {
	//			r = &validator{r: r}
	//		}
	//		_, err = io.Copy(w, r)
	//
	//		if err != nil {
	//			if err == errInvalidUTF8 {
	//				_ = src.WriteControl(websocket.CloseMessage,
	//					websocket.FormatCloseMessage(websocket.CloseInvalidFramePayloadData, ""),
	//					time.Time{})
	//			}
	//			log.Println("Copy:", err)
	//			return
	//		}
	//		err = w.Close()
	//		if err != nil {
	//			log.Println("Close:", err)
	//			return
	//		}
	//	}
	//}()
	//
	//<- exitChan
}

func main() {
	// 开启服务监听器，监听 services/ 目录下所有的服务器
	master, err := discovery.NewMaster([]string{ETCD_ADDR}, SERVICES_DIR)
	if err != nil {
		log.Fatal(err)
	}

	_ = http.ListenAndServe(LOCALHOST, &handler{master})
}

func init() {
	// 各种初始化动作，比如服务发现获取可提供服务的后端信息，这里写死来啊，作为例子展示
	bufpool = &sync.Pool{}
	bufpool.New = func() interface{} {
		return make([]byte, 32*1024)
	}
}

type validator struct {
	state int
	x     rune
	r     io.Reader
}

var errInvalidUTF8 = errors.New("invalid utf8")

func (r *validator) Read(p []byte) (int, error) {
	n, err := r.r.Read(p)
	state := r.state
	x := r.x
	for _, b := range p[:n] {
		state, x = decode(state, x, b)
		if state == utf8Reject {
			break
		}
	}
	r.state = state
	r.x = x
	if state == utf8Reject || (err == io.EOF && state != utf8Accept) {
		return n, errInvalidUTF8
	}
	return n, err
}

// UTF-8 decoder from http://bjoern.hoehrmann.de/utf-8/decoder/dfa/
//
// Copyright (c) 2008-2009 Bjoern Hoehrmann <bjoern@hoehrmann.de>
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to
// deal in the Software without restriction, including without limitation the
// rights to use, copy, modify, merge, publish, distribute, sublicense, and/or
// sell copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING
// FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS
// IN THE SOFTWARE.
var utf8d = [...]byte{
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, // 00..1f
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, // 20..3f
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, // 40..5f
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, // 60..7f
	1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, // 80..9f
	7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, // a0..bf
	8, 8, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, // c0..df
	0xa, 0x3, 0x3, 0x3, 0x3, 0x3, 0x3, 0x3, 0x3, 0x3, 0x3, 0x3, 0x3, 0x4, 0x3, 0x3, // e0..ef
	0xb, 0x6, 0x6, 0x6, 0x5, 0x8, 0x8, 0x8, 0x8, 0x8, 0x8, 0x8, 0x8, 0x8, 0x8, 0x8, // f0..ff
	0x0, 0x1, 0x2, 0x3, 0x5, 0x8, 0x7, 0x1, 0x1, 0x1, 0x4, 0x6, 0x1, 0x1, 0x1, 0x1, // s0..s0
	1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 0, 1, 1, 1, 1, 1, 0, 1, 0, 1, 1, 1, 1, 1, 1, // s1..s2
	1, 2, 1, 1, 1, 1, 1, 2, 1, 2, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 2, 1, 1, 1, 1, 1, 1, 1, 1, // s3..s4
	1, 2, 1, 1, 1, 1, 1, 1, 1, 2, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 3, 1, 3, 1, 1, 1, 1, 1, 1, // s5..s6
	1, 3, 1, 1, 1, 1, 1, 3, 1, 3, 1, 1, 1, 1, 1, 1, 1, 3, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, // s7..s8
}

const (
	utf8Accept = 0
	utf8Reject = 1
)

func decode(state int, x rune, b byte) (int, rune) {
	t := utf8d[b]
	if state != utf8Accept {
		x = rune(b&0x3f) | (x << 6)
	} else {
		x = rune((0xff >> t) & b)
	}
	state = int(utf8d[256+state*16+int(t)])
	return state, x
}
