package main

import (
	"fmt"
	"sync"
	"testing"
)

func doSomething() {
	myFind()
}

func doSomethingPrepare(n int) {

}

func BenchmarkDemo(b *testing.B) {
	doSomethingPrepare(b.N)
	b.ResetTimer() // 重新计时，那准备工作不计入测试时间

	fmt.Println(b.N)
	var wg sync.WaitGroup

	for i := 0; i != b.N; i++ {
		wg.Add(1)
		go func() {
			myFind()
			wg.Done()
		}()
	}
	wg.Wait()
}