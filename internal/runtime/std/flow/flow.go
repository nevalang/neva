package main

import (
	"context"

	"github.com/emil14/neva/internal/runtime"
)

func Trigger(ctx context.Context, io runtime.IO) error {
	slots, err := io.In.ArrPort("sigs")
	if err != nil {
		return err
	}

	v, err := io.In.Port("v")
	if err != nil {
		return err
	}

	out, err := io.Out.Port("v")
	if err != nil {
		return err
	}

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			for i := range slots {
				select {
				case <-ctx.Done():
					return ctx.Err()
				case <-slots[i]:
					msg := <-v
					out <- msg
				}
			}
		}
	}
}
