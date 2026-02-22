package funcs

import (
	"context"
	"unicode"

	"github.com/nevalang/neva/internal/runtime"
)

type stringFromIntCodepoint struct{}

func (stringFromIntCodepoint) Create(
	io runtime.IO,
	_ runtime.Msg,
) (func(ctx context.Context), error) {
	dataIn, err := io.In.Single("data")
	if err != nil {
		return nil, err
	}

	resOut, err := io.Out.Single("res")
	if err != nil {
		return nil, err
	}

	return func(ctx context.Context) {
		for {
			data, ok := dataIn.Receive(ctx)
			if !ok {
				return
			}

			res := codepointToString(data.Int())
			if !resOut.Send(ctx, runtime.NewStringMsg(res)) {
				return
			}
		}
	}, nil
}

func codepointToString(v int64) string {
	if v < 0 || v > unicode.MaxRune || (v >= 0xD800 && v <= 0xDFFF) {
		return string(unicode.ReplacementChar)
	}

	return string(rune(v))
}
