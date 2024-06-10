package funcs

import (
	"context"
	"os"

	"github.com/nevalang/neva/internal/runtime"
)

type args struct{}

func (a args) Create(io runtime.FuncIO, _ runtime.Msg) (func(ctx context.Context), error) {
	outport, err := io.Out.Port("data")
	if err != nil {
		return nil, err
	}

	return func(ctx context.Context) {
		var lst_args []runtime.Msg
		for i := 0; i < len(os.Args); i++ {
			lst_args = append(lst_args, runtime.NewStrMsg(os.Args[i]))
		}
		for {
			select {

			case <-ctx.Done():
				return
			case outport <- runtime.NewListMsg(lst_args...):
			}
		}
	}, nil
}
