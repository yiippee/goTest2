package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/scanner"
	"go/token"
	"go/types"
	"log"
	"reflect"
	"strconv"
)

//
//// The playground now supports parentheses or square brackets (only one at
//// a time) for generic type and function declarations and instantiations.
//// By default, parentheses are expected. To switch to square brackets,
//// the first generic declaration in the source must use square brackets.
//
//type Ordered interface {
//type int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, uintptr, float32, float64, string
//}
//
//func Print[type T](s []T) {
//	for _, v := range s {
//		fmt.Print(v)
//	}
//}
//
//// Map 对 []T1 的每个元素执行函数 f 得到新的 []T2
//func Map[T1, T2 Ordered](s []T1, f func(T1) T2) []T2 {
//	r := make([]T2, len(s))
//	for i, v := range s {
//		r[i] = f(v)
//	}
//	return r
//}
//
//func generics() {
//	s := []int{1, 2, 3}
//
//	floats := Map[int, float64](s, func(i int) float64 { return float64(i + 1) })
//	// 现在 floats 的值是 []float64{1.0, 2.0, 3.0}.
//	fmt.Println(floats)
//
//	floats2 := Map(s, func(i int) float64 { return float64(i + 2) })
//	fmt.Println(floats2)
//}

const (
	_MaxSmallSize   = 32768
	smallSizeDiv    = 8
	smallSizeMax    = 1024
	largeSizeDiv    = 128
	_NumSizeClasses = 67
	_PageShift      = 13
)

var class_to_size = [_NumSizeClasses]uint16{0, 8, 16, 32, 48, 64, 80, 96, 112, 128, 144, 160, 176, 192, 208, 224, 240, 256, 288, 320, 352, 384, 416, 448, 480, 512, 576, 640, 704, 768, 896, 1024, 1152, 1280, 1408, 1536, 1792, 2048, 2304, 2688, 3072, 3200, 3456, 4096, 4864, 5376, 6144, 6528, 6784, 6912, 8192, 9472, 9728, 10240, 10880, 12288, 13568, 14336, 16384, 18432, 19072, 20480, 21760, 24576, 27264, 28672, 32768}
var class_to_allocnpages = [_NumSizeClasses]uint8{0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 2, 1, 2, 1, 2, 1, 3, 2, 3, 1, 3, 2, 3, 4, 5, 6, 1, 7, 6, 5, 4, 3, 5, 7, 2, 9, 7, 5, 8, 3, 10, 7, 4}

type divMagic struct {
	shift    uint8
	shift2   uint8
	mul      uint16
	baseMask uint16
}

var class_to_divmagic = [_NumSizeClasses]divMagic{{0, 0, 0, 0}, {3, 0, 1, 65528}, {4, 0, 1, 65520}, {5, 0, 1, 65504}, {4, 11, 683, 0}, {6, 0, 1, 65472}, {4, 10, 205, 0}, {5, 9, 171, 0}, {4, 11, 293, 0}, {7, 0, 1, 65408}, {4, 13, 911, 0}, {5, 10, 205, 0}, {4, 12, 373, 0}, {6, 9, 171, 0}, {4, 13, 631, 0}, {5, 11, 293, 0}, {4, 13, 547, 0}, {8, 0, 1, 65280}, {5, 9, 57, 0}, {6, 9, 103, 0}, {5, 12, 373, 0}, {7, 7, 43, 0}, {5, 10, 79, 0}, {6, 10, 147, 0}, {5, 11, 137, 0}, {9, 0, 1, 65024}, {6, 9, 57, 0}, {7, 9, 103, 0}, {6, 11, 187, 0}, {8, 7, 43, 0}, {7, 8, 37, 0}, {10, 0, 1, 64512}, {7, 9, 57, 0}, {8, 6, 13, 0}, {7, 11, 187, 0}, {9, 5, 11, 0}, {8, 8, 37, 0}, {11, 0, 1, 63488}, {8, 9, 57, 0}, {7, 10, 49, 0}, {10, 5, 11, 0}, {7, 10, 41, 0}, {7, 9, 19, 0}, {12, 0, 1, 61440}, {8, 9, 27, 0}, {8, 10, 49, 0}, {11, 5, 11, 0}, {7, 13, 161, 0}, {7, 13, 155, 0}, {8, 9, 19, 0}, {13, 0, 1, 57344}, {8, 12, 111, 0}, {9, 9, 27, 0}, {11, 6, 13, 0}, {7, 14, 193, 0}, {12, 3, 3, 0}, {8, 13, 155, 0}, {11, 8, 37, 0}, {14, 0, 1, 49152}, {11, 8, 29, 0}, {7, 13, 55, 0}, {12, 5, 7, 0}, {8, 14, 193, 0}, {13, 3, 3, 0}, {7, 14, 77, 0}, {12, 7, 19, 0}, {15, 0, 1, 32768}}
var size_to_class8 = [smallSizeMax/smallSizeDiv + 1]uint8{0, 1, 2, 3, 3, 4, 4, 5, 5, 6, 6, 7, 7, 8, 8, 9, 9, 10, 10, 11, 11, 12, 12, 13, 13, 14, 14, 15, 15, 16, 16, 17, 17, 18, 18, 18, 18, 19, 19, 19, 19, 20, 20, 20, 20, 21, 21, 21, 21, 22, 22, 22, 22, 23, 23, 23, 23, 24, 24, 24, 24, 25, 25, 25, 25, 26, 26, 26, 26, 26, 26, 26, 26, 27, 27, 27, 27, 27, 27, 27, 27, 28, 28, 28, 28, 28, 28, 28, 28, 29, 29, 29, 29, 29, 29, 29, 29, 30, 30, 30, 30, 30, 30, 30, 30, 30, 30, 30, 30, 30, 30, 30, 30, 31, 31, 31, 31, 31, 31, 31, 31, 31, 31, 31, 31, 31, 31, 31, 31}
var size_to_class128 = [(_MaxSmallSize-smallSizeMax)/largeSizeDiv + 1]uint8{31, 32, 33, 34, 35, 36, 36, 37, 37, 38, 38, 39, 39, 39, 40, 40, 40, 41, 42, 42, 43, 43, 43, 43, 43, 44, 44, 44, 44, 44, 44, 45, 45, 45, 45, 46, 46, 46, 46, 46, 46, 47, 47, 47, 48, 48, 49, 50, 50, 50, 50, 50, 50, 50, 50, 50, 50, 51, 51, 51, 51, 51, 51, 51, 51, 51, 51, 52, 52, 53, 53, 53, 53, 54, 54, 54, 54, 54, 55, 55, 55, 55, 55, 55, 55, 55, 55, 55, 55, 56, 56, 56, 56, 56, 56, 56, 56, 56, 56, 57, 57, 57, 57, 57, 57, 58, 58, 58, 58, 58, 58, 58, 58, 58, 58, 58, 58, 58, 58, 58, 58, 59, 59, 59, 59, 59, 59, 59, 59, 59, 59, 59, 59, 59, 59, 59, 59, 60, 60, 60, 60, 60, 61, 61, 61, 61, 61, 61, 61, 61, 61, 61, 61, 62, 62, 62, 62, 62, 62, 62, 62, 62, 62, 63, 63, 63, 63, 63, 63, 63, 63, 63, 63, 63, 63, 63, 63, 63, 63, 63, 63, 63, 63, 63, 63, 64, 64, 64, 64, 64, 64, 64, 64, 64, 64, 64, 64, 64, 64, 64, 64, 64, 64, 64, 64, 64, 65, 65, 65, 65, 65, 65, 65, 65, 65, 65, 65, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66, 66}

func test_mspan() {
	s := 32 / 8
	s2 := size_to_class8[s]        // 得出这块内存大小所在的索引值是多少
	s3 := class_to_size[s2]        // 标识这个索引值对应的内存大小是多少，（会有一些浪费，比如需要30B，但也需要分配32B）
	s4 := class_to_allocnpages[s2] // 标识这个mspan由多少个page（8k）组成，大部分都是1页啊

	fmt.Println(s2, s3, s4)
}

//
//func Eval(exp ast.Expr) float64 {
//	switch exp := exp.(type) {
//	case *ast.BinaryExpr:
//		return EvalBinaryExpr(exp)
//	case *ast.BasicLit:
//		f, _ := strconv.ParseFloat(exp.Value, 64)
//		return f
//	}
//	return 0
//}
//
//// 二元表达式求值
//func EvalBinaryExpr(exp *ast.BinaryExpr) float64 {
//	switch exp.Op {
//	case token.ADD:
//		return Eval(exp.X) + Eval(exp.Y)
//	case token.MUL:
//		return Eval(exp.X) * Eval(exp.Y)
//	}
//	return 0
//}

// 带有变量的二元表达式求值
func Eval(exp ast.Expr, vars map[string]float64) float64 {
	switch exp := exp.(type) {
	case *ast.BinaryExpr:
		return EvalBinaryExpr(exp, vars)
	case *ast.BasicLit:
		f, _ := strconv.ParseFloat(exp.Value, 64)
		return f
	case *ast.Ident: // 代表一个标识符 identifier
		return vars[exp.Name] // 在Eval函数递归解析时，如果当前解析的表达式语法树结点是*ast.Ident类型，则直接从vars表格查询结果
	}
	return 0
}

func EvalBinaryExpr(exp *ast.BinaryExpr, vars map[string]float64) float64 {
	switch exp.Op {
	case token.ADD:
		return Eval(exp.X, vars) + Eval(exp.Y, vars)
	case token.MUL:
		return Eval(exp.X, vars) * Eval(exp.Y, vars)
	}
	return 0
}

func main() {
	// tcmalloc 内存分配
	test_mspan()

	// 类型检查
	main4()

	// 复合类型
	main3()

	// 包文件解析，函数、方法解析
	main2()

	// 引入变量
	expr2, _ := parser.ParseExpr(`1+2*3+x`)
	fmt.Println(Eval(expr2, map[string]float64{
		"x": 100,
	}))

	// 二元表达式求值
	expr, _ := parser.ParseExpr(`1+2*3`)
	ast.Print(nil, Eval(expr, nil))

	// 字面量
	var lit9527 = &ast.BasicLit{
		Kind:  token.INT,
		Value: "9527",
	}
	ast.Print(nil, lit9527)

	t := reflect.TypeOf(lit9527)
	fmt.Println(t)

	// token解析
	var src = []byte(`println("你好，世界") // test`)

	var fset = token.NewFileSet()
	var file = fset.AddFile("hello.go", fset.Base(), len(src))

	var s scanner.Scanner
	s.Init(file, src, nil, scanner.ScanComments)

	for {
		pos, tok, lit := s.Scan()
		if tok == token.EOF {
			break
		}
		fmt.Printf("%s\t%s\t%q\n", fset.Position(pos), tok, lit)
	}
}

func main2() {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "hello.go", src, parser.AllErrors)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(f)
}

const src = `package pkgname

import ("a"; "b")
type SomeType int
const PI = 3.14
var Length = 1

func main() {
	for {}
	for true {}
	for i := 0; true; i++ {}
	for i, v := range m {}
}

func (p *xType) Hello(arg1, arg2 int) (bool, error) { 
	a := 1
	fmt.Println(a)
}
`

func main3() {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "hello.go", src3, parser.AllErrors)
	if err != nil {
		log.Fatal(err)
	}

	for _, decl := range f.Decls {
		ast.Print(nil, decl.(*ast.GenDecl).Specs[0])
	}
}

const src3 = `package foo
type Int1 int
type Int2 pkg.int
type IntPtr *int
type IntPtrPtr **int

type MyStruct struct {
	a, b int "int value"
	string
}

type FuncType func(a, b int) bool

type IntReader interface {
	Read() int
}

type IntReader2 struct {
	Read func() int // 与接口的定义很像了
}

`

func main4() {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "hello.go", src4, parser.AllErrors)
	if err != nil {
		log.Fatal(err)
	}

	pkg, err := new(types.Config).Check("hello.go", fset, []*ast.File{f}, nil)
	if err != nil {
		log.Fatal(err)
	}

	_ = pkg
}

const src4 = `package pkg

func hello() {
	var _ = "a" + 1
}
`
