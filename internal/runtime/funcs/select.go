package funcs

import (
	"context"
	"errors"
	"sync"

	"github.com/nevalang/neva/internal/runtime"
	"github.com/nevalang/neva/internal/runtime/messages"
)

type selector struct{}

//nolint:varnamelen // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
func (selector) Create(io runtime.IO, _ messages.Msg) (func(ctx context.Context), error) {
	ifArrIn, err := io.In.Array("if")
	if err != nil {
		//nolint:wrapcheck // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
		return nil, err
	}

	thenArrIn, err := io.In.Array("then")
	if err != nil {
		//nolint:wrapcheck // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
		return nil, err
	}

	if ifArrIn.Len() != thenArrIn.Len() {
		return nil, errors.New("number of 'if' inports must match number of 'then' outports")
	}

	resOut, err := io.Out.Single("res")
	if err != nil {
		//nolint:wrapcheck // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
		return nil, err
	}

	return func(ctx context.Context) {
		for {
			var ifMsg runtime.SelectedMsg
			var ifOK, thenOK bool
			then := make([]messages.Msg, thenArrIn.Len())

			// The trigger and candidate values form one selection operation.
			// Receive them concurrently so either side may arrive first.
			var group sync.WaitGroup
			group.Go(func() {
				ifMsg, ifOK = ifArrIn.Select(ctx)
			})
			group.Go(func() {
				thenOK = thenArrIn.ReceiveAll(ctx, func(idx int, ordered runtime.OrderedMsg) bool {
					then[idx] = ordered.Msg
					return true
				})
			})
			group.Wait()

			if !ifOK || !thenOK {
				return
			}

			if !resOut.Send(ctx, then[ifMsg.SlotIdx], ifMsg.OrderedMsg) {
				return
			}
		}
	}, nil
}
