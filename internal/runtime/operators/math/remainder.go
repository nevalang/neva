package main

import (
	"github.com/emil14/respect/internal/runtime"
)

func Remainder(io runtime.IO) error {
	a, err := io.In.Port(runtime.PortAddr{Port: "a"})
	if err != nil {
		return err
	}

	b, err := io.In.Port(runtime.PortAddr{Port: "b"})
	if err != nil {
		return err
	}

	out, err := io.Out.Port(runtime.PortAddr{Port: "out"})
	if err != nil {
		return err
	}

	go func() {
		for {
			msgA := <-a
			msgB := <-b
			out <- runtime.NewIntMsg(
				msgA.Int() % msgB.Int(),
			)
		}
	}()

	return nil
}
