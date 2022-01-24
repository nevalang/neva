package main

import (
	"errors"
	"fmt"

	"github.com/emil14/neva/internal/new/core"
)

var ErrMul = errors.New("multiplication")

func Mul(io core.IO) error {
	in, err := io.In.ArrPort("in")
	if err != nil {
		return fmt.Errorf("%w: %v", ErrMul, err)
	}

	out, err := io.Out.Port("out")
	if err != nil {
		return fmt.Errorf("%w: %v", ErrMul, err)
	}

	go func() {
		for {
			buf := make(chan int, len(in))

			for i := range in {
				port := in[i]
				go func() {
					msg := <-port
					buf <- msg.Int()
				}()
			}

			mul := 1
			for i := 0; i < len(in); i++ {
				mul *= <-buf
			}

			close(buf)

			out <- core.NewIntMsg(mul)
		}
	}()

	return nil
}
