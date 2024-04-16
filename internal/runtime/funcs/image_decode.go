package funcs

import (
	"context"
	"image"
	"image/jpeg"
	"image/png"
	"strings"

	"github.com/nevalang/neva/internal/runtime"
)

type imageDecode struct{}

func (imageDecode) Create(io runtime.FuncIO, _ runtime.Msg) (func(ctx context.Context), error) {
	fmt, err := io.In.Port("fmt")
	if err != nil {
		return nil, err
	}

	data, err := io.In.Port("data")
	if err != nil {
		return nil, err
	}

	bounds, err := io.Out.Port("bounds")
	if err != nil {
		return nil, err
	}

	seq, err := io.Out.Port("seq")
	if err != nil {
		return nil, err
	}

	return func(ctx context.Context) {
		for {
			s := <-data
			if s == nil {
				// Something went wrong up the pipe.
				// Send a nil bounds and sentinel Pixel.
				bounds <- nil
				seq <- nil
				continue
			}
			dat := s.Str()
			var (
				f  formatMsg
				im image.Image
			)
			// Default to Raw if error.
			switch f.decode(<-fmt); f {
			case 1: // JPEG
				var err error
				if im, err = jpeg.Decode(strings.NewReader(dat)); err != nil {
					// Something went wrong. Send nils.
					bounds <- nil
					seq <- nil
					continue
				}
			default: // Raw, PNG
				var err error
				if im, err = png.Decode(strings.NewReader(dat)); err != nil {
					// Something went wrong. Send nils.
					bounds <- nil
					seq <- nil
					continue
				}
			}

			// Send the image properties.
			var b boundsMsg
			b.decodeRect(im.Bounds())
			bounds <- b.encode()

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
