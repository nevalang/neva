package main

import "github.com/emil14/neva/internal/old/runtime"

func More(io runtime.IO) error {
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
			out <- runtime.NewBoolMsg(
				msgA.Int() > msgB.Int(),
			)
		}
	}()

	return nil
}
