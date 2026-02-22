package funcs

import (
	"context"
	"strconv"

	"github.com/nevalang/neva/internal/runtime"
)

type formatFloat struct{}

func (formatFloat) Create(
	io runtime.IO,
	_ runtime.Msg,
) (func(ctx context.Context), error) {
	dataIn, err := io.In.Single("data")
	if err != nil {
		return nil, err
	}

	fmtIn, err := io.In.Single("fmt")
	if err != nil {
		return nil, err
	}

	precIn, err := io.In.Single("prec")
	if err != nil {
		return nil, err
	}

	bitsIn, err := io.In.Single("bits")
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

			fmtMsg, ok := fmtIn.Receive(ctx)
			if !ok {
				return
			}

			prec, ok := precIn.Receive(ctx)
			if !ok {
				return
			}

			bits, ok := bitsIn.Receive(ctx)
			if !ok {
				return
			}

			format := byte('g')
			formatStr := fmtMsg.Str()
			if len(formatStr) > 0 {
				format = formatStr[0]
			}

			res := strconv.FormatFloat(
				data.Float(),
				format,
				int(prec.Int()),
				int(bits.Int()),
			)
			if !resOut.Send(ctx, runtime.NewStringMsg(res)) {
				return
			}
		}
	}, nil
}
