package funcs

import (
	"context"
	"errors"
	"fmt"

	"github.com/nevalang/neva/internal/runtime"
)

func Print(ctx context.Context, io runtime.FuncIO) (func(), error) {
	vin, err := io.In.Port("v")
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
			case v := <-vin:
				fmt.Println(v.String())
				vout <- v
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
				fmt.Println(ctx.Err())
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

func Repo() map[string]runtime.Func {
	return map[string]runtime.Func{
		"Print": Print,
		"Lock":  Lock,
		"Const": Const,
	}
}
