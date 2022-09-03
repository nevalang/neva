package main

import (
	"errors"
	"fmt"

	"github.com/emil14/neva/internal/core"
)

var ErrFmt = errors.New("fmt")

func Fmt(io core.IO) error {
	strPort, err := io.In.Port("str")
	if err != nil {
		return fmt.Errorf("%w: %v", ErrFmt, err)
	}

	argPorts, err := io.In.ArrPort("args")
	if err != nil {
		return fmt.Errorf("%w: %v", ErrFmt, err)
	}

	outPort, err := io.Out.Port("out")
	if err != nil {
		return fmt.Errorf("%w: %v", ErrFmt, err)
	}

	go func() {
		for msg := range strPort {
			ss := make([]interface{}, 0, len(argPorts))

			for i := range argPorts {
				arg := <-argPorts[i]
				ss = append(ss, arg)
			}

			outPort <- core.NewStrMsg(
				fmt.Sprintf(msg.Str(), ss...),
			)
		}
	}()

	return nil
}
