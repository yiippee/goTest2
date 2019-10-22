package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/fatih/pool"
	"hash/crc32"
	"io"
	"net"
	"os"
	"strings"
	"sync"
)

var lock sync.Mutex
var IPList []string
var ip string
var list string
var bufpool *sync.Pool
var connPool pool.Pool

func main() {
	flag.StringVar(&ip, "l", ":9897", "-l=0.0.0.0:9897 指定服务监听的端口")
	flag.StringVar(&list, "d", "127.0.0.1:1789", "-d=127.0.0.1:1789,127.0.0.1:1788 指定后端的IP和端口,多个用','隔开")
	flag.Parse()
	IPList = strings.Split(list, ",")
	if len(IPList) <= 0 {
		fmt.Println("后端IP和端口不能空,或者无效")
		os.Exit(1)
	}
	server()
}

func server() {
	lis, err := net.Listen("tcp", ip)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer lis.Close()
	for {
		conn, err := lis.Accept()
		if err != nil {
			fmt.Println("建立连接错误:%v\n", err)
			continue
		}
		fmt.Println(conn.RemoteAddr(), conn.LocalAddr())
		go handle(conn)
	}
}

func handle(sconn net.Conn) {
	defer sconn.Close()

	// 选择需要连接的服务器ip,这里可以做一些路由策略
	// todo 这里考虑根据路由策略，获取哪一个连接池来获取具体的连接。具体的策略可以考虑一致性hash算法，待验证
	ip, ok := getIP(sconn)
	if !ok {
		return
	}

	// 此处可以设置一个连接池来优化
	// dconn, err := net.Dial("tcp", ip)
	dconn, err := connPool.Get()
	if err != nil {
		fmt.Printf("连接%v失败:%v\n", ip, err)
		return
	}
	ExitChan := make(chan bool, 1)

	//go func(sconn net.Conn, dconn net.Conn, Exit chan bool) {
	//	// io.Copy 性能不佳，
	//	_, err := Copy(dconn, sconn)
	//	fmt.Printf("往%v发送数据失败:%v\n", ip, err)
	//	ExitChan <- true
	//}(sconn, dconn, ExitChan)
	//
	//go func(sconn net.Conn, dconn net.Conn, Exit chan bool) {
	//	_, err := Copy(sconn, dconn)
	//
	//	fmt.Printf("从%v接收数据失败:%v\n", ip, err)
	//	ExitChan <- true
	//}(sconn, dconn, ExitChan)

	// 不需要传参了，因为这些变量只属于这一个goroutine了，可以直接捕获
	go func() {
		// io.Copy 性能不佳，所以自己实现copy，主要用到了buf池
		_, err := Copy(dconn, sconn)
		fmt.Printf("往%v发送数据失败:%v\n", ip, err)
		ExitChan <- true
	}()

	go func() {
		_, err := Copy(sconn, dconn)

		fmt.Printf("从%v接收数据失败:%v\n", ip, err)
		ExitChan <- true
	}()

	<-ExitChan

	dconn.Close()
}

func getIP(sconn net.Conn) (string, bool) {

	// 读取登录信息，可以读一行。因为对于tcp协议来说，基于字节流的，所以需要自己规定协议。而使用websocket则不需要
	reader := bufio.NewReader(sconn)
	buf, _, err := reader.ReadLine()

	if err != nil {
		return "", false
	}

	lock.Lock()
	defer lock.Unlock()

	if len(IPList) < 1 {
		return "", false
	}
	// 根据登录者的id，hash到具体的某一个后端服务器、
	// 对于有状态的服务器（比如数据存储服务器），则更适合采用一致性hash算法
	ip := IPList[crc32.ChecksumIEEE(buf)%1]
	ip = IPList[crc32.ChecksumIEEE([]byte("lizhanbin"))%1]
	ip = IPList[crc32.ChecksumIEEE([]byte("zhangsumin"))%1]

	return ip, true
}

func Copy(dst io.Writer, src io.Reader) (written int64, err error) {
	if wt, ok := src.(io.WriterTo); ok {
		return wt.WriteTo(dst)
	}
	if rt, ok := dst.(io.ReaderFrom); ok {
		return rt.ReadFrom(src)
	}

	buf := bufpool.Get().([]byte)
	defer bufpool.Put(buf)
	for {
		nr, er := src.Read(buf)
		if nr > 0 {
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

func init() {
	// 内存池
	bufpool = &sync.Pool{}
	bufpool.New = func() interface{} {
		return make([]byte, 32*1024)
	}

	// 连接池
	// create a factory() to be used with channel based pool
	// 可以根据多个后端服务器，创建多个连接池，然后通过一致性hash算法来选择某一个线程池
	// 一致性hash算法是为了找到具体某一台服务器的，本质上是一个找服务器的过程。
	// 如果服务器是有状态的，则很有效，因为可以明确定位到某一台服务器，虽然采用普通的hash算法也可以定位，但是
	// 一致性hash算法在服务器节点变动时，可以达到 1 / (n+1) 的修改操作。
	// 再结合虚拟节点，可以达到平衡性，虚拟节点的存在只是多了一个虚拟节点到实体节点的映射过程。
	factory := func() (net.Conn, error) { return net.Dial("tcp", "127.0.0.1:1789") }

	// create a new channel based pool with an initial capacity of 5 and maximum
	// capacity of 30. The factory will create 5 initial connections and put it
	// into the pool.
	var err error
	connPool, err = pool.NewChannelPool(5, 30, factory)
	if err != nil {
		panic(err)
	}
}
