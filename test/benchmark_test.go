package main

import (
	"fmt"
	"testing"
	"time"
)

func doSomething() {

}

func doSomethingPrepare(n int) {
	for i := 0; i != n; i++ {
		_ = time.Now()
	}
}

func BenchmarkDemo(b *testing.B) {
	doSomethingPrepare(b.N)
	b.ResetTimer() // 重新计时，那准备工作不计入测试时间

	fmt.Println(b.N)
	for i := 0; i != b.N; i++ {
		doSomething()
	}
}