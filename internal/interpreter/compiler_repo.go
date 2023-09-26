package interpreter

import (
	"context"

	"github.com/nevalang/neva/internal/pkg/disk"
	"github.com/nevalang/neva/pkg/ir"
)

// CompilerRepo is compiler.Repo implementation that
// allows compiler to read source code from disc but disallows it to write IR back.
type CompilerRepo struct {
	disk.Repository
}

func (c CompilerRepo) Save(ctx context.Context, path string, prog *ir.Program) error {
	return nil
}

func NewCompilerRepo(repo disk.Repository) CompilerRepo {
	return CompilerRepo{
		Repository: repo,
	}
}
