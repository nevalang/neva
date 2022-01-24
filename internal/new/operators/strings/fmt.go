package main

import (
	"errors"
	"fmt"

	"github.com/emil14/neva/internal/new/core"
)

var ErrFmt = errors.New("fmt")

func Fmt(io core.IO) error {
	str, err := io.In.Port("str")
	if err != nil {
		return fmt.Errorf("%w: %v", ErrFmt, err)
	}

	args, err := io.In.ArrPort("args")
	if err != nil {
		return fmt.Errorf("%w: %v", ErrFmt, err)
	}

	out, err := io.Out.Port("out")
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

			out <- core.NewStrMsg(
				fmt.Sprintf(s.Str(), ss...),
			)
		}
	}()

	return nil
}
