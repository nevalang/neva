package funcs

import (
	"context"
	"io"
	"net/http"

	"github.com/nevalang/neva/internal/runtime"
)

type httpGet struct{}

func (httpGet) Create(funcIO runtime.IO, _ runtime.Msg) (func(ctx context.Context), error) {
	urlIn, err := funcIO.In.Single("url")
	if err != nil {
		return nil, err
	}

	resOut, err := funcIO.Out.Single("res")
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
			closeErr := resp.Body.Close()
			if err != nil {
				if !errOut.Send(ctx, errFromErr(err)) {
					return
				}
				continue
			}
			if closeErr != nil {
				if !errOut.Send(ctx, errFromErr(closeErr)) {
					return
				}
				continue
			}

			if !resOut.Send(
				ctx,
				respMsg(resp.StatusCode, body),
			) {
				return
			}
		}
	}, nil
}

func respMsg(statusCode int, body []byte) runtime.StructMsg {
	return runtime.NewStructMsg([]runtime.StructField{
		runtime.NewStructField("body", runtime.NewBytesMsg(body)),
		runtime.NewStructField("statusCode", runtime.NewIntMsg(int64(statusCode))),
	})
}
