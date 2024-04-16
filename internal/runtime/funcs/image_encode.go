package funcs

import (
	"context"
	"image"
	"image/jpeg"
	"image/png"
	"strings"

	"github.com/nevalang/neva/internal/runtime"
)

type imageEncode struct{}

func (imageEncode) Create(io runtime.FuncIO, _ runtime.Msg) (func(ctx context.Context), error) {
	fmt, err := io.In.Port("fmt")
	if err != nil {
		return nil, err
	}

	bounds, err := io.In.Port("bounds")
	if err != nil {
		return nil, err
	}

	seq, err := io.In.Port("seq")
	if err != nil {
		return nil, err
	}

	data, err := io.Out.Port("data")
	if err != nil {
		return nil, err
	}

	return func(ctx context.Context) {
		for {
			var b boundsMsg
			// Default to 0x0 image if error.
			b.decode(<-bounds)
			// Create an RGBA image and set individual pixels.
			im := image.NewRGBA64(b.rect())
			for e := range seq {
				var p pixelMsg
				if p.decode(e) {
					im.Set(int(p.point.x), int(p.point.y), p.color.color())
				}
			}
			// Encode the image in the desired format to sb.
			var (
				f  formatMsg
				sb strings.Builder
			)
			// Default to Raw if error.
			switch f.decode(<-fmt); f {
			case 1: // JPEG
				if err := jpeg.Encode(&sb, im, &jpeg.Options{Quality: jpeg.DefaultQuality}); err != nil {
					// Something went wrong. Send a nil image.
					data <- nil
					continue
				}
			default: // Raw, PNG
				var sb strings.Builder
				if err := png.Encode(&sb, im); err != nil {
					// Something went wrong. Send a nil image.
					data <- nil
					continue
				}
			}
			// Send the image.
			data <- runtime.NewStrMsg(sb.String())
		}
	}, nil
}
