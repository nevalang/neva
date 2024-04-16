package funcs

import (
	"context"
	"image"

	"github.com/nevalang/neva/internal/runtime"
)

type imageDraw struct{}

func (imageDraw) Create(io runtime.FuncIO, _ runtime.Msg) (func(ctx context.Context), error) {
	bounds, err := io.In.Port("bounds")
	if err != nil {
		return nil, err
	}

	src, err := io.In.Port("src")
	if err != nil {
		return nil, err
	}

	over, err := io.In.Port("over")
	if err != nil {
		return nil, err
	}

	seq, err := io.Out.Port("seq")
	if err != nil {
		return nil, err
	}

	return func(ctx context.Context) {
		for {
			var b boundsMsg
			b.decode(<-bounds)

			// Create the src image.
			im := image.NewRGBA(b.rect())
			for e := range src {
				var p pixelMsg
				if e == nil {
					break // Sentinel.
				}
				if p.decode(e) {
					im.Set(int(p.point.x), int(p.point.y), p.color.color())
				}
			}

			// Draw over the image.
			for e := range over {
				var p pixelMsg
				if e == nil {
					break // Sentinel.
				}
				if p.decode(e) {
					im.Set(int(p.point.x), int(p.point.y), p.color.color())
				}
			}

			// Send raw Pixels.
			var p pixelMsg
			for x := 0; x < im.Bounds().Dx(); x++ {
				for y := 0; y < im.Bounds().Dy(); y++ {
					c := im.At(x, y)
					p.point = pointMsg{x: int64(x), y: int64(y)}
					p.color.decodeColor(c)
					seq <- p.encode()
				}
			}
			// Send a sentinel Pixel.
			seq <- nil
		}
	}, nil
}
