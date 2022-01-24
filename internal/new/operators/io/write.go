package main

import (
	"fmt"

	"github.com/emil14/neva/internal/new/core"
)

func Write(io core.IO) error {
	in, err := io.In.Port("in")
	if err != nil {
		return err
	}

	go func() {
		for {
			fmt.Print(<-in)
		}
	}()

	return nil
}
