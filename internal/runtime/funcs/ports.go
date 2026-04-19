package funcs

import (
	"fmt"

	"github.com/nevalang/neva/internal/runtime"
)

func singleIn(runtimeIO runtime.IO, portName string) (runtime.SingleInport, error) {
	port, err := runtimeIO.In.Single(portName)
	if err != nil {
		return runtime.SingleInport{}, fmt.Errorf("input %q: %w", portName, err)
	}
	return port, nil
}

func arrayIn(runtimeIO runtime.IO, portName string) (runtime.ArrayInport, error) {
	port, err := runtimeIO.In.Array(portName)
	if err != nil {
		return runtime.ArrayInport{}, fmt.Errorf("input %q: %w", portName, err)
	}
	return port, nil
}

func singleOut(runtimeIO runtime.IO, portName string) (runtime.SingleOutport, error) {
	port, err := runtimeIO.Out.Single(portName)
	if err != nil {
		return runtime.SingleOutport{}, fmt.Errorf("output %q: %w", portName, err)
	}
	return port, nil
}

func arrayOut(runtimeIO runtime.IO, portName string) (runtime.ArrayOutport, error) {
	port, err := runtimeIO.Out.Array(portName)
	if err != nil {
		return runtime.ArrayOutport{}, fmt.Errorf("output %q: %w", portName, err)
	}
	return port, nil
}
