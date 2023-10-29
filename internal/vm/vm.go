// Package vm contains `VirtualMachine` entity.
//
// `vm.Machine` is basically just a wrapper around `runtime.Runtime`
// that can read *.ir files, parse them and use runtime to execute them.
//
// Runtime needs some mechanism to load programs from the outer world - VM is exactly that.
// It has scary name like it's something complicated but it actually simplest thing in Nevalang.
// One might say it's not a "real VM" because it doesn't emulate a computer.
// These days VMs are defined more broadly. Even though there's no "machine" virtualized,
// we actually have a bytecode (a graph-like one, IR) to execute.
// Names are not that important as long as they are consistent.

package vm

import (
	"context"
	"fmt"

	"github.com/nevalang/neva/internal/runtime"
)

type Machine struct {
	loader  Loader
	decoder Decoder
	runtime Runtime
}

type (
	Loader interface {
		Load(ctx context.Context, path string) ([]byte, error)
	}
	Decoder interface {
		Decode([]byte) (runtime.Program, error)
	}
	Runtime interface {
		Run(context.Context, runtime.Program) (code int, err error)
	}
)

func (vm Machine) Exec(ctx context.Context, pathToIRFile string) (int, error) {
	bb, err := vm.loader.Load(ctx, pathToIRFile)
	if err != nil {
		return 0, err
	}

	runtimeProg, err := vm.decoder.Decode(bb)
	if err != nil {
		return 0, fmt.Errorf("decode: %w", err)
	}

	exitCode, err := vm.runtime.Run(context.Background(), runtimeProg)
	if err != nil {
		return 0, fmt.Errorf("run: %w", err)
	}

	return exitCode, nil
}

func New(
	loader Loader,
	decoder Decoder,
	runtime Runtime,
) Machine {
	return Machine{
		loader:  loader,
		decoder: decoder,
		runtime: runtime,
	}
}
