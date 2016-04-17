package main

import "fmt"

func main() {
	pow := make([]int, 10)
	for i := range pow {
		pow[i] = 1 << uint(i) // == 2**i
	}
	for i, v := range pow {
		fmt.Printf("2**%d = %d\n", i, v)
	}
	fmt.Println("Values:")
	for _, value := range pow {
		fmt.Printf("%d\n", value)
	}
	fmt.Println("Indices:")
	for i := range pow {
		fmt.Printf("%d\n", i)
	}
}
