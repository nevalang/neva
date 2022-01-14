package main

import (
	"errors"
	"fmt"

	"github.com/emil14/neva/internal/runtime"
)

var ErrLock = errors.New("lock")

func Lock(io runtime.IO) error {
	sig, err := io.In.Port(runtime.PortAddr{Port: "sig"})
	if err != nil {
		return fmt.Errorf("%w: %v", ErrLock, err)
	}

	data, err := io.In.Port(runtime.PortAddr{Port: "data"})
	if err != nil {
		return fmt.Errorf("%w: %v", ErrLock, err)
	}

	out, err := io.Out.Port(runtime.PortAddr{Port: "data"})
	if err != nil {
		return fmt.Errorf("%w: %v", ErrLock, err)
	}

	go func() {
		for msg := range data {
			<-sig
			out <- msg
		}
	}()

	return nil
}
