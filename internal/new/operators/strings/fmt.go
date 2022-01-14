package main

import (
	"errors"
	"fmt"

	"github.com/emil14/neva/internal/runtime"
)

var ErrFmt = errors.New("fmt")

func Fmt(io runtime.IO) error {
	str, err := io.In.Port(runtime.PortAddr{Port: "str"})
	if err != nil {
		return fmt.Errorf("%w: %v", ErrFmt, err)
	}

	args, err := io.In.PortArray("args")
	if err != nil {
		return fmt.Errorf("%w: %v", ErrFmt, err)
	}

	out, err := io.Out.Port(runtime.PortAddr{Port: "out"})
	if err != nil {
		return fmt.Errorf("%w: %v", ErrFmt, err)
	}

	go func() {
		for s := range str {
			ss := make([]interface{}, len(args))
			for _, ch := range args {
				arg := <-ch
				ss = append(ss, arg)
			}

			out <- runtime.NewStrMsg(
				fmt.Sprintf(s.Str(), ss...),
			)
		}
	}()

	return nil
}
