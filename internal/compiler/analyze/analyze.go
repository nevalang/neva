package analyze

import (
	"context"
	"errors"

	"github.com/emil14/neva/internal/compiler/src"
)

type Analyzer struct {
	c   Comparator
	r   Resolver
	ts  TypeScope
	std StdLib
}

type StdLib interface {
	Get(ref src.PkgRef) (src.PkgRef, error)
}

type TypeScope map[src.TypeRef]src.Type

func (a Analyzer) Analyze(ctx context.Context, prog src.Program) error {
	rootPkg, ok := prog.Pkgs[prog.RootPkg]
	if !ok {
		panic("no root pkg")
	}

	if rootPkg.Entities.Components.Root == nil {
		panic("no root component in root pkg")
	}

	seen, err := a.analyzePkg(rootPkg, prog.Pkgs)
	if err != nil {
		panic(err)
	}

	for ref := range prog.Pkgs {
		if _, ok := seen[ref]; !ok {
			return errors.New("!ok seenPackages")
		}
	}

	return nil
}

type PkgScope map[src.PkgRef]src.Pkg

func (a Analyzer) analyzePkg(pkg src.Pkg, scope PkgScope) (map[src.PkgRef]struct{}, error) {
	// is root or lib
	if pkg.Entities.Components.Root != nil {
		_, ok := pkg.Entities.Components.Scope[*pkg.Entities.Components.Root]
		if !ok {
			panic("no root component")
		}
	} else {
		if pkg.Entities.Components.Scope == nil &&
			pkg.Entities.Messages == nil &&
			pkg.Entities.Types == nil {
			return nil, errors.New("non-root pkg without exports")
		}
	}

	seenPkgs := map[src.PkgRef]struct{}{}
	for _, v := range pkg.Imports {
		seenPkgs[src.PkgRef{}]
	}

	seen := make(map[src.PkgEntityRef]struct{}, len(pkg.Exports))

	for _, ref := range pkg.Exports {
		name := ref.LocalName

		switch ref.Kind {
		case src.ComponentEntity:
			cmp, ok := pkg.Entities.Components[name]
			if !ok {
				panic(name)
			}
			if err := a.analyzeComponent(cmp, pkg, scope); err != nil {
				panic(err)
			}
		case src.TypeEntity:
			typ, ok := pkg.Entities.Types[name]
			if !ok {
				panic(name)
			}
			if err := a.analyzeType(typ, pkg, scope); err != nil {
				panic(err)
			}
		case src.MsgEntity:
			msg, ok := pkg.Entities.Messages[name]
			if !ok {
				panic(name)
			}
			if err := a.analyzeMsg(msg, pkg, scope); err != nil {
				panic(err)
			}
		}

		seen[ref] = struct{}{}
	}

	return nil, nil
}

func (a Analyzer) analyzeComponent(cmp src.Component, pkg src.Pkg, pkgScope PkgScope) error {
	for _, port := range cmp.Interface.IO.In {
		if _, err := a.r.Resolve(port.TypeExpr, a.ts); err != nil {
			return err
		}
	}
	return nil
}

func (a Analyzer) analyzeType(cmp src.Type, pkg src.Pkg, pkgScope PkgScope) error {
	return nil
}

func (a Analyzer) analyzeMsg(cmp src.Msg, pkg src.Pkg, pkgScope PkgScope) error {
	return nil
}
