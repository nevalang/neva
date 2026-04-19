package funcs

import (
	"context"
	"os"

	"github.com/nevalang/neva/internal/runtime"
)

type writeAll struct{}

//nolint:gocognit // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
func (c writeAll) Create(rio runtime.IO, _ runtime.Msg) (func(ctx context.Context), error) {
	filenameIn, err := singleIn(rio, "filename")
	if err != nil {
		//nolint:wrapcheck // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
		return nil, err
	}

	dataIn, err := singleIn(rio, "data")
	if err != nil {
		//nolint:wrapcheck // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
		return nil, err
	}

	resOut, err := singleOut(rio, "res")
	if err != nil {
		//nolint:wrapcheck // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
		return nil, err
	}

	errOut, err := singleOut(rio, "err")
	if err != nil {
		//nolint:wrapcheck // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
		return nil, err
	}

	return func(ctx context.Context) {
		for {
			//nolint:varnamelen // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
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
