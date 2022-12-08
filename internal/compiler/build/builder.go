package build

import (
	"context"

	"github.com/emil14/neva/internal/compiler/src"
)

type Builder struct{}

func (b Builder) Build(ctx context.Context, path string) (src.Program, error) {
	return src.Program{}, nil // TODO
}

func (b Builder) build(ctx context.Context, path string, seen map[string]src.Pkg) {
	// check seen first
	// if pkg not seen, download and parse pkg file
	// for every import call recursion
	// when there's no imports build src.Pkg and return
}
