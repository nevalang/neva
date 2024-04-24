package funcs

import (
	"context"
	"image/png"
	"strings"

	"github.com/nevalang/neva/internal/runtime"
)

type imageEncode struct{}

func (imageEncode) Create(io runtime.FuncIO, _ runtime.Msg) (func(ctx context.Context), error) {
	in, err := io.In.Port("img")
	if err != nil {
		return nil, err
	}

	data, err := io.Out.Port("data")
	if err != nil {
		return nil, err
	}

	errCh, err := io.Out.Port("err")
	if err != nil {
		return nil, err
	}

	return func(ctx context.Context) {
		for {
			var b imageMsg
			b.decode(<-in)
			im := b.createImage()
			// Encode the image in the desired format to sb.
			var sb strings.Builder // for encoded output.
			if err := png.Encode(&sb, im); err != nil {
				// Something went wrong. Send err.
				errCh <- runtime.NewMapMsg(map[string]runtime.Msg{
					"error": runtime.NewStrMsg(err.Error()),
				})
				continue
			}
			// Send the image.
			data <- runtime.NewStrMsg(sb.String())
		}
	}, nil
}
