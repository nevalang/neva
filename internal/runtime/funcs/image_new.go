package funcs

import (
	"context"

	"github.com/nevalang/neva/internal/runtime"
)

type imageNew struct{}

func (imageNew) Create(io runtime.FuncIO, _ runtime.Msg) (func(ctx context.Context), error) {
	bounds, err := io.In.Port("bounds")
	if err != nil {
		return nil, err
	}

	seq, err := io.Out.Port("seq")
	if err != nil {
		return nil, err
	}

	return func(ctx context.Context) {
		for msg := range bounds {
			var b boundsMsg
			if !b.decode(msg) {
				// A empty bounds is like a 0x0 image.
				// Send a sentinel Pixel and continue.
				seq <- nil
				continue
			}
			var dx, dy int64
			if dx = b.max.x - b.min.x; dx < 0 {
				dx = -dx
			}
			if dy = b.max.y - b.min.y; dy < 0 {
				dy = -dy
			}
			// Send raw Pixels.
			im := make([]pixelMsg, dx*dy)
			for _, p := range im {
				seq <- p.encode()
			}
			// Send a sentinel Pixel.
			seq <- nil
		}
	}, nil
}
