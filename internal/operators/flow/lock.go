package main

import (
	"github.com/emil14/neva/internal/core"
)

func Lock(io core.IO) error {
	sig, err := io.In.Port("sig")
	if err != nil {
		return err
	}

	data, err := io.In.Port("data")
	if err != nil {
		return err
	}

	out, err := io.Out.Port("out")
	if err != nil {
		return err
	}

	go func() {
		for msg := range data {
			<-sig
			out <- msg
		}
	}()

	return nil
}
