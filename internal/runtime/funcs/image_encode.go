package funcs

import (
	"context"
	"image/png"
	"strings"

	"github.com/nevalang/neva/internal/runtime"
)

type imageEncode struct{}

func (imageEncode) Create(io runtime.FuncIO, _ runtime.Msg) (func(ctx context.Context), error) {
	in, err := io.In.SingleInport("img")
	if err != nil {
		return nil, err
	}

	data, err := io.Out.SingleOutport("data")
	if err != nil {
		return nil, err
	}

	errCh, err := io.Out.SingleOutport("err")
	if err != nil {
		return nil, err
	}

	return func(ctx context.Context) {
		for {
			var b imageMsg
			select {
			case m := <-in:
				b.decode(m)
			case <-ctx.Done():
				return
			}
			im := b.createImage()
			// Encode the image in the desired format to sb.
			var sb strings.Builder // for encoded output.
			if err := png.Encode(&sb, im); err != nil {
				// Something went wrong. Send err.
				select {
				case errCh <- runtime.NewMapMsg(map[string]runtime.Msg{
					"error": runtime.NewStrMsg(err.Error()),
				}):
					continue
				case <-ctx.Done():
					return
				}
			}
			// Send the image.
			select {
			case data <- runtime.NewStrMsg(sb.String()):
			case <-ctx.Done():
				return
			}
		}
	}, nil
}
