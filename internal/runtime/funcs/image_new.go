package funcs

import (
	"context"
	"image"

	"github.com/nevalang/neva/internal/runtime"
)

type imageNew struct{}

//nolint:cyclop,gocognit,gocyclo,varnamelen // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
func (imageNew) Create(io runtime.IO, _ runtime.Msg) (func(ctx context.Context), error) {
	pixelsIn, err := singleIn(io, "pixels")
	if err != nil {
		//nolint:wrapcheck // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
		return nil, err
	}

	imgOut, err := singleOut(io, "img")
	if err != nil {
		//nolint:wrapcheck // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
		return nil, err
	}

	errOut, err := singleOut(io, "err")
	if err != nil {
		//nolint:wrapcheck // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
		return nil, err
	}

	return func(ctx context.Context) {
		for {
			//nolint:varnamelen // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
			im := make(map[pixelMsg]struct{})
			var (
				width  int64
				height int64
			)
		stream:
			for {
				//nolint:varnamelen // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
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
