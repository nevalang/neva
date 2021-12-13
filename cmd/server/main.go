package main

import (
	"fmt"

	"github.com/emil14/neva/internal/server"
)

func main() {
	srv := server.NewServer()
	fmt.Println(srv)
}
