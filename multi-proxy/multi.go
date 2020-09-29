// 多级代理

package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"net/url"
	"runtime/debug"
	"strings"
	"sync"
	"time"

	"github.com/robfig/cron"
)

var lock *sync.Mutex = &sync.Mutex{}
var proxyUrls map[string]bool = make(map[string]bool)
var mu sync.Mutex

var connHold map[string]net.Conn = make(map[string]net.Conn) //map[代理服务器url]tcp连接

var bufferPool = sync.Pool{
	New: func() interface{} { return make([]byte, 32*1024) },
}

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	refreshProxyAddr()

	cronTask := cron.New()
	cronTask.AddFunc("@every 1h", func() {
		mu.Lock()
		defer mu.Unlock()
		refreshProxyAddr()
	})
	cronTask.Start()
	//go func() {
	//	for {
	//		time.Sleep(3 * time.Second)
	//		for k, v := range proxyUrls {
	//			fmt.Println(k, v)
	//		}
	//	}

	//}()
}

func main() {
	l, err := net.Listen("tcp", ":9090")
	if err != nil {
		log.Panic(err)
	}

	for {
		client, err := l.Accept()
		if err != nil {
			log.Panic(err)
		}
		go handle(client)
	}
}

func handle(client net.Conn) {
	defer func() {
		if err := recover(); err != nil {
			log.Println(err)
			debug.PrintStack()
		}
	}()
	if client == nil {
		return
	}
	// log.Println("client tcp tunnel connection:", client.LocalAddr().String(), "->", client.RemoteAddr().String())
	// client.SetDeadline(time.Now().Add(time.Duration(10) * time.Second))
	defer client.Close()

	// var b [1024]byte
	b := bufferPool.Get().([]byte)
	defer bufferPool.Put(b)

	n, err := client.Read(b[:]) //读取应用层的所有数据
	index := bytes.IndexByte(b[:], '\n')
	if err != nil || index == -1 {
		log.Println(err, string(b)) // 传输层的连接是没有应用层的内容 比如：net.Dial()
		return
	}
	var method, host, address string
	fmt.Sscanf(string(b[:index]), "%s%s", &method, &host)
	if method == "" || host == "" {
		return
	}
	// log.Println(method, host)
	host = strings.TrimSpace(host)
	//if !strings.HasPrefix(host, "/") {
	//	return
	//}
	hostPortURL, err := url.Parse(host)
	if err != nil {
		log.Println(err, host)
		return
	}

	if hostPortURL.Opaque == "443" { //https访问
		address = hostPortURL.Scheme + ":443"
	} else { //http访问
		if strings.Index(hostPortURL.Host, ":") == -1 { //host不带端口， 默认80
			address = hostPortURL.Host + ":80"
		} else {
			address = hostPortURL.Host
		}
	}

	server, err := Dial("tcp", address)
	if err != nil {
		log.Println(err)
		return
	}
	//在应用层完成数据转发后，关闭传输层的通道
	defer server.Close()
	// log.Println("server tcp tunnel connection:", server.LocalAddr().String(), "->", server.RemoteAddr().String())
	// server.SetDeadline(time.Now().Add(time.Duration(10) * time.Second))

	method = strings.ToUpper(method)
	if method == "CONNECT" {
		fmt.Fprint(client, "HTTP/1.1 200 Connection established\r\n\r\n")
	} else {
		log.Println("server write", method) //其它协议
		server.Write(b[:n])
	}

	bufferPool.Put(b)

	//进行转发
	go func() {
		// io.Copy(server, client)
		Copy(server, client)
	}()
	// io.Copy(client, server) //阻塞转发
	Copy(client, server) //阻塞转发
}

func Copy(dst io.Writer, src io.Reader) (written int64, err error) {
	if wt, ok := src.(io.WriterTo); ok {
		return wt.WriteTo(dst)
	}
	if rt, ok := dst.(io.ReaderFrom); ok {
		return rt.ReadFrom(src)
	}

	buf := bufferPool.Get().([]byte)
	defer bufferPool.Put(buf)
	for {
		nr, er := src.Read(buf)
		if nr > 0 {
			fmt.Println(string(buf[:nr]))
			nw, ew := dst.Write(buf[0:nr])
			if nw > 0 {
				written += int64(nw)
			}
			if ew != nil {
				err = ew
				break
			}
			if nr != nw {
				err = io.ErrShortWrite
				break
			}
		}
		if er == io.EOF {
			break
		}
		if er != nil {
			err = er
			break
		}
	}
	return written, err
}

//refreshProxyAddr 刷新代理ip
func refreshProxyAddr() {
	//var proxyUrlsTmp map[string]string = make(map[string]string) // 获取代理ip地址逻辑
	//proxyUrls = proxyUrlsTmp                                     //可以手动设置测试代理ip
}

func getProxy(country string) ([]string, error) {
	// 'http://tiqu.linksocket.com:81/abroad?num=100&type=1&pro=0&city=0&yys=0&port=1&flow=1&ts=0&ys=0&cs=0&lb=1&sb=0&pb=4&mr=0&regions=my&n=0'
	proxyUrl := ConfigMap["proxy_source"]["my"]

	req, err := http.NewRequest("GET", proxyUrl, nil)
	if err != nil {
		return nil, err
	}
	tr := http.DefaultTransport
	res, err := tr.RoundTrip(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("server returned: %v", res.Status)
	}
	//b := bufferPool.Get().(*bytes.Buffer)
	//b.Reset()
	//defer bufferPool.Put(b)
	//_, err = io.Copy(b, res.Body)

	//if err != nil {
	//	return nil, fmt.Errorf("reading response body: %v", err)
	//}

	s := bufio.NewScanner(res.Body)
	s.Split(bufio.ScanLines)
	for s.Scan() {
		if -1 == strings.Index(s.Text(), "false") {
			lock.Lock()
			proxyUrls[s.Text()] = true
			lock.Unlock()
		} else {
			return nil, fmt.Errorf("reading response body: %v", err)
		}
	}

	// 获得的body
	/*
		156.251.125.135:3871
		156.251.125.135:3858
		156.251.125.145:5131
		156.251.125.145:5130
		156.251.125.142:5136
		156.251.125.142:5128
		156.251.125.145:5132
		156.251.125.145:5151
		156.251.125.142:5118
		156.251.125.135:3849
		156.251.125.145:5149
		156.251.125.135:3856
		156.251.125.135:3855
		156.251.125.135:3868
		156.251.125.145:5127
		156.251.125.145:5120
		156.251.125.135:3878
		156.251.125.145:5134
		156.251.125.135:3861
		156.251.125.135:3880
		156.251.125.145:5125
		156.251.125.142:5125
		156.251.125.145:5141
		156.251.125.135:3851
	*/

	return nil, nil
}

//DialSimple 直接通过发送数据报与二级代理服务器建立连接
//func DialSimple(network, addr string) (net.Conn, error) {
//	var proxyAddr string
//	for proxyAddr = range proxyUrls { //随机获取一个代理地址
//		break
//	}
//	c, err := func() (net.Conn, error) {
//		u, _ := url.Parse(proxyAddr)
//		log.Println("代理host", u.Host)
//		// Dial and create client connection.
//		c, err := net.DialTimeout("tcp", u.Host, time.Second*5)
//		if err != nil {
//			log.Println(err)
//			return nil, err
//		}
//		_, err = c.Write([]byte("CONNECT w.xxxx.com:443 HTTP/1.1\r\n Host: w.xxxx.com:443\r\n User-Agent: Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/79.0.3945.88 Safari/537.3\r\n\r\n")) // w.xxxx.com:443 替换成实际的地址
//		if err != nil {
//			panic(err)
//		}
//		c.Write([]byte(`GET www.baidu.com HTTP/1.1\r\n\r\n`))
//		io.Copy(os.Stdout, c)
//		return c, err
//	}()
//	return c, err
//}

//Dial 建立一个传输通道
func Dial(network, addr string) (net.Conn, error) {
	var proxyAddr string
	var ok bool
	lock.Lock()
	for proxyAddr, ok = range proxyUrls { //随机获取一个代理地址
		if ok {
			break
		}
	}
	lock.Unlock()

	if proxyAddr == "" || !ok {
		_, err := getProxy("cn")
		if err != nil {
			proxyAddr = addr
		} else {
			lock.Lock()
			for proxyAddr, ok = range proxyUrls { //随机获取一个代理地址
				if ok {
					break
				}
			}
			lock.Unlock()
		}
	}
	//建立到代理服务器的传输层通道
	c, err := func() (net.Conn, error) {
		//u, err := url.Parse(proxyAddr)
		//if err != nil {
		//	panic(err)
		//}
		h := proxyAddr
		// log.Println("代理地址", h)
		// Dial and create client connection.
		// conn, err := net.DialTimeout("tcp", u.Host, time.Second*5)
		conn, err := net.DialTimeout("tcp", h, time.Second*5)
		if err != nil {
			return nil, err
		}

		reqURL, err := url.Parse("http://" + addr)
		if err != nil {
			return nil, err
		}
		req, err := http.NewRequest(http.MethodConnect, reqURL.String(), nil)
		if err != nil {
			return nil, err
		}
		req.Close = false
		req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/79.0.3945.88 Safari/537.3")

		err = req.Write(conn)
		if err != nil {
			return nil, err
		}

		resp, err := http.ReadResponse(bufio.NewReader(conn), req)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		// log.Println(resp.StatusCode, resp.Status, resp.Proto, resp.Header)
		if resp.StatusCode != 200 {
			err = fmt.Errorf("Connect server using proxy error, StatusCode [%d]", resp.StatusCode)
			return nil, err
		}
		return conn, err
	}()
	if c == nil || err != nil { //代理异常
		log.Println("代理异常：", proxyAddr, err)
		lock.Lock()
		proxyUrls[proxyAddr] = false
		lock.Unlock()

		return net.Dial(network, addr)
		// return nil, fmt.Errorf("server returned: %v", "xxxx")
	}
	// log.Println("代理正常,tunnel信息", c.LocalAddr().String(), "->", c.RemoteAddr().String())
	return c, err
}
