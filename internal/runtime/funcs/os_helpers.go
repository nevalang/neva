package funcs

import (
	"context"
	"fmt"

	"github.com/nevalang/neva/internal/runtime"
)

func singleInport(io runtime.IO, name string) (runtime.SingleInport, error) {
	inport, err := io.In.Single(name)
	if err != nil {
		return runtime.SingleInport{}, fmt.Errorf("get inport %q: %w", name, err)
	}

	return inport, nil
}

func singleOutport(io runtime.IO, name string) (runtime.SingleOutport, error) {
	outport, err := io.Out.Single(name)
	if err != nil {
		return runtime.SingleOutport{}, fmt.Errorf("get outport %q: %w", name, err)
	}

	return outport, nil
}

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

func runBinaryLoop(
	ctx context.Context,
	firstIn runtime.SingleInport,
	secondIn runtime.SingleInport,
	resOut runtime.SingleOutport,
	errOut *runtime.SingleOutport,
	step func(runtime.OrderedMsg, runtime.OrderedMsg) (runtime.Msg, error),
) {
	for {
		firstInput, firstReceived := firstIn.Receive(ctx)
		if !firstReceived {
			return
		}

		secondInput, secondReceived := secondIn.Receive(ctx)
		if !secondReceived {
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
