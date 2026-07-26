package funcs

import (
	"context"
	"fmt"

	"github.com/nevalang/neva/internal/runtime"
	"github.com/nevalang/neva/internal/runtime/messages"
)

type printf struct{}

//nolint:varnamelen // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
func (p printf) Create(io runtime.IO, _ messages.Msg) (func(ctx context.Context), error) {
	tplIn, err := io.In.Single("tpl")
	if err != nil {
		//nolint:wrapcheck // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
		return nil, err
	}

	argsIn, err := io.In.Array("args")
	if err != nil {
		//nolint:perfsprint // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
		return nil, fmt.Errorf("missing required input port 'args'")
	}

	sigOut, err := io.Out.Single("sig")
	if err != nil {
		//nolint:perfsprint // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
		return nil, fmt.Errorf("missing required output port 'args'")
	}

	errOut, err := io.Out.Single("err")
	if err != nil {
		//nolint:wrapcheck // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
		return nil, err
	}

	return p.handle(tplIn, argsIn, errOut, sigOut)
}

//nolint:gocognit // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
func (printf) handle(
	tplIn runtime.SingleInport,
	//nolint:gocritic // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
	argsIn runtime.ArrayInport,
	errOut runtime.SingleOutport,
	sigOut runtime.SingleOutport,
) (func(ctx context.Context), error) {
	return func(ctx context.Context) {
		for {
			templateMsg, received := tplIn.Receive(ctx)
			if !received {
				return
			}

			args, causes, received := receivePrintfArgs(ctx, &argsIn)
			if !received {
				return
			}

			res, err := messages.StringFormat(templateMsg.Str(), args)
			if err != nil {
				if !errOut.Send(ctx, errFromErr(err), append(causes, templateMsg)...) {
					return
				}
				continue
			}

			if _, err := fmt.Print(res); err != nil {
				if !errOut.Send(ctx, errFromErr(err), append(causes, templateMsg)...) {
					return
				}
				continue
			}

			if !sigOut.Send(ctx, messages.NewStringMsg(res), append(causes, templateMsg)...) {
				return
			}
		}
	}, nil
}

func receivePrintfArgs(ctx context.Context, argsIn *runtime.ArrayInport) ([]messages.Msg, []runtime.OrderedMsg, bool) {
	args := make([]messages.Msg, argsIn.Len())
	causes := make([]runtime.OrderedMsg, argsIn.Len())
	if !argsIn.ReceiveAll(ctx, func(idx int, ordered runtime.OrderedMsg) bool {
		args[idx] = ordered.Msg
		causes[idx] = ordered
		return true
	}) {
		return nil, nil, false
	}
	return args, causes, true
}
