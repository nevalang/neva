package funcs

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/nevalang/neva/internal/runtime"
)

type printf struct{}

func (p printf) Create(io runtime.FuncIO, _ runtime.Msg) (func(ctx context.Context), error) {
	tplIn, err := io.In.Port("tpl")
	if err != nil {
		return nil, err
	}

	argsIn, ok := io.In["args"]
	if !ok {
		return nil, fmt.Errorf("missing required input port 'args'")
	}

	argsOut, ok := io.Out["args"]
	if !ok {
		return nil, fmt.Errorf("missing required output port 'args'")
	}

	if len(argsIn) != len(argsOut) {
		return nil, fmt.Errorf("input and output ports 'args' must have the same length")
	}

	errOut, err := io.Out.Port("err")
	if err != nil {
		return nil, err
	}

	return p.handle(tplIn, argsIn, errOut, argsOut)
}

func (printf) handle(
	tplIn chan runtime.Msg,
	argsIn []chan runtime.Msg,
	errOut chan runtime.Msg,
	argsOuts []chan runtime.Msg,
) (func(ctx context.Context), error) {
	return func(ctx context.Context) {
		var (
			tpl  runtime.Msg
			args = make([]runtime.Msg, len(argsIn))
		)

		for {
			select {
			case <-ctx.Done():
				return
			case tpl = <-tplIn:
			}

			for i, argIn := range argsIn {
				select {
				case <-ctx.Done():
					return
				case arg := <-argIn:
					args[i] = arg
				}
			}

			res, err := format(tpl.Str(), args)
			if err != nil {
				select {
				case <-ctx.Done():
					return
				case errOut <- errorFromString(err.Error()):
				}
				continue
			}

			if _, err := fmt.Print(res); err != nil {
				select {
				case <-ctx.Done():
					return
				case errOut <- errorFromString(err.Error()):
				}
				continue
			}

			for i, argOut := range argsOuts {
				select {
				case <-ctx.Done():
					return
				case argOut <- args[i]:
				}
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
		return "", fmt.Errorf("not all arguments were used in the template")
	}

	return result.String(), nil
}
