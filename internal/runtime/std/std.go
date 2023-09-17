package std

import (
	"context"
	"errors"
	"fmt"

	"github.com/nevalang/neva/internal/runtime"
)

func Print(ctx context.Context, io runtime.FuncIO) (func(), error) {
	in, err := io.In.Port("v")
	if err != nil {
		return nil, err
	}
	out, err := io.Out.Port("v")
	if err != nil {
		return nil, err
	}
	return func() {
		for {
			select {
			case <-ctx.Done():
				return
			case v := <-in:
				fmt.Println(v.String())
				out <- v
			}
		}
	}, nil
}

func Lock(ctx context.Context, io runtime.FuncIO) (func(), error) {
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
	return func() {
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

func Const(ctx context.Context, io runtime.FuncIO) (func(), error) {
	msg, ok := ctx.Value("msg").(runtime.Msg)
	if !ok {
		return nil, errors.New("ctx msg not found")
	}

	out, err := io.Out.Port("v")
	if err != nil {
		return nil, err
	}

	return func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				out <- msg
			}
		}
	}, nil
}
