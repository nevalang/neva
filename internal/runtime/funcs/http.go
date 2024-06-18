package funcs

import (
	"context"
	"io"
	"net/http"

	"github.com/nevalang/neva/internal/runtime"
)

type httpGet struct{}

func (httpGet) Create(funcIO runtime.FuncIO, _ runtime.Msg) (func(ctx context.Context), error) {
	urlIn, err := funcIO.In.Single("url")
	if err != nil {
		return nil, err
	}

	respOut, err := funcIO.Out.Single("resp")
	if err != nil {
		return nil, err
	}

	errOut, err := funcIO.Out.Single("err")
	if err != nil {
		return nil, err
	}

	return func(ctx context.Context) {
		for {
			urlMsg, ok := urlIn.Receive(ctx)
			if !ok {
				return
			}

			resp, err := http.Get(urlMsg.Str())
			if err != nil {
				if !errOut.Send(ctx, errFromErr(err)) {
					return
				}
				continue
			}

			body, err := io.ReadAll(resp.Body)
			if err != nil {
				if !errOut.Send(ctx, errFromErr(err)) {
					return
				}
				continue
			}

			if !respOut.Send(
				ctx,
				respMsg(resp.StatusCode, body),
			) {
				return
			}
		}
	}, nil
}

func respMsg(statusCode int, body []byte) runtime.MapMsg {
	return runtime.NewMapMsg(map[string]runtime.Msg{
		"body":       runtime.NewStrMsg(string(body)),
		"statusCode": runtime.NewIntMsg(int64(statusCode)),
	})
}
