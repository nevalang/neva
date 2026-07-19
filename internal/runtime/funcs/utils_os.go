package funcs

import (
	"context"

	"github.com/nevalang/neva/internal/runtime"
)

// runSignalLoop handles signal-triggered runtime functions.
func runSignalLoop(
	ctx context.Context,
	sigIn runtime.SingleInport,
	resOut runtime.SingleOutport,
	errOut *runtime.SingleOutport,
	step func() (runtime.Msg, error),
) {
	for {
		if _, received := sigIn.Receive(ctx); !received {
			return
		}

		result, err := step()
		if err != nil {
			if errOut == nil {
				panic(err)
			}

			if !errOut.Send(ctx, errFromErr(err)) {
				return
			}

			continue
		}

		if !resOut.Send(ctx, result) {
			return
		}
	}
}

// createSignalLoop prepares a signal-triggered runtime function loop.
func createSignalLoop(
	rio runtime.IO,
	withErr bool,
	step func() (runtime.Msg, error),
) (func(ctx context.Context), error) {
	sigIn, err := singleInport(rio, "sig")
	if err != nil {
		return nil, err
	}

	resOut, err := singleOutport(rio, "res")
	if err != nil {
		return nil, err
	}

	var errOut *runtime.SingleOutport
	if withErr {
		singleErrOut, err := singleOutport(rio, "err")
		if err != nil {
			return nil, err
		}

		errOut = &singleErrOut
	}

	return func(ctx context.Context) {
		runSignalLoop(ctx, sigIn, resOut, errOut, step)
	}, nil
}

// createUnaryLoop prepares a one-input runtime function loop.
func createUnaryLoop(
	rio runtime.IO,
	inName string,
	withErr bool,
	step func(runtime.OrderedMsg) (runtime.Msg, error),
) (func(ctx context.Context), error) {
	inputIn, err := singleInport(rio, inName)
	if err != nil {
		return nil, err
	}

	resOut, err := singleOutport(rio, "res")
	if err != nil {
		return nil, err
	}

	var errOut *runtime.SingleOutport
	if withErr {
		singleErrOut, err := singleOutport(rio, "err")
		if err != nil {
			return nil, err
		}

		errOut = &singleErrOut
	}

	return func(ctx context.Context) {
		runUnaryLoop(ctx, inputIn, resOut, errOut, step)
	}, nil
}

// createBinaryLoop prepares a two-input runtime function loop.
func createBinaryLoop(
	rio runtime.IO,
	firstName string,
	secondName string,
	step func(runtime.OrderedMsg, runtime.OrderedMsg) (runtime.Msg, error),
) (func(ctx context.Context), error) {
	firstIn, err := singleInport(rio, firstName)
	if err != nil {
		return nil, err
	}

	secondIn, err := singleInport(rio, secondName)
	if err != nil {
		return nil, err
	}

	resOut, err := singleOutport(rio, "res")
	if err != nil {
		return nil, err
	}

	singleErrOut, err := singleOutport(rio, "err")
	if err != nil {
		return nil, err
	}

	return func(ctx context.Context) {
		runBinaryLoop(ctx, firstIn, secondIn, resOut, &singleErrOut, step)
	}, nil
}

// runUnaryLoop handles runtime functions that consume one input per step.
func runUnaryLoop(
	ctx context.Context,
	inputIn runtime.SingleInport,
	resOut runtime.SingleOutport,
	errOut *runtime.SingleOutport,
	step func(runtime.OrderedMsg) (runtime.Msg, error),
) {
	for {
		input, received := inputIn.Receive(ctx)
		if !received {
			return
		}

		result, err := step(input)
		if err != nil {
			if errOut == nil {
				panic(err)
			}

			if !errOut.Send(ctx, errFromErr(err)) {
				return
			}

			continue
		}

		if !resOut.Send(ctx, result) {
			return
		}
	}
}

// runBinaryLoop handles runtime functions that consume two inputs per step.
func runBinaryLoop(
	ctx context.Context,
	firstIn runtime.SingleInport,
	secondIn runtime.SingleInport,
	resOut runtime.SingleOutport,
	errOut *runtime.SingleOutport,
	step func(runtime.OrderedMsg, runtime.OrderedMsg) (runtime.Msg, error),
) {
	for {
		firstInput, secondInput, received := receive2(ctx, firstIn, secondIn)
		if !received {
			return
		}

		result, err := step(firstInput, secondInput)
		if err != nil {
			if errOut == nil {
				panic(err)
			}

			if !errOut.Send(ctx, errFromErr(err)) {
				return
			}

			continue
		}

		if !resOut.Send(ctx, result) {
			return
		}
	}
}
