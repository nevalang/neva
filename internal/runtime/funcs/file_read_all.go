package funcs

import (
	"context"
	"io"
	"os"

	"github.com/nevalang/neva/internal/runtime"
)

type fileReadAll struct{}

//nolint:cyclop,gocognit,gocyclo // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
func (c fileReadAll) Create(rio runtime.IO, _ runtime.Msg) (func(ctx context.Context), error) {
	filenameIn, err := rio.In.Single("filename")
	if err != nil {
		//nolint:wrapcheck // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
		return nil, err
	}

	resOut, err := rio.Out.Single("res")
	if err != nil {
		//nolint:wrapcheck // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
		return nil, err
	}

	errOut, err := rio.Out.Single("err")
	if err != nil {
		//nolint:wrapcheck // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
		return nil, err
	}

	return func(ctx context.Context) {
		for {
			name, ok := filenameIn.Receive(ctx)
			if !ok {
				return
			}

			//nolint:varnamelen // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
			f, err := os.Open(name.Str())
			if err != nil {
				if !errOut.Send(ctx, errFromErr(err)) {
					return
				}
				continue
			}

			data, err := io.ReadAll(f)
			if err != nil {
				if !errOut.Send(ctx, errFromErr(err)) {
					return
				}
				continue
			}

			if err := f.Close(); err != nil {
				if !errOut.Send(ctx, errFromErr(err)) {
					return
				}
				continue
			}

			if !resOut.Send(ctx, runtime.NewBytesMsg(data)) {
				return
			}
		}
	}, nil
}
