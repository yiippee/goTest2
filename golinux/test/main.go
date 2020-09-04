package main

import "fmt"

func test(a int) func(i int) int {
	return func(i int) int {
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

/*
# 下面的指令标明把 main.go 生成 linux 下的 amd64 二进制文件
# 其中 -N 指定编译器不要进行优化，-l 指定编译器不要对函数进行内联处理
# 其中 -o testl 指定输出二进制文件到 testl 中
# -gcflags 的参数可以通过 go tool compile --help 获取
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build --gcflags "-N -l" -o testl main.go

# 可以通过 go tool objdump --help 来查看 objdump 的 -s 用法
# 比如 go tool objdump -s "^main.main$" testl 只返回 main.main 函数的汇编代码
# 下面的指令标明把 上一步生成的 testl 提取汇编代码到 ojbl.S 文件中
go tool objdump -S testl > objl.S

著作权归作者所有。
商业转载请联系作者获得授权,非商业转载请注明出处。
原文: https://jingwei.link/2019/06/01/golang-outer-variable-in-clousure.html


go tool compile -S .\test.go
go tool objdump .\test.o

*/
