package funcs

import (
	"context"
	"fmt"
	"sync"

	"github.com/nevalang/neva/internal/runtime"
)

type binaryOperator func(left runtime.Msg, right runtime.Msg) runtime.Msg

type unaryOperator func(input runtime.Msg) runtime.Msg

func createBinaryFuncConcurrent(runtimeIO runtime.IO, apply binaryOperator) (func(context.Context), error) {
	leftInput, rightInput, resultOutput, err := resolveBinaryPorts(runtimeIO)
	if err != nil {
		return nil, err
	}

	return func(ctx context.Context) {
		for {
			var leftValue runtime.Msg
			var rightValue runtime.Msg
			var leftAvailable bool
			var rightAvailable bool

			var group sync.WaitGroup
			group.Go(func() {
				leftValue, leftAvailable = leftInput.Receive(ctx)
			})
			group.Go(func() {
				rightValue, rightAvailable = rightInput.Receive(ctx)
			})
			group.Wait()

			if !leftAvailable || !rightAvailable {
				return
			}

			if !resultOutput.Send(ctx, apply(leftValue, rightValue)) {
				return
			}
		}
	}, nil
}

func createUnaryFunc(runtimeIO runtime.IO, apply unaryOperator) (func(context.Context), error) {
	dataInput, err := getSingleInport(runtimeIO.In, "data")
	if err != nil {
		return nil, err
	}

	resultOutput, err := getSingleOutport(runtimeIO.Out, "res")
	if err != nil {
		return nil, err
	}

	return func(ctx context.Context) {
		for {
			inputValue, available := dataInput.Receive(ctx)
			if !available {
				return
			}

			if !resultOutput.Send(ctx, apply(inputValue)) {
				return
			}
		}
	}, nil
}

func resolveBinaryPorts(runtimeIO runtime.IO) (runtime.SingleInport, runtime.SingleInport, runtime.SingleOutport, error) {
	leftInput, err := getSingleInport(runtimeIO.In, "left")
	if err != nil {
		return runtime.SingleInport{}, runtime.SingleInport{}, runtime.SingleOutport{}, err
	}

	rightInput, err := getSingleInport(runtimeIO.In, "right")
	if err != nil {
		return runtime.SingleInport{}, runtime.SingleInport{}, runtime.SingleOutport{}, err
	}

	resultOutput, err := getSingleOutport(runtimeIO.Out, "res")
	if err != nil {
		return runtime.SingleInport{}, runtime.SingleInport{}, runtime.SingleOutport{}, err
	}

	return leftInput, rightInput, resultOutput, nil
}

func getSingleInport(inports runtime.Inports, portName string) (runtime.SingleInport, error) {
	inport, err := inports.Single(portName)
	if err != nil {
		return runtime.SingleInport{}, fmt.Errorf("resolve inport %q: %w", portName, err)
	}

	return inport, nil
}

func getSingleOutport(outports runtime.Outports, portName string) (runtime.SingleOutport, error) {
	outport, err := outports.Single(portName)
	if err != nil {
		return runtime.SingleOutport{}, fmt.Errorf("resolve outport %q: %w", portName, err)
	}

	return outport, nil
}
