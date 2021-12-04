package main

import (
	"errors"
	"fmt"

	"github.com/emil14/neva/internal/runtime"
)

var ErrWrite = errors.New("multiplication")

func Write(io runtime.IO) error {
	in, err := io.In.Port(runtime.PortAddr{Port: "in"})
	if err != nil {
		return fmt.Errorf("%w: %v", ErrWrite, err)
	}

	go func() {
		for {
			fmt.Print(<-in)
		}
	}()

	return nil
}
