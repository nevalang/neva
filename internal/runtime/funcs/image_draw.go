package funcs

import (
	"context"
	"image"

	"github.com/nevalang/neva/internal/runtime"
)

type imageDraw struct{}

func (imageDraw) Create(io runtime.FuncIO, _ runtime.Msg) (func(ctx context.Context), error) {
	_, err := io.In.Port("src")
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

			// Create the src image.

			// Draw over the image.
			for e := range over {
				var p pixelMsg
				if e == nil {
					break // Sentinel.
				}
				if p.decode(e) {
				}
			}

			// Send raw Pixels.
			im := image.RGBA{}
			var p pixelMsg
			for x := 0; x < im.Bounds().Dx(); x++ {
				for y := 0; y < im.Bounds().Dy(); y++ {
					c := im.At(x, y)
					p.color.decodeColor(c)
					seq <- p.encode()
				}
			}
			// Send a sentinel Pixel.
			seq <- nil
		}
	}, nil
}
