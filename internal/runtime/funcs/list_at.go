package funcs

import (
	"context"

	"github.com/nevalang/neva/internal/runtime"
)

type listAt struct{}

//nolint:cyclop,gocognit,gocyclo,varnamelen,funlen,nestif // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
func (listAt) Create(io runtime.IO, _ runtime.Msg) (func(ctx context.Context), error) {
	dataIn, err := io.In.Single("data")
	if err != nil {
		//nolint:wrapcheck // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
		return nil, err
	}

	idxIn, err := io.In.Single("idx")
	if err != nil {
		//nolint:wrapcheck // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
		return nil, err
	}

	resOut, err := io.Out.Single("res")
	if err != nil {
		//nolint:wrapcheck // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
		return nil, err
	}

	errOut, err := io.Out.Single("err")
	if err != nil {
		//nolint:wrapcheck // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
		return nil, err
	}

	return func(ctx context.Context) {
		for {
			dataMsg, idxMsg, ok := receive2(ctx, dataIn, idxIn)
			if !ok {
				return
			}

			idx := idxMsg.Int()
			list := dataMsg.List()
			if data, ok := runtime.AsListInts(list); ok {
				item, valid := listItem(data, idx)
				if !valid {
					if !errOut.Send(ctx, errFromString("index out of bounds")) {
						return
					}
					continue
				}
				if !resOut.Send(ctx, runtime.NewIntMsg(item)) {
					return
				}
				continue
			}
			if data, ok := runtime.AsListStrings(list); ok {
				item, valid := listItem(data, idx)
				if !valid {
					if !errOut.Send(ctx, errFromString("index out of bounds")) {
						return
					}
					continue
				}
				if !resOut.Send(ctx, runtime.NewStringMsg(item)) {
					return
				}
				continue
			}
			if data, ok := runtime.AsListBools(list); ok {
				item, valid := listItem(data, idx)
				if !valid {
					if !errOut.Send(ctx, errFromString("index out of bounds")) {
						return
					}
					continue
				}
				if !resOut.Send(ctx, runtime.NewBoolMsg(item)) {
					return
				}
				continue
			}
			if data, ok := runtime.AsListFloats(list); ok {
				item, valid := listItem(data, idx)
				if !valid {
					if !errOut.Send(ctx, errFromString("index out of bounds")) {
						return
					}
					continue
				}
				if !resOut.Send(ctx, runtime.NewFloatMsg(item)) {
					return
				}
				continue
			}

			data := list.Msgs()

			l := int64(len(data))
			if idx < -l || idx >= l {
				if !errOut.Send(ctx, errFromString("index out of bounds")) {
					return
				}
			}

			if idx < 0 {
				// support negative indexing:
				//	$l = [1, 2, 3]
				//	$l[-1] // 3
				idx += int64(len(data))
			}

			if !resOut.Send(ctx, data[idx]) {
				return
			}
		}
	}, nil
}

//nolint:ireturn // Generic helper returns scalar value type parameter.
func listItem[T any](data []T, idx int64) (T, bool) {
	dataLen := int64(len(data))
	if idx < -dataLen || idx >= dataLen {
		var zero T
		return zero, false
	}
	if idx < 0 {
		idx += dataLen
	}
	return data[idx], true
}
