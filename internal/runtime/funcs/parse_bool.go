package funcs

import (
	"context"
	"strconv"
	"strings"

	"github.com/nevalang/neva/internal/runtime"
)

type parseBool struct{}

//nolint:gocognit // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
func (parseBool) Create(
	//nolint:varnamelen // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
	io runtime.IO,
	_ runtime.Msg,
) (func(ctx context.Context), error) {
	dataIn, err := singleIn(io, "data")
	if err != nil {
		//nolint:wrapcheck // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
		return nil, err
	}

	resOut, err := singleOut(io, "res")
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
			data, ok := dataIn.Receive(ctx)
			if !ok {
				return
			}

			parsed, parseErr := strconv.ParseBool(data.Str())
			if parseErr != nil {
				parseErrMsg := strings.TrimPrefix(parseErr.Error(), "strconv.ParseBool: ")
				if !errOut.Send(ctx, errFromString(parseErrMsg)) {
					return
				}
				continue
			}

			if !resOut.Send(ctx, runtime.NewBoolMsg(parsed)) {
				return
			}
		}
	}, nil
}
