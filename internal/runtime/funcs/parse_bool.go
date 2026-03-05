package funcs

import (
	"context"
	"strconv"
	"strings"

	"github.com/nevalang/neva/internal/runtime"
)

type parseBool struct{}

func (parseBool) Create(
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

	errOut, err := io.Out.Single("err")
	if err != nil {
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
