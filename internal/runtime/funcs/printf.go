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
			tpl, ok := tplIn.Receive(ctx)
			if !ok {
				return
			}

			args := make([]runtime.Msg, argsIn.Len())
			if !argsIn.ReceiveAll(ctx, func(idx int, msg runtime.Msg) bool {
				args[idx] = msg
				return true
			}) {
				return
			}

			res, err := format(tpl.Str(), args)
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
	// Use a map to keep track of which arguments have been used
	usedArgs := make(map[int]bool)

	// Builder to construct the final result
	var result strings.Builder
	result.Grow(len(tpl)) // Optimistically assume no increase in length

	// Scan through the template to find and replace placeholders
	//nolint:varnamelen // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
	i := 0
	for i < len(tpl) {
		//nolint:nestif // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
		if tpl[i] == '$' {
			// Attempt to read an argument index after the '$'
			//nolint:varnamelen // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
			j := i + 1
			var argIndexStr strings.Builder
			for j < len(tpl) && tpl[j] >= '0' && tpl[j] <= '9' {
				argIndexStr.WriteString(string(tpl[j]))
				j++
			}

			if argIndexStr.String() != "" {
				argIndex, err := strconv.Atoi(argIndexStr.String())
				if err != nil {
					// Handle the error if the conversion fails
					return "", fmt.Errorf("invalid placeholder %s: %w", argIndexStr.String(), err)
				}

				if argIndex < 0 || argIndex >= len(args) {
					return "", fmt.Errorf("template refers to arg %d, but only %d args given", argIndex, len(args))
				}

				// Mark this arg as used
				usedArgs[argIndex] = true

				// Replace the placeholder with the argument's string representation
				fmt.Fprint(&result, args[argIndex])

				// Move past the current placeholder in the template
				i = j
				continue
			}
		}

		// If not processing a placeholder, just copy the current character
		result.WriteByte(tpl[i])
		i++
	}

	// Check if all arguments were used
	if len(usedArgs) != len(args) {
		return "", fmt.Errorf(
			"not all arguments are used in the template: got %v, used %v",
			len(args), len(usedArgs),
		)
	}

	return result.String(), nil
}
