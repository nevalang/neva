package funcs

import (
	"bytes"
	"context"
	"image/png"

	"github.com/nevalang/neva/internal/runtime"
)

type imageEncode struct{}

//nolint:gocognit,varnamelen // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
func (imageEncode) Create(io runtime.IO, _ runtime.Msg) (func(ctx context.Context), error) {
	imgIn, err := io.In.Single("img")
	if err != nil {
		//nolint:wrapcheck // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
		return nil, err
	}

	resOut, err := io.Out.Single("res")
	if err != nil {
		//nolint:wrapcheck // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
		return nil, err
	}

	errOut, err := io.Out.Single("err")
	if err != nil {
		//nolint:wrapcheck // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
		return nil, err
	}

	return func(ctx context.Context) {
		for {
			imgMsg, ok := imgIn.Receive(ctx)
			if !ok {
				return
			}

			imgStructMsg := imgMsg.Struct()
			//nolint:varnamelen // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
			b := imageMsg{
				pixels: imgStructMsg.Get("pixels").Bytes(),
				width:  imgStructMsg.Get("width").Int(),
				height: imgStructMsg.Get("height").Int(),
			}

			im := b.createImage()

			// Encode the image in the desired format to sb.
			//nolint:varnamelen // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
			var sb bytes.Buffer // for encoded output.
			if err := png.Encode(&sb, im); err != nil {
				if !errOut.Send(ctx, errFromErr(err)) {
					return
				}
			}

			if !resOut.Send(
				ctx,
				runtime.NewBytesMsg(sb.Bytes()),
			) {
				return
			}
		}
	}, nil
}
