package main

import "fmt"

func main() {
	fmt.Println(
		f(-42, 42),
	)
}

type Constraint interface {
	string | int
}

func f[T Constraint](x, y T) T {
	return x + y
}
