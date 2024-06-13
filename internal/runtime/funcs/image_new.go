package funcs

import (
	"context"
	"image"

	"github.com/nevalang/neva/internal/runtime"
)

type imageNew struct{}

func (imageNew) Create(io runtime.FuncIO, _ runtime.Msg) (func(ctx context.Context), error) {
	pixelsIn, err := io.In.SingleInport("pixels")
	if err != nil {
		return nil, err
	}

	imgOut, err := io.Out.SingleOutport("img")
	if err != nil {
		return nil, err
	}

	errOut, err := io.Out.SingleOutport("err")
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
					errOut.Send(ctx, errFromErr())

					select {
					case errOut <- runtime.NewMapMsg(map[string]runtime.Msg{
						"text": runtime.NewStrMsg("image.New: Pixel out of bounds"),
					}):
					case <-ctx.Done():
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

			var i imageMsg
			i.decodeImage(img)
			if !imgOut.Send(ctx, i.encode()) {
				return
			}
		}
	}, nil
}
