package main

import "github.com/emil14/neva/internal/core"

func And(io core.IO) error {
	inports, err := io.In.ArrPort("in")
	if err != nil {
		return err
	}

	outport, err := io.Out.Port("out")
	if err != nil {
		return err
	}

	go func() {
		for {
			res := true

			for _, port := range inports {
				msg := <-port
				if !msg.Bool() {
					res = false
					break
				}
			}

			outport <- core.NewBoolMsg(res)
		}
	}()

	return nil
}
