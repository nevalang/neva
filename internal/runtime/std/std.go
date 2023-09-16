package std

import (
	"context"
	"fmt"

	"github.com/nevalang/neva/internal/runtime"
)

func Print(io runtime.FuncIO) (func(context.Context), error) {
	ch, err := io.In.Port("v")
	if err != nil {
		return nil, err
	}
	return func(ctx context.Context) {
		for {
			select {
			case <-ctx.Done():
				return
			case v := <-ch:
				fmt.Println(v.String())
			}
		}
	}, nil
}

func Lock(io runtime.FuncIO) (func(context.Context), error) {
	vin, err := io.In.Port("v")
	if err != nil {
		return nil, err
	}
	sig, err := io.In.Port("sig")
	if err != nil {
		return nil, err
	}
	vout, err := io.Out.Port("v")
	if err != nil {
		return nil, err
	}
	return func(ctx context.Context) {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				<-sig
				vout <- <-vin
			}
		}
	}, nil
}

func Const(io runtime.FuncIO) (func(context.Context), error) {
	vout, err := io.Out.Port("v")
	if err != nil {
		return nil, err
	}
	return func(ctx context.Context) {
		msg := ctx.Value("msg").(runtime.Msg)
		for {
			select {
			case <-ctx.Done():
				return
			default:
				vout <- msg
			}
		}
	}, nil
}
