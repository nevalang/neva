package main

import (
	"github.com/emil14/neva/internal/core"
)

func Lock(io core.IO) error {
	sig, err := io.In.Port("sig")
	if err != nil {
		return err
	}

	dataIn, err := io.In.Port("data")
	if err != nil {
		return err
	}

	dataOut, err := io.Out.Port("data")
	if err != nil {
		return err
	}

	go func() {
		for {
			<-sig // read sig first to avoid unnecessary sendings
			msg := <-dataIn
			dataOut <- msg
		}
	}()

	return nil
}
