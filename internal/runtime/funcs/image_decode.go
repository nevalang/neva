package funcs

import (
	"context"
	"image"
	"image/png"
	"strings"

	"github.com/nevalang/neva/internal/runtime"
)

type imageDecode struct{}

func (imageDecode) Create(io runtime.FuncIO, _ runtime.Msg) (func(ctx context.Context), error) {
	data, err := io.In.Port("data")
	if err != nil {
		return nil, err
	}

	out, err := io.Out.Port("img")
	if err != nil {
		return nil, err
	}

	errCh, err := io.Out.Port("err")
	if err != nil {
		return nil, err
	}

	return func(ctx context.Context) {
		for {
			s := <-data
			dat := s.Str()
			im, err := png.Decode(strings.NewReader(dat))
			if err != nil {
				// Something went wrong. Send err.
				errCh <- runtime.NewMapMsg(map[string]runtime.Msg{
					"error": runtime.NewStrMsg(err.Error()),
				})
				continue
			}

			// Send the image.
			var i imageMsg
			i.decodeImage(im.(*image.RGBA))
			out <- i.encode()
		}
	}, nil
}
