package main

import "fmt"

func main() {
	i, j := 42, 2701

	p := &i         // point to i
	fmt.Println(*p) // read i thro the pointer
	*p = 21         // set i through pointer
	fmt.Println(i)  // see new value of i

	p = &j          // point to j
	*p = *p / 37    // divide j thro pointer
	fmt.Println(j) 
}
