package main

import "github.com/emil14/neva/internal/new/core"

func Filter(io core.IO) error {
	dataIn, err := io.In.Port("data")
	if err != nil {
		return err
	}

	condIn, err := io.In.Port("cond")
	if err != nil {
		return err
	}

	dataOut, err := io.Out.Port("data")
	if err != nil {
		return err
	}

	go func() {
		for {
			data := <-dataIn
			cond := <-condIn

			if core.Equal(data, cond) {
				dataOut <- data
			}
		}
	}()

	return nil
}
