package main

import "github.com/emil14/neva/internal/runtime"

func Filter(io runtime.IO) error {
	data, err := io.In.Port("data")
	if err != nil {
		return err
	}

	marker, err := io.In.Port("data")
	if err != nil {
		return err
	}

	acc, err := io.Out.Port("accepted")
	if err != nil {
		return err
	}

	rej, err := io.Out.Port("rejected")
	if err != nil {
		return err
	}

	go func() {
		for {
			dataMsg := <-data
			markerMsg := <-marker

			if markerMsg.Bool() {
				acc <- dataMsg
				continue
			}

			rej <- dataMsg
		}
	}()

	return nil
}
