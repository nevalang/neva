package funcs

import (
	"context"
	"strings"

	"github.com/nevalang/neva/internal/runtime"
)

type stringSplit struct{}

func (p stringSplit) Create(io runtime.FuncIO, _ runtime.Msg) (func(ctx context.Context), error) {
	dataIn, err := io.In.Port("data")
	if err != nil {
		return nil, err
	}

	delimIn, err := io.In.Port("delim")
	if err != nil {
		return nil, err
	}

	resOut, err := io.Out.Port("res")
	if err != nil {
		return nil, err
	}

	return func(ctx context.Context) {
		var (
			data  runtime.Msg
			delim runtime.Msg
		)

		for {
			select {
			case <-ctx.Done():
				return
			case data = <-dataIn:
			}

			select {
			case <-ctx.Done():
				return
			case delim = <-delimIn:
			}

			splitted := strings.Split(data.Str(), delim.Str())
			res := make([]runtime.Msg, len(splitted))
			for i, s := range splitted {
				res[i] = runtime.NewStrMsg(s)
			}

			select {
			case <-ctx.Done():
				return
			case resOut <- runtime.NewListMsg(res...):
			}
		}
	}, nil
}
