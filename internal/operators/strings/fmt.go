package main

import (
	"fmt"

	"github.com/emil14/neva/internal/core"
)

func Fmt(io core.IO) error { // TODO rename maybe?
	strPort, err := io.In.Port("str")
	if err != nil {
		return err
	}

	argPortSlots, err := io.In.ArrPortSlots("args")
	if err != nil {
		return err
	}

	outPort, err := io.Out.Port("out")
	if err != nil {
		return err
	}

	go func() {
		for msg := range strPort {
			ss := make([]interface{}, 0, len(argPortSlots))

			for i := range argPortSlots {
				arg := <-argPortSlots[i]
				ss = append(ss, arg)
			}

			outPort <- core.NewStrMsg(
				fmt.Sprintf(msg.Str(), ss...), // FIXME go's interface not the best
			)
		}
	}()

	return nil
}
