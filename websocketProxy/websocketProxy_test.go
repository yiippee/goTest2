package main

import (
	"bufio"
	"fmt"
	"golang.org/x/net/http/httpguts"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"net/http/httputil"
	"net/url"
	"strings"
	"testing"
)

func upgradeType(h http.Header) string {
	if !httpguts.HeaderValuesContainsToken(h["Connection"], "Upgrade") {
		return ""
	}
	return strings.ToLower(h.Get("Upgrade"))
}

func TestWebsocketProxy(t *testing.T) {
	backendServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if upgradeType(r.Header) != "websocket" {
			t.Error("unexpected backend request")
			http.Error(w, "unexpected request", 400)
			return
		}
		c, _, err := w.(http.Hijacker).Hijack()
		if err != nil {
			t.Error(err)
			return
		}
		defer c.Close()
		io.WriteString(c, "HTTP/1.1 101 Switching Protocols\r\nConnection: upgrade\r\nUpgrade: WebSocket\r\n\r\n")
		bs := bufio.NewScanner(c)
		if !bs.Scan() {
			t.Errorf("backend failed to read line from client: %v", bs.Err())
			return
		}
		fmt.Fprintf(c, "backend got %q\n", bs.Text())
	}))
	defer backendServer.Close()

	backURL, _ := url.Parse(backendServer.URL)
	rproxy := httputil.NewSingleHostReverseProxy(backURL)
	rproxy.ErrorLog = log.New(ioutil.Discard, "", 0) // quiet for tests
	rproxy.ModifyResponse = func(res *http.Response) error {
		res.Header.Add("X-Modified", "true")
		return nil
	}

	handler := http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.Header().Set("X-Header", "X-Value")
		rproxy.ServeHTTP(rw, req)
	})

	frontendProxy := httptest.NewServer(handler)
	defer frontendProxy.Close()

	req, _ := http.NewRequest("GET", frontendProxy.URL, nil)
	req.Header.Set("Connection", "Upgrade")
	req.Header.Set("Upgrade", "websocket")

	c := frontendProxy.Client()
	res, err := c.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	if res.StatusCode != 101 {
		t.Fatalf("status = %v; want 101", res.Status)
	}

	got := res.Header.Get("X-Header")
	want := "X-Value"
	if got != want {
		t.Errorf("Header(XHeader) = %q; want %q", got, want)
	}

	if upgradeType(res.Header) != "websocket" {
		t.Fatalf("not websocket upgrade; got %#v", res.Header)
	}
	rwc, ok := res.Body.(io.ReadWriteCloser)
	if !ok {
		t.Fatalf("response body is of type %T; does not implement ReadWriteCloser", res.Body)
	}
	defer rwc.Close()

	if got, want := res.Header.Get("X-Modified"), "true"; got != want {
		t.Errorf("response X-Modified header = %q; want %q", got, want)
	}

	io.WriteString(rwc, "Hello\n")
	bs := bufio.NewScanner(rwc)
	if !bs.Scan() {
		t.Fatalf("Scan: %v", bs.Err())
	}
	got = bs.Text()
	want = `backend got "Hello"`
	if got != want {
		t.Errorf("got %#q, want %#q", got, want)
	}
}
