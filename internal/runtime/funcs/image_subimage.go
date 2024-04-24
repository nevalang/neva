package funcs

import (
	"context"

	"github.com/nevalang/neva/internal/runtime"
)

type imageSubImage struct{}

func (imageSubImage) Create(io runtime.FuncIO, _ runtime.Msg) (func(ctx context.Context), error) {
	_, err := io.Out.Port("seq")
	if err != nil {
		return nil, err
	}

	return func(ctx context.Context) {
		for {
			// var dx, dy int64
			// if dx = b.max.x - b.min.x; dx < 0 {
			// 	dx = -dx
			// }
			// if dy = b.max.y - b.min.y; dy < 0 {
			// 	dy = -dy
			// }
			// // Send raw Pixels.
			// im := make([]pixelMsg, dx*dy)
			// for _, p := range im {
			// 	seq <- p.encode()
			// }
			// // Send a sentinel Pixel.
			// seq <- nil
		}
	}, nil
}
