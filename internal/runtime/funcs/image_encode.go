package funcs

import (
	"context"
	"image/png"
	"strings"

	"github.com/nevalang/neva/internal/runtime"
)

type imageEncode struct{}

func (imageEncode) Create(io runtime.IO, _ runtime.Msg) (func(ctx context.Context), error) {
	imgIn, err := io.In.Single("img")
	if err != nil {
		return nil, err
	}

	dataOut, err := io.Out.Single("data")
	if err != nil {
		return nil, err
	}

	errOut, err := io.Out.Single("err")
	if err != nil {
		return nil, err
	}

	return func(ctx context.Context) {
		for {
			imgMsg, ok := imgIn.Receive(ctx)
			if !ok {
				return
			}

			var b imageMsg
			b.decode(imgMsg.Map())

			im := b.createImage()

			// Encode the image in the desired format to sb.
			var sb strings.Builder // for encoded output.
			if err := png.Encode(&sb, im); err != nil {
				if !errOut.Send(ctx, errFromErr(err)) {
					return
				}
			}

			if !dataOut.Send(
				ctx,
				runtime.NewStrMsg(sb.String()),
			) {
				return
			}
		}
	}, nil
}
