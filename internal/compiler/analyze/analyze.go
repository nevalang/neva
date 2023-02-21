package analyze

import (
	"context"

	"github.com/emil14/neva/internal/compiler/src"
	ts "github.com/emil14/neva/pkg/types"
)

var h ts.Helper

type Analyzer struct {
	Resolver TypeResolver
}

type (
	TypeResolver interface {
		Resolve(ts.Expr, TypeEnv) (ts.Expr, error)
	}
	TypeEnv interface {
		Get(src.EntityRef) (ts.Def, error)
	}
)

// Analyze checks that:
// Program has ref to root pkg;
// Root pkg exist;
// Root pkg has root component ref;
// All pkgs are analyzed;
func (a Analyzer) Analyze(ctx context.Context, prog src.Prog) (src.Prog, error) {
	if prog.RootPkg == "" {
		panic("program must have root pkg")
	}

	rootPkg, ok := prog.Pkgs[prog.RootPkg]
	if !ok {
		panic("root pkg not found")
	}

	if rootPkg.RootComponent == "" { // is executable
		panic("root pkg must have root component")
	}

	resolvedPkgs := make(map[string]src.Pkg, len(prog.Pkgs))
	for pkgName := range prog.Pkgs { // we know main component must be there
		resolvedPkg, err := a.analyzePkg(pkgName, prog.Pkgs)
		if err != nil {
			panic(err)
		}
		resolvedPkgs[pkgName] = resolvedPkg
	}

	return src.Prog{
		Pkgs:    resolvedPkgs,
		RootPkg: prog.RootPkg,
	}, nil
}

// analyzePkg checks that:
// If pkg has ref to root component then it satisfies the pkg-with-root-component-specific requirements;
// There's no imports of not found pkgs;
// There's no unused imports;
// All entities are analyzed and;
// Used (exported or referenced by exported entities or root component).
func (a Analyzer) analyzePkg(pkgName string, pkgs map[string]src.Pkg) (src.Pkg, error) { //nolint:unparam
	pkg := pkgs[pkgName]

	if pkg.RootComponent != "" { // is executable
		if err := a.analyzePkgWithRootComponent(pkg, pkgs); err != nil {
			panic(err)
		}
	} else if len(a.getExports(pkg.Entities)) == 0 {
		panic("package must have exported entities if it doesn't have a root component")
	}

	imports, err := a.getPkgImports(pkg.Imports, pkgs)
	if err != nil {
		panic(err)
	} // at this we know all pkg's imports points to existing pkgs

	resolvedEntities, entitiesUsedByPkg, err := a.analyzeEntities(pkg, imports)
	if err != nil {
		panic(err)
	}

	if err := a.analyzeUsed(pkg, entitiesUsedByPkg); err != nil {
		panic(err)
	}

	return src.Pkg{
		Entities:      resolvedEntities,
		Imports:       pkg.Imports,
		RootComponent: pkg.RootComponent,
	}, nil
}

// getExports returns only exported entities
func (a Analyzer) getExports(entities map[string]src.Entity) map[string]src.Entity {
	exports := make(map[string]src.Entity, len(entities))
	for name, entity := range entities {
		exports[name] = entity
	}
	return exports
}

// analyzeUsed returns error if there're unused imports or entities
func (Analyzer) analyzeUsed(pkg src.Pkg, usedEntities map[src.EntityRef]struct{}) error {
	usedImports := map[string]struct{}{}
	usedLocalEntities := map[string]struct{}{}

	for ref := range usedEntities {
		if ref.Pkg == "" {
			usedLocalEntities[ref.Name] = struct{}{}
		} else {
			usedImports[ref.Pkg] = struct{}{}
		}
	}

	for alias := range pkg.Imports {
		if _, ok := usedImports[alias]; !ok {
			panic("unused imports")
		}
	}

	for entityName := range pkg.Entities {
		if _, ok := usedLocalEntities[entityName]; !ok {
			panic("unused local entity")
		}
	}

	return nil
}

func (a Analyzer) analyzeEntities(pkg src.Pkg, imports map[string]src.Pkg) (map[string]src.Entity, map[src.EntityRef]struct{}, error) {
	resolvedPkgEntities := make(map[string]src.Entity, len(pkg.Entities))
	entitiesUsedByPkg := map[src.EntityRef]struct{}{} // both local and imported

	for entityName, entity := range pkg.Entities {
		if entity.Exported || entityName == pkg.RootComponent {
			entitiesUsedByPkg[src.EntityRef{Name: entityName}] = struct{}{}
		}

		resolvedEntity, entitiesUsedByEntity, err := a.analyzeEntity(entity, imports)
		if err != nil {
			panic(err)
		}

		for entityRef := range entitiesUsedByEntity {
			entitiesUsedByPkg[entityRef] = struct{}{}
		}

		resolvedPkgEntities[entityName] = resolvedEntity
	}

	return resolvedPkgEntities, entitiesUsedByPkg, nil
}

// getPkgImports maps aliases to packages
func (Analyzer) getPkgImports(pkgImports map[string]string, pkgs map[string]src.Pkg) (map[string]src.Pkg, error) {
	imports := make(map[string]src.Pkg, len(pkgImports))
	for alias, pkgRef := range pkgImports {
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

func (a Analyzer) analyzeEntity(entity src.Entity, imports map[string]src.Pkg) (
	src.Entity,
	map[src.EntityRef]struct{},
	error,
) { //nolint:unparam
	switch entity.Kind { // https://github.com/emil14/neva/issues/186
	case src.TypeEntity:
		resolvedType, referencedTypes, err := a.analyzeType(entity.Type)
		if err != nil {
			return src.Entity{}, nil, err
		}
		return src.Entity{
			Type:     resolvedType,
			Kind:     src.TypeEntity,
			Exported: entity.Exported,
		}, referencedTypes, nil
	case src.MsgEntity:
		resolvedMsg, usedEntities, err := a.analyzeMsg(entity.Msg)
		if err != nil {
			return src.Entity{}, nil, err
		}
		return src.Entity{
			Msg:      resolvedMsg,
			Kind:     src.MsgEntity,
			Exported: entity.Exported,
		}, usedEntities, nil
	case src.InterfaceEntity:
	case src.ComponentEntity:
		_, err := a.analyzeComponent(entity.Component)
		return src.Entity{}, nil, err
	default:
		panic("unknown entity type")
	}

	return src.Entity{}, map[src.EntityRef]struct{}{}, nil
}

func (a Analyzer) analyzeMsg(msg src.Msg) (src.Msg, map[src.EntityRef]struct{}, error) {
	return src.Msg{}, nil, nil
}

func (a Analyzer) analyzeType(def ts.Def) (ts.Def, map[src.EntityRef]struct{}, error) {
	// arg=constr
	

	expr := ts.Expr{
		Lit:  ts.LitExpr{},
		Inst: ts.InstExpr{},
	}

	a.Resolver.Resolve(expr, nil)

	return ts.Def{}, nil, nil
}

func (Analyzer) newMethod(def ts.Def) []ts.Expr {
	args := make([]ts.Expr, len(def.Params))
	for _, param := range def.Params {
		if param.Constr.Empty() {
			args = append(args, h.Inst("int"))
		} else {
			args = append(args, param.Constr)
		}
	}
	return args
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
	// pp := make(map[string]struct{}, len(params))

	// for _, param := range params {
	// 	if param.Name == "" {
	// 		panic("param name cannot be empty")
	// 	}
	// 	if _, ok := pp[param.Name]; ok {
	// 		panic("parameter names must be unique")
	// 	}
	// 	pp[param.Name] = struct{}{}
	// }

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
func (a Analyzer) analyzeRootComponent(rootComp src.Component, pkg src.Pkg, pkgs map[string]src.Pkg) error { //nolint:funlen,unparam,lll
	if len(rootComp.TypeParams) != 1 {
		panic("root component must have one type parameter")
	}

	typeParam := rootComp.TypeParams[0]

	if !typeParam.Constr.Empty() {
		panic("type parameter of root component must not have constraint")
	}

	if err := a.analyzeRootComponentIO(rootComp.IO, typeParam); err != nil {
		return err
	}

	if err := a.analyzeRootComponentNodes(rootComp.Nodes, pkg, pkgs); err != nil {
		return err
	}

	return nil
}

func (Analyzer) analyzeRootComponentIO(io src.IO, typeParam ts.Param) error {
	if len(io.Out) != 0 {
		panic("root component can't have outports")
	}

	if len(io.In) != 1 {
		panic("root component must have 1 inport")
	}

	sigInport, ok := io.In["sig"]
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

	return nil
}

func (Analyzer) analyzeRootComponentNodes(nodes map[string]src.Node, pkg src.Pkg, pkgs map[string]src.Pkg) error {
	if len(nodes) == 0 {
		panic("component must have nodes")
	}

	for _, node := range nodes {
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
		panic("entity not found")
		if !ok {
		}

		if entity.Kind != src.ComponentEntity {
			panic("root component nodes can only refer to other components")
		}
	}

	return nil
}
