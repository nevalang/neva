package funcs

import (
	"context"
	"image"

	"github.com/nevalang/neva/internal/runtime"
)

type imageNew struct{}

func (imageNew) Create(io runtime.IO, _ runtime.Msg) (func(ctx context.Context), error) {
	pixelsIn, err := io.In.Single("pixels")
	if err != nil {
		return nil, err
	}

	imgOut, err := io.Out.Single("img")
	if err != nil {
		return nil, err
	}

	errOut, err := io.Out.Single("err")
	if err != nil {
		return nil, err
	}

	return func(ctx context.Context) {
		for {
			im := make(map[pixelMsg]struct{})
			var (
				width  int64
				height int64
			)
		stream:
			for {
				m, ok := pixelsIn.Receive(ctx)
				if !ok {
					return
				}

				var pix pixelStreamMsg
				pix.decode(m)
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
				im[pix.pixelMsg] = struct{}{}
				if pix.last {
					break stream
				}
			}

			img := image.NewRGBA(image.Rect(0, 0, int(width), int(height)))
			for pix := range im {
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
