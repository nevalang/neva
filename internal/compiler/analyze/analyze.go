package analyze

import (
	"context"

	"github.com/emil14/neva/internal/compiler/src"
	ts "github.com/emil14/neva/pkg/types"
)

type Analyzer struct {
	// resolver ts.ExprResolver
}

// Analyze checks that:
// Program has ref to root pkg;
// Root pkg exist;
// Root pkg has root component ref;
// All pkgs are analyzed;
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

	for pkgName := range prog.Pkgs {
		if err := a.analyzePkg(pkgName, prog.Pkgs); err != nil {
			panic(err)
		}
	}

	return nil
}

// analyzePkg checks that:
// If pkg has ref to root component then it satisfies the pkg-with-root-component-specific requirements;
// There's no imports of not found pkgs;
// There's no unused imports;
// All entities are analyzed and;
// Used (exported or referenced by exported entities or root component).
func (a Analyzer) analyzePkg(pkgName string, pkgs map[string]src.Pkg) error { //nolint:unparam
	pkg := pkgs[pkgName]

	if pkg.RootComponent != "" {
		if err := a.analyzePkgWithRootComponent(pkg, pkgs); err != nil {
			panic(err)
		}
	}

	imports, err := a.makeImports(pkg, pkgs)
	if err != nil {
		panic(err)
	}

	usedImports, usedLocalEntities, err := a.analyzeEntities(pkg, imports)
	if err != nil {
		panic(err)
	}

	if err := a.analyzeUsedImportsAndEntities(
		imports, pkg.Entities,
		usedImports, usedLocalEntities,
	); err != nil {
		panic(err)
	}

	if pkg.RootComponent == "" && len(usedLocalEntities) == 0 {
		panic("package must have exported entities if it doesn't have a root component")
	}

	return nil
}

func (Analyzer) analyzeUsedImportsAndEntities(
	imports map[string]src.Pkg, entities map[string]src.Entity,
	usedImports, usedEntities map[string]struct{},
) error {
	for alias := range imports {
		if _, ok := usedImports[alias]; !ok {
			panic("unused imports")
		}
	}

	for entityName := range entities {
		if _, ok := usedEntities[entityName]; !ok {
			panic("unused local entity")
		}
	}

	return nil
}

func (a Analyzer) analyzeEntities(pkg src.Pkg, imports map[string]src.Pkg) (map[string]struct{}, map[string]struct{}, error) {
	usedImports := make(map[string]struct{}, len(pkg.Imports))
	usedLocalEntities := make(map[string]struct{}, len(pkg.Entities))

	for entityName, entity := range pkg.Entities {
		entitiesUsedByEntity, err := a.analyzeEntity(entity, imports)
		if err != nil {
			panic(err)
		}

		if entity.Exported {
			usedLocalEntities[entityName] = struct{}{}
		}

		for entityRef := range entitiesUsedByEntity {
			if entityRef.Import == "" {
				usedLocalEntities[entityName] = struct{}{}
				continue
			}
			usedImports[entityRef.Import] = struct{}{}
		}
	}

	return usedImports, usedLocalEntities, nil
}

// makeImports maps aliases to packages
func (Analyzer) makeImports(pkg src.Pkg, pkgs map[string]src.Pkg) (map[string]src.Pkg, error) {
	imports := make(map[string]src.Pkg, len(pkg.Imports))
	for alias, pkgRef := range pkg.Imports {
		importedPkg, ok := pkgs[pkgRef]
		if !ok {
			panic("imported pkg not found")
		}
		imports[alias] = importedPkg
	}
	return imports, nil
}

// analyzePkgWithRootComponent checks that:
// Entity referenced as root component exist;
// That entity is a component;
// It's not exported and;
// It satisfies root-component-specific requirements;
func (a Analyzer) analyzePkgWithRootComponent(pkg src.Pkg, pkgs map[string]src.Pkg) error {
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

	if err := a.analyzeRootComponent(entity.Component, pkg, pkgs); err != nil {
		panic(err)
	}

	return nil
}

func (a Analyzer) analyzeEntity(entity src.Entity, imports map[string]src.Pkg) (map[src.EntityRef]struct{}, error) { //nolint:unparam,lll
	// usedImports := map[string]struct{}{}

	switch entity.Kind { // https://github.com/emil14/neva/issues/186
	case src.ComponentEntity:
		_, err := a.analyzeComponent(entity.Component)
		return nil, err // todo
	case src.TypeEntity:
		// TODO
	case src.MsgEntity:
		// TODO
	case src.InterfaceEntity:
		// TODO
	default:
		panic("unknown entity type")
	}

	return map[src.EntityRef]struct{}{}, nil
}

func (a Analyzer) analyzeComponent(component src.Component) (map[string]struct{}, error) {
	if err := a.analyzeTypeParameters(component.TypeParams); err != nil {
		panic(err)
	}
	if err := a.analyzeIO(component.IO); err != nil {
		panic(err)
	}
	if err := a.analyzeNodes(component.Nodes); err != nil {
		panic(err)
	}
	if err := a.analyzeNet(component.Net); err != nil {
		panic(err)
	}
	return nil, nil
}

func (a Analyzer) analyzeTypeParameters(params []ts.Param) error {
	pp := make(map[string]struct{}, len(params))
	for _, param := range params {
		if param.Name == "" {
			panic("param name cannot be empty")
		}
		if _, ok := pp[param.Name]; ok {
			panic("parameter names must be unique")
		}
		pp[param.Name] = struct{}{}

		ts.NewDefaultResolver().Resolve()
	}

	// 	func<
	//     T,
	//     Y T,
	//     Z std.vec<Y>
	// >() {

	// }

	// <T int, Y arr<T>>

	return nil
}

func (a Analyzer) analyzeIO(src.IO) error {
	return nil
}

func (a Analyzer) analyzeNodes(map[string]src.Node) error {
	return nil
}

func (a Analyzer) analyzeNet([]src.Connection) error {
	return nil
}

// analyzeRootComponent checks root-component-specific requirements:
// Has only one type parameter without constraint;
// Doesn't have any outports;
// Has one not-array inport named `sig` that refers to that single type-parameter;
// Has at least one node (non-root components can have no nodes to implement just routing);
// All nodes has no static inports or references to anything but components (reason why pkg and pkgs needed);
func (Analyzer) analyzeRootComponent(rootComp src.Component, pkg src.Pkg, pkgs map[string]src.Pkg) error { //nolint:funlen,unparam,lll
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
