package main

import (
	"fmt"

	"github.com/emil14/neva/internal/core"
)

func Print(io core.IO) error {
	dataIn, err := io.In.Port("data")
	if err != nil {
		return err
	}

	dataOut, err := io.Out.Port("data")
	if err != nil {
		return err
	}

	go func() {
		for msg := range dataIn {
			fmt.Print(msg)
			dataOut <- msg
		}
	}()

	return nil
}
