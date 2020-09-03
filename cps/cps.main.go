package main

import "fmt"

// CPS，是Continuation Passing Style的缩写，它是一种编码风格，函数执行完以后，并不通过返回值，而是调用它自己的Continuation来完成计算。

// 在一般的写法中我们定义了一个相加函数然后向控制台打印了a+b的结果.
func add(x, y int) int {
	z := x + y
	fmt.Println(z)

	return x + y
}

func _k(x int) int {
	return x
}

type f func(x int) int

func fib_cps(n int, k f) int {
	if n < 2 {
		return k(n)
	}

	return fib_cps(n-1, func(x int) int {
		return fib_cps(n-2, func(y int) int {
			return k(x + y)
		})
	})
}

// 递归实现求阶乘
func factorial(n int) int {
	if n <= 2 {
		return n
	}
	return n * factorial(n-1)
}

// cps 风格
// func factorial_cps()

func main() {
	ret := factorial(5)
	fmt.Println("ret: ", ret)

	//ret := fib_cps(10, _k)
	//fmt.Println("ret: ", ret)
}
