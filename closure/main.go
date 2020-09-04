package main

import "fmt"

func test(a int) func(i int) int {
	// 类是有行为的数据，而闭包是有数据的行为；
	return func(i int) int {
		// a必须在堆上分配，如果在栈上分配，函数结束，a也被回收了；然后会定义出一个匿名结构体:
		// type.struct {
		//	 F uintptr // 这个就是闭包调用的函数指针
		// 	 a *int    // 这就是闭包的上下文数据
		// }
		/*
				接着生成一个该对象，并将之前在堆上分配的整型对象a的地址赋值给结构体中的a指针，
			接下来将闭包调用的func函数地址赋值给结构体中F指针；这样，每生成一个闭包函数，
			其实就是生成一个上述结构体对象，每个闭包对象也就有自己的数据a和调用函数F；
			最后将这个结构体的地址返回给main函数；
		*/
		a = a + i
		return a
	}
}
func main() {
	f := test(1)
	a := f(2)
	fmt.Println(a)
	b := f(3)
	fmt.Println(b)
}
