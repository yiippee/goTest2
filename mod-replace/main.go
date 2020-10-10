package main

import (
	"fmt"

	"github.com/lzb/replace"
)

func main() {
	fmt.Println("test")
	replace.Test()

	a := testmod.A
	fmt.Println(a)
}
