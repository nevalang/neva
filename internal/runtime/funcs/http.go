package funcs

import (
	"context"
	goio "io"
	"net/http"

	"github.com/nevalang/neva/internal/runtime"
)

type httpGet struct{}

func (httpGet) Create(io runtime.FuncIO, _ runtime.Msg) (func(ctx context.Context), error) {
	urlIn, err := io.In.SingleInport("url")
	if err != nil {
		return nil, err
	}

	respOut, err := io.Out.SingleOutport("resp")
	if err != nil {
		return nil, err
	}

	errOut, err := io.Out.SingleOutport("err")
	if err != nil {
		return nil, err
	}

	return func(ctx context.Context) {
		for {
			var u string
			select {
			case m := <-urlIn:
				u = m.Str()
			case <-ctx.Done():
				return
			}
			resp, err := http.Get(u)
			if err != nil {
				select {
				case errOut <- runtime.NewMapMsg(map[string]runtime.Msg{
					"text": runtime.NewStrMsg(err.Error()),
				}):
					continue
				case <-ctx.Done():
					return
				}
			}
			body, err := goio.ReadAll(resp.Body)
			if err != nil {
				select {
				case errOut <- runtime.NewMapMsg(map[string]runtime.Msg{
					"text": runtime.NewStrMsg(err.Error()),
				}):
					continue
				case <-ctx.Done():
					return
				}
			}
			select {
			case respOut <- runtime.NewMapMsg(map[string]runtime.Msg{
				"statusCode": runtime.NewIntMsg(int64(resp.StatusCode)),
				"body":       runtime.NewStrMsg(string(body)),
			}):
			case <-ctx.Done():
				return
			}
		}
	}, nil
}
