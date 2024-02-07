package funcs

import (
	"bufio"
	"context"
	"os"

	"github.com/nevalang/neva/internal/runtime"
)

type reader struct{}

func (r reader) Create(io runtime.FuncIO, _ runtime.Msg) (func(ctx context.Context), error) {
	sig, err := io.In.Port("sig")
	if err != nil {
		return nil, err
	}

	vout, err := io.Out.Port("v")
	if err != nil {
		return nil, err
	}

	return func(ctx context.Context) {
		reader := bufio.NewReader(os.Stdin)
		for {
			select {
			case <-ctx.Done():
				return
			case <-sig:
				select {
				case <-ctx.Done():
					return
				default:
					bb, _, err := reader.ReadLine()
					if err != nil {
						panic(err)
					}
					select {
					case <-ctx.Done():
						return
					case vout <- runtime.NewStrMsg(string(bb)):
					}
				}
			}
		}
	}, nil
}
