package analyze

import (
	"context"

	"github.com/emil14/neva/internal/compiler/src"
)

type Analyzer struct {
	// resolver ts.ExprResolver
}

func (a Analyzer) Analyze(ctx context.Context, prog src.Prog) error {
	if prog.RootPkg == "" {
		panic("program must have root pkg")
	}

	rootPkg, ok := prog.Pkgs[prog.RootPkg]
	if !ok {
		panic("root pkg not found")
	}

	if rootPkg.RootComponent == "" {
		panic("root pkg must have root component")
	}

	var c int
	for _, p := range prog.Pkgs {
		c += len(p.Entities)
	}
	visited := make(map[src.EntityRef]struct{}, c)

	if err := a.analyzePkg(prog.RootPkg, prog.Pkgs, visited); err != nil {
		panic(err)
	}

	return nil
}

// analyzePkg checks
func (a Analyzer) analyzePkg(
	pkgName string,
	pkgs map[string]src.Pkg,
	visited map[src.EntityRef]struct{},
) error { //nolint:unparam
	pkg := pkgs[pkgName]

	if pkg.RootComponent != "" {
		entity, ok := pkg.Entities[pkg.RootComponent]
		if !ok {
			panic("root component not found")
		}

		if entity.Kind != src.ComponentEntity {
			panic("entity with name of the root component is not component")
		}

		if entity.Exported {
			panic("root component must not be exported")
		}

		if err := a.analyzeRootComp(entity.Component, pkg, pkgs); err != nil {
			panic(err)
		}

		return nil
	}

	usedImports := make(map[string]struct{}, len(pkg.Imports))

	for entityName, entity := range pkg.Entities {
		ref := src.EntityRef{Pkg: pkgName, Name: entityName}
		if _, ok := visited[ref]; ok {
			continue
		}

		imports, err := a.analyzeEntity(entity, pkgs)
		if err != nil {
			panic(err)
		}

		for use := range imports {
			usedImports[use] = struct{}{}
		}

		visited[ref] = struct{}{}
	}

	return nil
}

func (a Analyzer) analyzeEntity(entity src.Entity, pkgs map[string]src.Pkg) (map[string]struct{}, error) { //nolint:unparam,lll
	// usedImports := map[string]struct{}{}

	// https://github.com/emil14/neva/issues/186

	switch entity.Kind {
	case src.TypeEntity:
		// a.resolver.Resolve(entity.Type)
		// if entity.Type.Body.Empty() { // TODO move to validator?
		// 	panic("type entity body cannot be empty")
		// }
		// entity.Type.
	case src.MsgEntity:
	case src.InterfaceEntity:
	case src.ComponentEntity:
	default:
		panic("unknown entity type")
	}

	return map[string]struct{}{}, nil
}

// analyzeRootComp checks root-component-specific requirements:
// Has only one type parameter without constraint; Doesn't have any outports;
// Has one not-array inport named `sig` that refers to type-parameter; Has at least one node;
// All nodes has no static inports or references to anything but components (reason why pkg and pkgs needed);
func (Analyzer) analyzeRootComp(rootComp src.Component, pkg src.Pkg, pkgs map[string]src.Pkg) error { //nolint:funlen,unparam,lll
	if len(rootComp.TypeParams) != 1 {
		panic("root component must have one type parameter")
	}

	typeParam := rootComp.TypeParams[0]

	if !typeParam.Constr.Empty() {
		panic("type parameter of root component must not have constraint")
	}

	if len(rootComp.IO.Out) != 0 {
		panic("root component can't have outports")
	}

	if len(rootComp.IO.In) != 1 {
		panic("root component must have 1 inport")
	}

	sigInport, ok := rootComp.IO.In["sig"]
	if !ok {
		panic("root component must have 'sig' inport")
	}

	if sigInport.IsArr {
		panic("sig inport of root component can't be array port")
	}

	if !sigInport.Type.Lit.Empty() {
		panic("sig inport of root component can't have literal as type")
	}

	if sigInport.Type.Inst.Ref != typeParam.Name {
		panic("sig inport of root component must refer to type parameter")
	}

	if len(rootComp.Nodes) == 0 {
		panic("component must have nodes")
	}

	for _, node := range rootComp.Nodes {
		if len(node.StaticInports) != 0 {
			panic("root component can't have static inports")
		}

		var pkgWithEntity src.Pkg
		if node.Instance.Ref.Pkg != "" {
			p, ok := pkgs[node.Instance.Ref.Pkg]
			if !ok {
				panic("pkg not found")
			}
			pkgWithEntity = p
		} else {
			pkgWithEntity = pkg
		}

		entity, ok := pkgWithEntity.Entities[node.Instance.Ref.Name]
		if !ok {
			panic("entity not found")
		}

		if entity.Kind != src.ComponentEntity {
			panic("root component nodes can only refer to other components")
		}
	}

	return nil
}
