package funcs

import (
	"context"
	"io"
	"net/http"

	"github.com/nevalang/neva/internal/runtime"
	"github.com/nevalang/neva/internal/runtime/messages"
)

type httpGet struct{}

//nolint:cyclop,gocognit,gocyclo // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
func (httpGet) Create(funcIO runtime.IO, _ messages.Msg) (func(ctx context.Context), error) {
	urlIn, err := funcIO.In.Single("url")
	if err != nil {
		//nolint:wrapcheck // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
		return nil, err
	}

	resOut, err := funcIO.Out.Single("res")
	if err != nil {
		//nolint:wrapcheck // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
		return nil, err
	}

	errOut, err := funcIO.Out.Single("err")
	if err != nil {
		//nolint:wrapcheck // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
		return nil, err
	}

	return func(ctx context.Context) {
		for {
			urlMsg, ok := urlIn.Receive(ctx)
			if !ok {
				return
			}

			//nolint:noctx // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
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

func respMsg(statusCode int, body []byte) messages.StructMsg {
	return messages.NewStructMsg([]messages.StructField{
		messages.NewStructField("body", messages.NewBytesMsg(body)),
		messages.NewStructField("statusCode", messages.NewIntMsg(int64(statusCode))),
	})
}
