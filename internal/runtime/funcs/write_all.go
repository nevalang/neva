package funcs

import (
	"context"
	"os"

	"github.com/nevalang/neva/internal/runtime"
)

type writeAll struct{}

//nolint:gocognit
func (c writeAll) Create(rio runtime.IO, _ runtime.Msg) (func(ctx context.Context), error) {
	filenameIn, err := rio.In.Single("filename")
	if err != nil {
		//nolint:wrapcheck
		return nil, err
	}

	dataIn, err := rio.In.Single("data")
	if err != nil {
		//nolint:wrapcheck
		return nil, err
	}

	resOut, err := rio.Out.Single("res")
	if err != nil {
		//nolint:wrapcheck
		return nil, err
	}

	errOut, err := rio.Out.Single("err")
	if err != nil {
		//nolint:wrapcheck
		return nil, err
	}

	return func(ctx context.Context) {
		for {
			//nolint:varnamelen
			filenameMsg, ok := filenameIn.Receive(ctx)
			if !ok {
				return
			}

			dataMsg, ok := dataIn.Receive(ctx)
			if !ok {
				return
			}

			// #nosec G306 -- file is user-controlled output
			err := os.WriteFile(filenameMsg.Str(), dataMsg.Bytes(), 0o644)
			if err != nil {
				if !errOut.Send(ctx, errFromErr(err)) {
					return
				}
				continue
			}

			if !resOut.Send(ctx, emptyStruct()) {
				return
			}
		}
	}, nil
}
