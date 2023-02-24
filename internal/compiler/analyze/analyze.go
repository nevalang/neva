package analyze

import (
	"context"
	"errors"
	"fmt"

	"github.com/emil14/neva/internal/compiler/src"
	ts "github.com/emil14/neva/pkg/types"
)

var (
	ErrAnalyzePkg    = errors.New("analyze package")
	ErrAnalyzeUsed   = errors.New("analyze used")
	ErrUnusedImport  = errors.New("unused import")
	ErrAnalyzeEntity = errors.New("analyze entity")
	ErrUnusedEntity  = errors.New("unused entity")
)

var h ts.Helper

type Analyzer struct {
	Resolver TypeResolver
}

type (
	TypeResolver interface {
		Resolve(ts.Expr, ts.Scope) (ts.Expr, error)
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

	if rootPkg.RootComponent == "" {
		panic("root pkg must have root component")
	}

	resolvedPkgs := make(map[string]src.Pkg, len(prog.Pkgs))
	for pkgName := range prog.Pkgs {
		resolvedPkg, err := a.analyzePkg(pkgName, prog.Pkgs)
		if err != nil {
			return src.Prog{}, fmt.Errorf("%w: pkg name %v, err %v", ErrAnalyzePkg, pkgName, err)
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

	imports, err := a.getImports(pkg.Imports, pkgs)
	if err != nil {
		panic(err)
	} // at this we know all pkg's imports points to existing pkgs

	resolvedEntities, allUsedEntities, err := a.analyzeEntities(pkg, imports)
	if err != nil {
		panic(err)
	}

	if err := a.analyzeUsed(pkg, allUsedEntities); err != nil {
		return src.Pkg{}, fmt.Errorf("%w: %v", ErrAnalyzeUsed, err)
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

	for ref := range usedEntities { // FIXME no pkg_2.*
		if ref.Pkg == "" {
			usedLocalEntities[ref.Name] = struct{}{}
		} else {
			usedImports[ref.Pkg] = struct{}{}
		}
	}

	for alias := range pkg.Imports {
		if _, ok := usedImports[alias]; !ok {
			return fmt.Errorf("%w: %v", ErrUnusedImport, alias)
		}
	}

	for entityName := range pkg.Entities {
		if _, ok := usedLocalEntities[entityName]; !ok {
			return fmt.Errorf("%w: %v", ErrUnusedEntity, entityName)
		}
	}

	return nil
}

func (a Analyzer) analyzeEntities(pkg src.Pkg, imports map[string]src.Pkg) (map[string]src.Entity, map[src.EntityRef]struct{}, error) {
	resolvedPkgEntities := make(map[string]src.Entity, len(pkg.Entities))
	allUsedEntities := map[src.EntityRef]struct{}{} // both local and imported

	for entityName, entity := range pkg.Entities {
		if entity.Exported || entityName == pkg.RootComponent {
			allUsedEntities[src.EntityRef{Name: entityName}] = struct{}{} // normalize?
		}

		resolvedEntity, entitiesUsedByEntity, err := a.analyzeEntity(entityName, pkg.Entities, imports)
		if err != nil {
			return nil, nil, fmt.Errorf("%w: name %v, err %v", ErrAnalyzeEntity, entityName, err)
		}

		for entityRef := range entitiesUsedByEntity {
			allUsedEntities[entityRef] = struct{}{}
		}

		resolvedPkgEntities[entityName] = resolvedEntity
	}

	return resolvedPkgEntities, allUsedEntities, nil
}

// getImports maps aliases to packages
func (Analyzer) getImports(pkgImports map[string]string, pkgs map[string]src.Pkg) (map[string]src.Pkg, error) {
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

func (a Analyzer) analyzeEntity(
	entityName string,
	entities map[string]src.Entity,
	imports map[string]src.Pkg,
) (
	src.Entity,
	map[src.EntityRef]struct{},
	error,
) { //nolint:unparam
	entity := entities[entityName]

	switch entity.Kind { // https://github.com/emil14/neva/issues/186
	case src.TypeEntity:
		builtins := a.builtinTypes()
		resolvedDef, usedTypeEntities, err := a.analyzeType(
			entityName,
			imports,
			a.getTypes(entities), // TODO move out to reduce number of calls
			builtins,
		) // TODO move builtins
		if err != nil {
			return src.Entity{}, nil, err
		}
		return src.Entity{
			Type:     resolvedDef,
			Kind:     src.TypeEntity,
			Exported: entity.Exported,
		}, usedTypeEntities, nil
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

func (Analyzer) builtinTypes() map[string]ts.Def { // TODO move?
	return map[string]ts.Def{
		"int": h.BaseDef(),
		"vec": h.BaseDef(h.ParamWithNoConstr("t")),
	}
}

func (a Analyzer) getTypes(entities map[string]src.Entity) map[string]ts.Def {
	types := make(map[string]ts.Def, len(entities))
	for name, entity := range entities {
		if entity.Kind == src.TypeEntity {
			types[name] = entity.Type
		}
	}
	return types
}

func (a Analyzer) analyzeMsg(msg src.Msg) (src.Msg, map[src.EntityRef]struct{}, error) {
	return src.Msg{}, nil, nil
}

func (a Analyzer) analyzeType(
	typeName string,
	imports map[string]src.Pkg,
	localTypes, builtinTypes map[string]ts.Def,
) (
	ts.Def,
	map[src.EntityRef]struct{},
	error,
) {
	def := localTypes[typeName]

	testExpr := ts.Expr{
		Inst: ts.InstExpr{
			Ref:  typeName,
			Args: a.getTestExprArgs(def),
		},
	}

	scope := Scope{
		imports:  imports,
		local:    localTypes,
		builtins: builtinTypes,
		visited:  map[src.EntityRef]struct{}{},
	}

	_, err := a.Resolver.Resolve(testExpr, scope) // TODO return simplified defs (type t1 pkg1.t0<t0> // t1<int> -> vec<int>)
	if err != nil {
		return ts.Def{}, nil, err
	}

	return def, scope.visited, nil
}

func (Analyzer) getTestExprArgs(def ts.Def) []ts.Expr {
	args := make([]ts.Expr, 0, len(def.Params))
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

// Makes sure that:
// All ports are used;
// All nodes are used;
// All nodes are known;
func (a Analyzer) analyzeNet([]src.Connection) error {
	return nil
}
