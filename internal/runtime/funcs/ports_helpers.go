package funcs

import (
	"fmt"

	"github.com/nevalang/neva/internal/runtime"
)

func singleInport(runtimeIO runtime.IO, name string) (runtime.SingleInport, error) {
	inport, err := runtimeIO.In.Single(name)
	if err != nil {
		return runtime.SingleInport{}, fmt.Errorf("inport %q: %w", name, err)
	}

	return inport, nil
}

func singleOutport(runtimeIO runtime.IO, name string) (runtime.SingleOutport, error) {
	outport, err := runtimeIO.Out.Single(name)
	if err != nil {
		return runtime.SingleOutport{}, fmt.Errorf("outport %q: %w", name, err)
	}

	return outport, nil
}
