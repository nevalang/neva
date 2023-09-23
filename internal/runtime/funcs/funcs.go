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
				select {
				case <-ctx.Done():
					return
				default:
					fmt.Println(v.String())
					select {
					case <-ctx.Done():
						return
					case vout <- v:
					}
				}
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
			case <-sig:
				select {
				case <-ctx.Done():
					return
				case v := <-vin:
					select {
					case <-ctx.Done():
						return
					case vout <- v:
					}
				}
			}
		}
	}, nil
}

func Const(ctx context.Context, io runtime.FuncIO) (func(), error) {
	msg := ctx.Value("msg")
	if msg == nil {
		return nil, errors.New("ctx msg not found")
	}

	v, ok := msg.(runtime.Msg)
	if !ok {
		return nil, errors.New("ctx value is not runtime message")
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
			case vout <- v:
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
