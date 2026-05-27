package funcs

import (
	"context"

	"github.com/nevalang/neva/internal/runtime"
)

type listAt struct{}

// Create builds runtime closure for list item access by index.
//
// Invariants:
//  1. `idx` must be an integer message.
//  2. Positive and negative indexing are supported (`-1` means last element).
//  3. Out-of-bounds indexes are returned via `err` outport.
//  4. Typed scalar lists are handled first (int/string/bool/float) to avoid
//     generic materialization in hot paths.
//  5. Generic list fallback is used for non-scalar or mixed-value lists.
//
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
			switch handled, sendOK := sendTypedListAt(ctx, list, idx, resOut, errOut); {
			case handled && !sendOK:
				return
			case handled:
				continue
			default:
			}

			data := listToMsgs(list)

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

func sendTypedListAt(
	ctx context.Context,
	list runtime.ListMsg,
	idx int64,
	resOut runtime.SingleOutport,
	errOut runtime.SingleOutport,
) (bool, bool) {
	if data, ok := runtime.AsListInts(list); ok {
		return true, sendTypedListAtValue(ctx, data, idx, func(v int64) runtime.Msg { return runtime.NewIntMsg(v) }, resOut, errOut)
	}
	if data, ok := runtime.AsListStrings(list); ok {
		return true, sendTypedListAtValue(ctx, data, idx, func(v string) runtime.Msg { return runtime.NewStringMsg(v) }, resOut, errOut)
	}
	if data, ok := runtime.AsListBools(list); ok {
		return true, sendTypedListAtValue(ctx, data, idx, func(v bool) runtime.Msg { return runtime.NewBoolMsg(v) }, resOut, errOut)
	}
	if data, ok := runtime.AsListFloats(list); ok {
		return true, sendTypedListAtValue(ctx, data, idx, func(v float64) runtime.Msg { return runtime.NewFloatMsg(v) }, resOut, errOut)
	}
	return false, true
}

//nolint:ireturn // Generic helper converts scalar to runtime.Msg via constructor.
func sendTypedListAtValue[T any](
	ctx context.Context,
	data []T,
	idx int64,
	toMsg func(T) runtime.Msg,
	resOut runtime.SingleOutport,
	errOut runtime.SingleOutport,
) bool {
	item, valid := listItem(data, idx)
	if !valid {
		return errOut.Send(ctx, errFromString("index out of bounds"))
	}
	return resOut.Send(ctx, toMsg(item))
}
