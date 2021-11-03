package main

import (
	"github.com/emil14/respect/internal/core"
)

func Mul(io core.IO) error {
	inportGroup, err := io.In.Group("in")
	if err != nil {
		return err
	}

	out, err := io.Out.Port("out")
	if err != nil {
		return err
	}

	go func() {
		for {
			buf := make(chan int, len(inportGroup))

			for i := range inportGroup {
				inport := inportGroup[i]
				go func() {
					buf <- (<-inport).Int()
				}()
			}

			mul := 1
			for i := 0; i < len(inportGroup); i++ {
				mul *= <-buf
			}

			close(buf)

			out <- core.NewIntMsg(mul)
		}
	}()

	return nil
}
