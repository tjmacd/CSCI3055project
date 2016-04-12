package main

import ("fmt"
		"stringutil"
)

func add(x, y int) int {
	return x + y
}

func main() {
	fmt.Println(add(3,4))
	fmt.Println(stringutil.Reverse("!oG ,olleH"))
}
