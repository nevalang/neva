package funcs

import (
	"context"

	"github.com/nevalang/neva/internal/runtime"
)

// streamSplitController serializes Split predicate work and preserves stream
// termination for each non-empty output branch.
type streamSplitController struct{}

type pendingStreamItem struct {
	data      runtime.Msg
	source    runtime.OrderedMsg
	predicate runtime.OrderedMsg
	idx       int64
	ok        bool
}

//nolint:cyclop,gocognit,gocyclo // The controller coordinates two branches and stream termination.
func (streamSplitController) Create(input runtime.IO, _ runtime.Msg) (func(ctx context.Context), error) {
	dataIn, err := input.In.Single("data")
	if err != nil {
		//nolint:wrapcheck // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
		return nil, err
	}

	predicateIn, err := input.In.Single("predicate")
	if err != nil {
		//nolint:wrapcheck // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
		return nil, err
	}

	itemOut, err := input.Out.Single("item")
	if err != nil {
		//nolint:wrapcheck // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
		return nil, err
	}

	thenOut, err := input.Out.Single("then")
	if err != nil {
		//nolint:wrapcheck // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
		return nil, err
	}

	elseOut, err := input.Out.Single("else")
	if err != nil {
		//nolint:wrapcheck // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
		return nil, err
	}

	return func(ctx context.Context) {
		var thenPending pendingStreamItem
		var elsePending pendingStreamItem

		for {
			msg, dataOK := dataIn.Receive(ctx)
			if !dataOK {
				return
			}

			item := msg.Struct()
			if !itemOut.Send(ctx, item.Get("data"), msg) {
				return
			}

			predicateMsg, predicateOK := predicateIn.Receive(ctx)
			if !predicateOK {
				return
			}

			targetOut := thenOut
			targetPending := &thenPending
			if !predicateMsg.Bool() {
				targetOut = elseOut
				targetPending = &elsePending
			}

			if !flushPendingStreamItem(ctx, targetOut, targetPending, false) {
				return
			}
			*targetPending = pendingStreamItem{
				data:      item.Get("data"),
				source:    msg,
				predicate: predicateMsg,
				idx:       item.Get("idx").Int(),
				ok:        true,
			}

			if !item.Get("last").Bool() {
				continue
			}

			if !flushPendingStreamItem(ctx, thenOut, &thenPending, true) {
				return
			}
			if !flushPendingStreamItem(ctx, elseOut, &elsePending, true) {
				return
			}
		}
	}, nil
}

func flushPendingStreamItem(
	ctx context.Context,
	out runtime.SingleOutport,
	pending *pendingStreamItem,
	last bool,
) bool {
	if !pending.ok {
		return true
	}

	if !out.Send(ctx, streamItem(pending.data, pending.idx, last), pending.source, pending.predicate) {
		return false
	}

	*pending = pendingStreamItem{}
	return true
}
