package main

import (
	"errors"
	"fmt"

	"github.com/emil14/neva/internal/new/core"
)

var ErrMul = errors.New("multiplication")

func Mul(io core.IO) error {
	inports, err := io.In.ArrPort("in")
	if err != nil {
		return fmt.Errorf("%w: %v", ErrMul, err)
	}

	out, err := io.Out.Port("out")
	if err != nil {
		return fmt.Errorf("%w: %v", ErrMul, err)
	}

	go func() {
		for {
			mul := 1
			for i := range inports {
				mul *= (<-inports[i]).Int()
			}
			out <- core.NewIntMsg(mul)
		}
	}()

	return nil
}
