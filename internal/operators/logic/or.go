package main

import "github.com/emil14/neva/internal/core"

func Or(io core.IO) error {
	in, err := io.In.ArrPortSlots("in")
	if err != nil {
		return err
	}

	out, err := io.Out.Port("out")
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

			out <- core.NewBoolMsg(res)
		}
	}()

	return nil
}
