package funcs

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"os"
	"strconv"

	"github.com/nevalang/neva/internal/runtime"
)

func Read(ctx context.Context, io runtime.FuncIO) (func(), error) {
	sig, err := io.In.Port("sig")
	if err != nil {
		return nil, err
	}
	vout, err := io.Out.Port("v")
	if err != nil {
		return nil, err
	}
	return func() {
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
					text, err := reader.ReadString('\n')
					if err != nil {
						panic(err) // TODO handle
					}
					select {
					case <-ctx.Done():
						return
					case vout <- runtime.NewStrMsg(text):
					}
				}
			}
		}
	}, nil
}

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
					fmt.Print(v.String())
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

func Void(ctx context.Context, io runtime.FuncIO) (func(), error) {
	vin, err := io.In.Port("v")
	if err != nil {
		return nil, err
	}

	return func() {
		for {
			select {
			case <-ctx.Done():
				return
			case <-vin:
			}
		}
	}, nil
}

func Add(ctx context.Context, io runtime.FuncIO) (func(), error) {
	msg := ctx.Value("msg")
	if msg == nil {
		return nil, errors.New("ctx msg not found")
	}

	typ, ok := msg.(runtime.Msg)
	if !ok {
		return nil, errors.New("ctx value is not runtime message")
	}

	var handler func(a, b runtime.Msg) runtime.Msg
	switch typ.Type() {
	case runtime.IntMsgType:
		handler = func(a, b runtime.Msg) runtime.Msg {
			return runtime.NewIntMsg(a.Int() + b.Int())
		}
	case runtime.FloatMsgType:
		handler = func(a, b runtime.Msg) runtime.Msg {
			return runtime.NewFloatMsg(a.Float() + b.Float())
		}
	case runtime.StrMsgType:
		handler = func(a, b runtime.Msg) runtime.Msg {
			return runtime.NewStrMsg(a.Str() + b.Str())
		}
	default:
		return nil, errors.New("unknown msg type")
	}

	a, err := io.In.Port("a")
	if err != nil {
		return nil, err
	}
	b, err := io.In.Port("b")
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
			case v1 := <-a:
				select {
				case <-ctx.Done():
					return
				case v2 := <-b:
					select {
					case <-ctx.Done():
						return
					case vout <- handler(v1, v2):
					}
				}
			}
		}
	}, nil
}

func ParseInt(ctx context.Context, io runtime.FuncIO) (func(), error) {
	vin, err := io.In.Port("v")
	if err != nil {
		return nil, err
	}

	vout, err := io.Out.Port("v")
	if err != nil {
		return nil, err
	}

	errout, err := io.Out.Port("err")
	if err != nil {
		return nil, err
	}

	parseFunc := func(str string) (runtime.Msg, error) {
		v, err := strconv.Atoi(str)
		if err != nil {
			return nil, err
		}
		return runtime.NewIntMsg(int64(v)), nil
	}

	return func() {
		for {
			select {
			case <-ctx.Done():
				return
			case str := <-vin:
				v, err := parseFunc(str.Str())
				if err != nil {
					select {
					case <-ctx.Done():
						return
					case errout <- runtime.NewStrMsg(err.Error()):
					}
					continue
				}
				select {
				case <-ctx.Done():
					return
				case vout <- v:
				}
			}
		}
	}, nil
}

func Repo() map[string]runtime.Func {
	return map[string]runtime.Func{
		"Read":     Read,
		"Print":    Print,
		"Lock":     Lock,
		"Const":    Const,
		"Add":      Add,
		"ParseInt": ParseInt,
		"Void":     Void,
	}
}
