package funcs

import (
	"context"
	"errors"

	"github.com/nevalang/neva/internal/runtime"
)

type matchSelector struct{}

//nolint:cyclop,gocognit,gocyclo,varnamelen // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
func (matchSelector) Create(io runtime.IO, _ runtime.Msg) (func(ctx context.Context), error) {
	dataIn, err := singleIn(io, "data")
	if err != nil {
		//nolint:wrapcheck // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
		return nil, err
	}

	ifIn, err := arrayIn(io, "if")
	if err != nil {
		//nolint:wrapcheck // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
		return nil, err
	}

	thenOut, err := arrayIn(io, "then")
	if err != nil {
		//nolint:wrapcheck // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
		return nil, err
	}

	if ifIn.Len() != thenOut.Len() {
		return nil, errors.New("number of 'if' inports must match number of 'then' outports")
	}

	elseIn, err := singleIn(io, "else")
	if err != nil {
		//nolint:wrapcheck // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
		return nil, err
	}

	resOut, err := singleOut(io, "res")
	if err != nil {
		//nolint:wrapcheck // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
		return nil, err
	}

	return func(ctx context.Context) {
		for {
			//nolint:varnamelen // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
			dataMsg, ok := dataIn.Receive(ctx)
			if !ok {
				return
			}

			ifMsgs := make([]runtime.Msg, ifIn.Len())
			if !ifIn.ReceiveAll(ctx, func(idx int, msg runtime.Msg) bool {
				ifMsgs[idx] = msg
				return true
			}) {
				return
			}

			thenMsgs := make([]runtime.Msg, thenOut.Len())
			if !thenOut.ReceiveAll(ctx, func(idx int, msg runtime.Msg) bool {
				thenMsgs[idx] = msg
				return true
			}) {
				return
			}

			elseInMsg, ok := elseIn.Receive(ctx)
			if !ok {
				return
			}

			resMsg := elseInMsg
			for i, ifMsg := range ifMsgs {
				if runtime.Match(dataMsg, ifMsg) {
					resMsg = thenMsgs[i]
					break
				}
			}

			if u, ok := resMsg.(runtime.UnionMsg); ok {
				resMsg = u.Data()
			}

			if !resOut.Send(ctx, resMsg) {
				return
			}
		}
	}, nil
}
