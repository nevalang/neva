package main

import "github.com/emil14/neva/internal/old/runtime"

func Or(io runtime.IO) error {
	in, err := io.In.PortArray("in")
	if err != nil {
		return err
	}

	out, err := io.Out.Port(runtime.PortAddr{Port: "out"})
	if err != nil {
		return err
	}

	go func() {
		for {
			buf := make(chan bool, len(in))
			for _, ch := range in {
				msg := <-ch
				buf <- msg.Bool()
			}

			res := false
			for b := range buf {
				if b {
					res = true
					break
				}
			}

			close(buf)

			out <- runtime.NewBoolMsg(res)
		}
	}()

	return nil
}
