package main

import "github.com/emil14/neva/internal/old/runtime"

func Filter(io runtime.IO) error {
	data, err := io.In.Port(runtime.PortAddr{Port: "data"})
	if err != nil {
		return err
	}

	marker, err := io.In.Port(runtime.PortAddr{Port: "marker"})
	if err != nil {
		return err
	}

	acc, err := io.Out.Port(runtime.PortAddr{Port: "acc"})
	if err != nil {
		return err
	}

	rej, err := io.Out.Port(runtime.PortAddr{Port: "rej"})
	if err != nil {
		return err
	}

	go func() {
		for {
			d := <-data
			m := <-marker

			if m.Bool() {
				acc <- d
				continue
			}

			rej <- d
		}
	}()

	return nil
}
