package funcs

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/nevalang/neva/internal/runtime"
)

type printf struct{}

//nolint:varnamelen // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
func (p printf) Create(io runtime.IO, _ runtime.Msg) (func(ctx context.Context), error) {
	tplIn, err := singleIn(io, "tpl")
	if err != nil {
		//nolint:wrapcheck // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
		return nil, err
	}

	argsIn, err := arrayIn(io, "args")
	if err != nil {
		//nolint:perfsprint // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
		return nil, fmt.Errorf("missing required input port 'args'")
	}

	sigOut, err := singleOut(io, "sig")
	if err != nil {
		//nolint:perfsprint // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
		return nil, fmt.Errorf("missing required output port 'args'")
	}

	errOut, err := singleOut(io, "err")
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

			args, received := receivePrintfArgs(ctx, &argsIn)
			if !received {
				return
			}

			res, err := format(templateMsg.Str(), args)
			if err != nil {
				if !errOut.Send(ctx, errFromErr(err)) {
					return
				}
				continue
			}

			if _, err := fmt.Print(res); err != nil {
				if !errOut.Send(ctx, errFromErr(err)) {
					return
				}
				continue
			}

			if !sigOut.Send(ctx, runtime.NewStringMsg(res)) {
				return
			}
		}
	}, nil
}

//nolint:gocognit // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
func format(tpl string, args []runtime.Msg) (string, error) {
	usedArgs := make(map[int]bool)
	var result strings.Builder
	result.Grow(len(tpl))

	//nolint:varnamelen // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
	templatePos := 0
	for templatePos < len(tpl) {
		placeholderIndex, nextIndex, hasPlaceholder, err := parsePlaceholderAt(tpl, templatePos)
		if err != nil {
			return "", err
		}
		//nolint:nestif // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
		if hasPlaceholder {
			if placeholderIndex >= len(args) {
				return "", fmt.Errorf("template refers to arg %d, but only %d args given", placeholderIndex, len(args))
			}
			usedArgs[placeholderIndex] = true
			fmt.Fprint(&result, args[placeholderIndex])
			templatePos = nextIndex
			continue
		}

		result.WriteByte(tpl[templatePos])
		templatePos++
	}

	if len(usedArgs) != len(args) {
		return "", fmt.Errorf(
			"not all arguments are used in the template: got %v, used %v",
			len(args), len(usedArgs),
		)
	}

	return result.String(), nil
}

func receivePrintfArgs(ctx context.Context, argsIn *runtime.ArrayInport) ([]runtime.Msg, bool) {
	args := make([]runtime.Msg, argsIn.Len())
	if !argsIn.ReceiveAll(ctx, func(idx int, msg runtime.Msg) bool {
		args[idx] = msg
		return true
	}) {
		return nil, false
	}
	return args, true
}

func parsePlaceholderAt(template string, startIdx int) (int, int, bool, error) {
	if startIdx >= len(template) || template[startIdx] != '$' {
		return 0, startIdx, false, nil
	}

	//nolint:varnamelen // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
	digitStart := startIdx + 1
	digitEnd := digitStart
	for digitEnd < len(template) && template[digitEnd] >= '0' && template[digitEnd] <= '9' {
		digitEnd++
	}

	if digitStart == digitEnd {
		return 0, startIdx, false, nil
	}

	argIndex, err := strconv.Atoi(template[digitStart:digitEnd])
	if err != nil {
		return 0, 0, false, fmt.Errorf("invalid placeholder %q: %w", template[digitStart:digitEnd], err)
	}

	return argIndex, digitEnd, true, nil
}
