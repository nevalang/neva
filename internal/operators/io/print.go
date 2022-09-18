package main

import (
	"context"
	"fmt"

	"github.com/emil14/neva/internal/core"
)

func Print(io core.IO) (func(context.Context), error) {
	dataIn, err := io.In.Port("data")
	if err != nil {
		return nil, err
	}

	dataOut, err := io.Out.Port("data")
	if err != nil {
		return nil, err
	}

	return func(ctx context.Context) {
		for {
			select {
			case <-ctx.Done():
				return
			case msg := <-dataIn:
				fmt.Println(<-dataIn)
				dataOut <- msg
			}
		}
	}, nil
}
