package funcs

import (
	"context"

	"github.com/nevalang/neva/internal/runtime"
)

type imageNew struct{}

func (imageNew) Create(io runtime.FuncIO, _ runtime.Msg) (func(ctx context.Context), error) {
	width, err := io.In.Port("width")
	if err != nil {
		return nil, err
	}

	height, err := io.In.Port("height")
	if err != nil {
		return nil, err
	}

	out, err := io.Out.Port("img")
	if err != nil {
		return nil, err
	}

	return func(ctx context.Context) {
		for {
			var i imageMsg
			if w := <-width; !decodeInt(&i.width, w) {
				continue
			}
			if h := <-height; !decodeInt(&i.height, h) {
				continue
			}
			out <- i.encode()
		}
	}, nil
}
