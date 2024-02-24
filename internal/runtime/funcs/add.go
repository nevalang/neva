package funcs

import (
	"context"

	"github.com/nevalang/neva/internal/runtime"
)

type intAdder struct{}

func (intAdder) Create(_ runtime.FuncIO, _ runtime.Msg) (func(ctx context.Context), error) {
	// vvin, err := io.In.Port("vv")
	// if err != nil {
	// 	return nil, err
	// }

	// vout, err := io.Out.Port("v")
	// if err != nil {
	// 	return nil, err
	// }

	return func(ctx context.Context) {
		for range ctx.Done() {
			return
		}
	}, nil
}

// type addFloats struct{}

// func (a addFloats) Create(io runtime.FuncIO, _ runtime.Msg) (func(ctx context.Context), error) {
// 	vvin, ok := io.In["vv"]
// 	if !ok {
// 		return nil, errors.New("vv port not found")
// 	}

// 	vout, err := io.Out.Port("v")
// 	if err != nil {
// 		return nil, err
// 	}

// 	return func(ctx context.Context) {
// 		for {
// 			var sum float64
// 			for _, vin := range vvin {
// 				select {
// 				case <-ctx.Done():
// 					return
// 				case msg := <-vin:
// 					sum += msg.Float()
// 				}
// 			}
// 			select {
// 			case <-ctx.Done():
// 				return
// 			case vout <- runtime.NewFloatMsg(sum):
// 			}
// 		}
// 	}, nil
// }

// type addStrings struct{}

// func (a addStrings) Create(io runtime.FuncIO, _ runtime.Msg) (func(ctx context.Context), error) {
// 	vvin, ok := io.In["vv"]
// 	if !ok {
// 		return nil, errors.New("vv port not found")
// 	}

// 	vout, err := io.Out.Port("v")
// 	if err != nil {
// 		return nil, err
// 	}

// 	return func(ctx context.Context) {
// 		for {
// 			var sum string
// 			for _, vin := range vvin {
// 				select {
// 				case <-ctx.Done():
// 					return
// 				case msg := <-vin:
// 					sum += msg.Str()
// 				}
// 			}
// 			select {
// 			case <-ctx.Done():
// 				return
// 			case vout <- runtime.NewStrMsg(sum):
// 			}
// 		}
// 	}, nil
// }

// type addLists struct{}

// func (a addLists) Create(io runtime.FuncIO, _ runtime.Msg) (func(ctx context.Context), error) {
// 	vvin, ok := io.In["vv"]
// 	if !ok {
// 		return nil, errors.New("vv port not found")
// 	}

// 	vout, err := io.Out.Port("v")
// 	if err != nil {
// 		return nil, err
// 	}

// 	return func(ctx context.Context) {
// 		for {
// 			var sum []runtime.Msg
// 			for _, vin := range vvin {
// 				select {
// 				case <-ctx.Done():
// 					return
// 				case msg := <-vin:
// 					sum = append(sum, msg.List()...)
// 				}
// 			}
// 			select {
// 			case <-ctx.Done():
// 				return
// 			case vout <- runtime.NewStrMsg(sum):
// 			}
// 		}
// 	}, nil
// }
