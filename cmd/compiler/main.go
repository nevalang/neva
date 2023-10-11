package main

import (
	"context"
	"os"

	"github.com/nevalang/neva/internal/compiler"
)

func main() {
	comp := compiler.New(nil, nil, nil, nil, nil, nil, nil)
	if _, err := comp.Compile(context.Background(), os.Args[1], os.Args[2]); err != nil {
		panic(err)
	}
}
