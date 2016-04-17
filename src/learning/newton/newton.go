package main

import (
	"fmt"
	"math"
)

func Sqrt(x float64) float64 {
	z := 1.0
	for i := 0; i < 10; i++ {
		z = z - (z*z-x)/(2*z)
	}
	return z
}

func main() {
	num := 644.0
	fmt.Println("Estimated: ", Sqrt(num))
	fmt.Println("Actual: ", math.Sqrt(num))
}
