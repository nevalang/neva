package funcs

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/nevalang/neva/internal/runtime"
)

type printf struct{}

func (p printf) Create(io runtime.IO, _ runtime.Msg) (func(ctx context.Context), error) {
	tplIn, err := io.In.Single("tpl")
	if err != nil {
		return nil, err
	}

	argsIn, err := io.In.Array("args")
	if err != nil {
		return nil, fmt.Errorf("missing required input port 'args'")
	}

	sigOut, err := io.Out.Single("sig")
	if err != nil {
		return nil, fmt.Errorf("missing required output port 'args'")
	}

	errOut, err := io.Out.Single("err")
	if err != nil {
		return nil, err
	}

	return p.handle(tplIn, argsIn, errOut, sigOut)
}

func (printf) handle(
	tplIn runtime.SingleInport,
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

func format(tpl string, args []runtime.Msg) (string, error) {
	// Use a map to keep track of which arguments have been used
	usedArgs := make(map[int]bool)

	// Builder to construct the final result
	var result strings.Builder
	result.Grow(len(tpl)) // Optimistically assume no increase in length

	// Scan through the template to find and replace placeholders
	i := 0
	for i < len(tpl) {
		if tpl[i] == '$' {
			// Attempt to read an argument index after the '$'
			j := i + 1
			argIndexStr := ""
			for j < len(tpl) && tpl[j] >= '0' && tpl[j] <= '9' {
				argIndexStr += string(tpl[j])
				j++
			}

			if argIndexStr != "" {
				argIndex, err := strconv.Atoi(argIndexStr)
				if err != nil {
					// Handle the error if the conversion fails
					return "", fmt.Errorf("invalid placeholder %s: %v", argIndexStr, err)
				}

				if argIndex < 0 || argIndex >= len(args) {
					return "", fmt.Errorf("template refers to arg %d, but only %d args given", argIndex, len(args))
				}

				// Mark this arg as used
				usedArgs[argIndex] = true

				// Replace the placeholder with the argument's string representation
				result.WriteString(fmt.Sprint(args[argIndex]))

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
