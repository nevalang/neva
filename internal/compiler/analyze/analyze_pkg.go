package analyze

import (
	"errors"

	"github.com/emil14/neva/internal/compiler/src"
)

var (
	ErrEntities      = errors.New("analyze entities")
	ErrUsed          = errors.New("analyze used")
	ErrExecutablePkg = errors.New("analyze package with root component")
	ErrUselessPkg    = errors.New("package without root component must have exports")
)

// analyzePkg checks that:
// If pkg has ref to root component then it satisfies the pkg-with-root-component-specific requirements;
// There's no imports of not found pkgs;
// There's no unused imports;
// All entities are analyzed and;
// Used (exported or referenced by exported entities or root component).
func (a Analyzer) analyzePkg(pkgName string, pkgs map[string]src.Pkg) (src.Pkg, error) { //nolint:unparam
	pkg := pkgs[pkgName]

	if pkg.RootComponent != "" {
		if err := a.analyzeExecutablePkg(pkg, pkgs); err != nil {
			return src.Pkg{}, errors.Join(ErrExecutablePkg, err)
		}
	} else if len(a.getExports(pkg.Entities)) == 0 {
		return src.Pkg{}, ErrUselessPkg
	}

	scope := Scope{
		base:     pkgName,
		pkgs:     pkgs,
		builtins: a.builtinEntities(),
		visited:  map[src.EntityRef]struct{}{},
	}

	resolvedEntities, used, err := a.analyzeEntities(pkg, scope)
	if err != nil {
		return src.Pkg{}, errors.Join(ErrEntities, err)
	}

	if err := a.analyzeUsed(pkg, used); err != nil {
		return src.Pkg{}, errors.Join(ErrUsed, err)
	}

	return src.Pkg{
		Entities:      resolvedEntities,
		Imports:       pkg.Imports,
		RootComponent: pkg.RootComponent,
	}, nil
}
