package vm

import (
	"context"
	"fmt"

	"github.com/nevalang/neva/internal/runtime"
)

type VirtualMachine struct {
	repo    Repository
	decoder Decoder
	runtime Runtime
}

type (
	Repository interface {
		ByPath(ctx context.Context, path string) ([]byte, error)
	}
	Decoder interface {
		Decode([]byte) (runtime.Program, error)
	}
	Runtime interface {
		Run(context.Context, runtime.Program) (code int, err error)
	}
)

func (vm VirtualMachine) Exec(ctx context.Context, path string) (int, error) {
	bb, err := vm.repo.ByPath(ctx, path)
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
	repo Repository,
	decoder Decoder,
	runtime Runtime,
) VirtualMachine {
	return VirtualMachine{
		repo:    repo,
		decoder: decoder,
		runtime: runtime,
	}
}
