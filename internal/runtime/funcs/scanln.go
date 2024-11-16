package funcs

import (
	"context"
	"fmt"

	"github.com/nevalang/neva/internal/runtime"
)

type scanln struct{}

// TODO add `:err` outport
func (r scanln) Create(rio runtime.IO, _ runtime.Msg) (func(ctx context.Context), error) {
	sigIn, err := rio.In.Single("sig")
	if err != nil {
		return nil, err
	}

	resOut, err := rio.Out.Single("res")
	if err != nil {
		return nil, err
	}

	// errOut, err := rio.Out.Single("err")
	// if err != nil {
	// 	return nil, err
	// }

	return func(ctx context.Context) {
		for {
			if _, ok := sigIn.Receive(ctx); !ok {
				return
			}

			var input string
			if _, err := fmt.Scanln(&input); err != nil {
				// if !errOut.Send(ctx, errFromErr(err)) {
				// 	return
				// }
				continue
			}

			if !resOut.Send(ctx, runtime.NewStringMsg(input)) {
				return
			}
		}
	}, nil
}
