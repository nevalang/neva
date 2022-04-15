package main

import "github.com/emil14/neva/internal/core"

func Remainder(io core.IO) error {
	a, err := io.In.Port("a")
	if err != nil {
		return err
	}

	b, err := io.In.Port("b")
	if err != nil {
		return err
	}

	out, err := io.Out.Port("out")
	if err != nil {
		return err
	}

	go func() {
		for {
			msgA := <-a
			msgB := <-b
			out <- core.NewIntMsg(
				msgA.Int() % msgB.Int(),
			)
		}
	}()

	return nil
}
