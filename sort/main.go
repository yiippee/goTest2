package main

import "sort"

type Person struct {
	Name string
	Age  int
}

func main() {
	keys := []int{1, 2, 3, 4, 5, 6, 7}
	sort.Search(7, func(i int) bool {
		return keys[i] >= 4
	})

	sort.Ints([]int{19, 2, 3, 4, 5, 6, 7})

	str := []string{"a", "f", "b"}
	sort.Strings(str)

	data := []Person{
		{"Alice", 20},
		{"Bob", 15},
		{"Jane", 30},
	}
	sort.Slice(data, func(i, j int) bool {
		return data[i].Age < data[j].Age
	})
}
