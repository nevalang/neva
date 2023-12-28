package funcs

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/nevalang/neva/internal/runtime"
)

func read(io runtime.FuncIO, _ runtime.Msg) (func(ctx context.Context), error) {
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

func print(io runtime.FuncIO, _ runtime.Msg) (func(ctx context.Context), error) { //nolint:predeclared
	vin, err := io.In.Port("v")
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

func lock(io runtime.FuncIO, _ runtime.Msg) (func(ctx context.Context), error) {
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

func constant(io runtime.FuncIO, msg runtime.Msg) (func(ctx context.Context), error) {
	v, ok := msg.(runtime.Msg)
	if !ok {
		return nil, errors.New("ctx value is not runtime message")
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
			case vout <- v:
			}
		}
	}, nil
}

func void(io runtime.FuncIO, _ runtime.Msg) (func(ctx context.Context), error) {
	vin, err := io.In.Port("v")
	if err != nil {
		return nil, err
	}

	return func(ctx context.Context) {
		for {
			select {
			case <-ctx.Done():
				return
			case <-vin:
			}
		}
	}, nil
}

func add(io runtime.FuncIO, _ runtime.Msg) (func(ctx context.Context), error) {
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

	handler := func(a, b runtime.Msg) runtime.Msg {
		return runtime.NewIntMsg(a.Int() + b.Int())
	}

	return func(ctx context.Context) {
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

func parseInt(io runtime.FuncIO, _ runtime.Msg) (func(ctx context.Context), error) {
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
			return nil, errors.New(strings.TrimPrefix(err.Error(), "strconv.Atoi: "))
		}
		return runtime.NewIntMsg(int64(v)), nil
	}

	return func(ctx context.Context) {
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

func Registry() map[string]runtime.Func {
	return map[string]runtime.Func{
		"Read":     read,
		"Print":    print,
		"Lock":     lock,
		"Const":    constant,
		"Add":      add,
		"ParseInt": parseInt,
		"Void":     void,
	}
}
