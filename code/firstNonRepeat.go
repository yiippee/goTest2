package main

import "fmt"

func main() {
	arr := []int{
		9, 4, 9, 6, 7, 4,
	}

	m := make(map[int]int)
	for _, v := range arr {
		m[v]++
	}

	for _, v := range m {
		if v == 1 {
			fmt.Println(v)
		}
	}
}
