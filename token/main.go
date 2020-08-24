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
