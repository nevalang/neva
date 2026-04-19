package funcs

import (
	"context"
	"unicode"

	"github.com/nevalang/neva/internal/runtime"
)

type stringFromIntCodepoint struct{}

func (stringFromIntCodepoint) Create(
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

	return func(ctx context.Context) {
		for {
			data, ok := dataIn.Receive(ctx)
			if !ok {
				return
			}

			res := codePointString(data.Int())
			if !resOut.Send(ctx, runtime.NewStringMsg(res)) {
				return
			}
		}
	}, nil
}

//nolint:varnamelen // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
func codePointString(v int64) string {
	if v < 0 || v > unicode.MaxRune || (v >= 0xD800 && v <= 0xDFFF) {
		return string(unicode.ReplacementChar)
	}

	// #nosec G115 -- guarded by unicode range checks above
	return string(rune(v))
}
