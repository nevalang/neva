package main

import (
	"context"
	"fmt"
)

func main() {
	c := MustCreateCompiler()

	fmt.Println(
		c.Compile(
			context.Background(),
			helloWorld(),
		),
	)
}
