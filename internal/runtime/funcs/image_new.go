package funcs

import (
	"context"
	"fmt"
	"image"

	"github.com/nevalang/neva/internal/runtime"
)

type imageNew struct{}

//nolint:cyclop,gocognit,gocyclo // Stream framing, image accumulation, and error forwarding share one lifecycle.
func (imageNew) Create(runtimeIO runtime.IO, _ runtime.Msg) (func(ctx context.Context), error) {
	pixelsIn, err := runtimeIO.In.Single("pixels")
	if err != nil {
		return nil, fmt.Errorf("get pixels inport: %w", err)
	}

	imgOut, err := runtimeIO.Out.Single("img")
	if err != nil {
		return nil, fmt.Errorf("get img outport: %w", err)
	}

	errOut, err := runtimeIO.Out.Single("err")
	if err != nil {
		return nil, fmt.Errorf("get err outport: %w", err)
	}

	return func(ctx context.Context) {
		for {
			if !waitStreamOpen(ctx, pixelsIn) {
				return
			}

			pixels := make(map[pixelMsg]struct{})
			var (
				width  int64
				height int64
			)

		stream:
			for {
				msg, ok := pixelsIn.Receive(ctx)
				if !ok {
					return
				}

				if runtime.IsStreamClose(msg.Msg) {
					break stream
				}
				if !runtime.IsStreamData(msg.Msg) {
					continue
				}

				var pix pixelMsg
				pix.decode(runtime.StreamDataValue(msg.Msg))
				if pix.x < 0 || pix.y < 0 {
					if !errOut.Send(ctx, errFromString("image.New: Pixel out of bounds")) {
						return
					}
				}
				if pix.x >= width {
					width = pix.x + 1
				}
				if pix.y >= height {
					height = pix.y + 1
				}
				pixels[pix] = struct{}{}
			}

			img := image.NewRGBA(image.Rect(0, 0, int(width), int(height)))
			for pix := range pixels {
				img.Set(int(pix.x), int(pix.y), pix.color.color())
			}

			if !imgOut.Send(ctx, runtime.NewStructMsg([]runtime.StructField{
				runtime.NewStructField("pixels", runtime.NewBytesMsg(img.Pix)),
				runtime.NewStructField("width", runtime.NewIntMsg(int64(img.Rect.Dx()))),
				runtime.NewStructField("height", runtime.NewIntMsg(int64(img.Rect.Dy()))),
			})) {
				return
			}
		}
	}, nil
}
