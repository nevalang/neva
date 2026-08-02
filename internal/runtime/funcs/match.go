package funcs

import (
	"context"
	"errors"
	"sync"

	"github.com/nevalang/neva/internal/runtime"
	"github.com/nevalang/neva/internal/runtime/messages"
)

type matchSelector struct{}

//nolint:cyclop,gocognit,gocyclo,varnamelen // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
func (matchSelector) Create(io runtime.IO, _ messages.Msg) (func(ctx context.Context), error) {
	dataIn, err := io.In.Single("data")
	if err != nil {
		//nolint:wrapcheck // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
		return nil, err
	}

	ifIn, err := io.In.Array("if")
	if err != nil {
		//nolint:wrapcheck // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
		return nil, err
	}

	thenOut, err := io.In.Array("then")
	if err != nil {
		//nolint:wrapcheck // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
		return nil, err
	}

	if ifIn.Len() != thenOut.Len() {
		return nil, errors.New("number of 'if' inports must match number of 'then' outports")
	}

	elseIn, err := io.In.Single("else")
	if err != nil {
		//nolint:wrapcheck // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
		return nil, err
	}

	resOut, err := io.Out.Single("res")
	if err != nil {
		//nolint:wrapcheck // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
		return nil, err
	}

	return func(ctx context.Context) {
		for {
			var dataMsg, elseInMsg runtime.OrderedMsg
			var dataOK, ifOK, thenOK, elseOK bool
			ifMsgs := make([]runtime.OrderedMsg, ifIn.Len())
			thenMsgs := make([]runtime.OrderedMsg, thenOut.Len())

			// All inputs select one result. Receive them concurrently so their
			// producers may send in any order.
			var group sync.WaitGroup
			group.Go(func() {
				dataMsg, dataOK = dataIn.Receive(ctx)
			})
			group.Go(func() {
				ifOK = ifIn.ReceiveAll(ctx, func(idx int, ordered runtime.OrderedMsg) bool {
					ifMsgs[idx] = ordered
					return true
				})
			})
			group.Go(func() {
				thenOK = thenOut.ReceiveAll(ctx, func(idx int, ordered runtime.OrderedMsg) bool {
					thenMsgs[idx] = ordered
					return true
				})
			})
			group.Go(func() {
				elseInMsg, elseOK = elseIn.Receive(ctx)
			})
			group.Wait()

			if !dataOK || !ifOK || !thenOK || !elseOK {
				return
			}

			resMsg := elseInMsg.Msg
			causes := []runtime.OrderedMsg{dataMsg, elseInMsg}
			for i, ifMsg := range ifMsgs {
				if messages.Match(dataMsg.Msg, ifMsg.Msg) {
					resMsg = thenMsgs[i].Msg
					causes = []runtime.OrderedMsg{dataMsg, ifMsg, thenMsgs[i]}
					break
				}
			}

			resMsg = tryToUnboxIfUnion(resMsg)

			if !resOut.Send(ctx, resMsg, causes...) {
				return
			}
		}
	}, nil
}
