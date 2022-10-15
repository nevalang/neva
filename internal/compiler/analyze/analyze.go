package analyze

import (
	"context"

	"github.com/emil14/neva/internal/compiler/src"
)

type Analyzer struct {
}

func (a Analyzer) Analyze(ctx context.Context, prog src.Program) error {
	pkg, ok := prog.Packages[prog.RootPkg]
	if !ok {
		panic("no root pkg")
	}

	if pkg.RootComponent == "" {
		panic("no root component")
	}

	if err := a.checkPkg(pkg, prog.Packages); err != nil {
		panic(err)
	}

	return nil
}

func (a Analyzer) checkPkg(pkg src.Package, scope map[src.PkgRef]src.Package) error {
	_, ok := pkg.Components[pkg.RootComponent]
	if pkg.RootComponent != "" && !ok {
		panic("no root component")
	}

	usedTypes := make(map[string]struct{}, len(pkg.Types))

	

	return nil
}
