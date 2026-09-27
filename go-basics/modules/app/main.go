package main

import (
	"fmt"

	"example.com/mymath"
)

func main() {
	a, b := 7, 5

	fmt.Printf("%d + %d = %d\n", a, b, mymath.Add(a, b))
	fmt.Printf("%d * %d = %d\n", a, b, mymath.Multiply(a, b))
}
