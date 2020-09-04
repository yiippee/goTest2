TEXT main.main(SB) /golang/src/jingwei.link/main.go
func main() {
  0x488300		64488b0c25f8ffffff	MOVQ FS:0xfffffff8, CX
  0x488309		483b6110		CMPQ 0x10(CX), SP
  0x48830d		0f8690000000		JBE 0x4883a3		; 上面三是对栈进行扩容判定，如果栈不够用了，会进行扩容
  0x488313		4883ec40		SUBQ $0x40, SP		; 预留出 0x40 的栈空间供 main 函数使用
  0x488317		48896c2438		MOVQ BP, 0x38(SP)
  0x48831c		488d6c2438		LEAQ 0x38(SP), BP	; 上面两句待探究，应该是为了保存某个场景为未来恢复某个状态做准备
	a1 := 1
  0x488321		48c744240801000000	MOVQ $0x1, 0x8(SP)	; 把 1 赋值到 0x8(SP) 的地址，即 a1
	a2 := 2
  0x48832a		48c7042402000000	MOVQ $0x2, 0(SP)	; 把 2 赋值到 0x8(SP) 的地址，即 a2
	myfunc := func() {
  0x488332		48c744242000000000	MOVQ $0x0, 0x20(SP)
  0x48833b		0f57c0			XORPS X0, X0
  0x48833e		0f11442428		MOVUPS X0, 0x28(SP)
  0x488343		488d542420		LEAQ 0x20(SP), DX		; 把 0x20(SP) 的地址加载到 DX 中
  0x488348		4889542418		MOVQ DX, 0x18(SP)		; 把 DX 的值，即 0x20(SP) 的值，赋值到 0x18(SP) 中; 0x18(SP) 中保存的是 0x20(SP) 的地址
  0x48834d		8402			TESTB AL, 0(DX)
  0x48834f		488d05ca000000		LEAQ main.main.func1(SB), AX	; 把 func1（我们定义的闭包函数体） 的地址赋值到 AX
  0x488356		4889442420		MOVQ AX, 0x20(SP)		; 把 AX 的值，即 func1 的地址，赋值到 0x20(SP) 中； 0x20(SP) 中保存的是 func1 的调用地址
  0x48835b		8402			TESTB AL, 0(DX)
  0x48835d		488d442408		LEAQ 0x8(SP), AX		; 把 0x8(SP) 的地址，即 a1 的地址（指针）赋值到 AX
  0x488362		4889442428		MOVQ AX, 0x28(SP)		; 把 a1 赋值到 0x28(SP) 中；0x28(SP) 中保存的是 a1 的地址
  0x488367		8402			TESTB AL, 0(DX)
  0x488369		488d0424		LEAQ 0(SP), AX			; 把 0(SP) 的地址，即 a2 的地址（指针）赋值到 AX
  0x48836d		4889442430		MOVQ AX, 0x30(SP)		; 把 a2 赋值到 0x30(SP) 中；0x30(SP) 中保存的是 a2 的地址
  0x488372		4889542410		MOVQ DX, 0x10(SP)		; 把 DX 的值，即 0x20(SP) 的地址，赋值到 0x10(SP) 中；0x10(SP) 中保存的是 0x20(SP) 的地址
	myfunc()
  0x488377		488b442420		MOVQ 0x20(SP), AX	; 把 0x20(SP)  中的内容，即 func1 的地址加载到 AX 寄存器
  0x48837c		ffd0			CALL AX			; 调用 func1 函数
	a1 = 11
  0x48837e		48c74424080b000000	MOVQ $0xb, 0x8(SP)	; 把 11 赋值到 0x8(SP) 的地址，即更新 a1
	a2 = 22
  0x488387		48c7042416000000	MOVQ $0x16, 0(SP)	; 把 22 赋值到 0(SP) 的地址，即更新 a2
	myfunc()
  0x48838f		488b542410		MOVQ 0x10(SP), DX	; 这里把 0x10(SP) 中的值，即 0x20(SP) 的地址加载到 DX 寄存器
  0x488394		488b02			MOVQ 0(DX), AX		; 把 0(DX) 中的值，即 func1 的地址加载到 AX 寄存器
  0x488397		ffd0			CALL AX			; 调用 func 1 函数。
}
  0x488399		488b6c2438		MOVQ 0x38(SP), BP
  0x48839e		4883c440		ADDQ $0x40, SP
  0x4883a2		c3			RET
func main() {
  0x4883a3		e83869fcff		CALL runtime.morestack_noctxt(SB)	; 申请更多的栈空间的地方，也是 goroutine 抢占的检查点
  0x4883a8		e953ffffff		JMP main.main(SB)












TEXT main.main.func1(SB) /golang/src/jingwei.link/main.go
	myfunc := func() {
  0x488420		64488b0c25f8ffffff	MOVQ FS:0xfffffff8, CX
  0x488429		488d4424a8		LEAQ -0x58(SP), AX
  0x48842e		483b4110		CMPQ 0x10(CX), AX
  0x488432		0f86ab010000		JBE 0x4885e3		; 上面三是对栈进行扩容判定，如果栈不够用了，会进行扩容
  0x488438		4881ecd8000000		SUBQ $0xd8, SP		; 预留出 0xd8 的栈空间供 func1(myfunc) 函数使用
  0x48843f		4889ac24d0000000	MOVQ BP, 0xd0(SP)
  0x488447		488dac24d0000000	LEAQ 0xd0(SP), BP	; 上面两句待探究，应该是为了保存某个场景为恢复某个状态做准备
  ; 下面重点关注 DX 的值，是 main.mian 中 0x20(SP) 的地址（区别于本函数的 SP 地址，本函数的 SP 地址已经由 SUBQ 改变过了）
  0x48844f		488b4208		MOVQ 0x8(DX), AX	; 0x8(DX)，其实就是 main.main 中的 0x28(SP)，即 a1 的地址，把这个地址里的值赋值到 AX
  0x488453		4889842480000000	MOVQ AX, 0x80(SP)	; 把 a1 的值赋值到 0x80(SP)
  0x48845b		488b4210		MOVQ 0x10(DX), AX	; 0x10(DX)，其实就是 main.main 中的 0x30(SP)，即 a2 的地址，把这个地址里的值赋值到 AX
  0x48845f		4889442478		MOVQ AX, 0x78(SP)	; 把 a2 的值赋值到 0x80(SP)
	sum := a1 + a2
  0x488464		488b8c2480000000	MOVQ 0x80(SP), CX	; 接下来就是很容易理解的加法运算了
  0x48846c		488b09			MOVQ 0(CX), CX
  0x48846f		480308			ADDQ 0(AX), CX
  0x488472		48894c2440		MOVQ CX, 0x40(SP)
	fmt.Printf("a1: %d, a2:%d, sum: %d\n", a1, a2, sum)
; 再往下就是复杂的 fmt.Printf 函数了，代码很长很臭，就不贴了
