package main

import (
	"github.com/emil14/respect/internal/core"
	"github.com/emil14/respect/internal/runtime"
)

func Or(io core.IO) error {
	in, err := io.In.Group("in")
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

			out <- runtime.NewBoolMsg(res)
		}
	}()

	return nil
}
