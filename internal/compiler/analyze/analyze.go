package analyze

import (
	"context"
	"errors"
	"fmt"

	"github.com/emil14/neva/internal/compiler/src"
	ts "github.com/emil14/neva/pkg/types"
)

var (
	ErrPkg              = errors.New("analyze package")
	ErrUsed             = errors.New("analyze used")
	ErrUnusedImport     = errors.New("unused import")
	ErrEntity           = errors.New("analyze entity")
	ErrUnusedEntity     = errors.New("unused entity")
	ErrUnknownMsgType   = errors.New("unknown msg type")
	ErrUnwantedMsgField = errors.New("unwanted msg field")
	ErrMissingMsgField  = errors.New("missing msg field")
	ErrVecEl            = errors.New("vector element")
	ErrNestedMsg        = errors.New("sub message")
	ErrReferencedMsg    = errors.New("msg not found by ref")
	ErrInassignMsg      = errors.New("msg is not assignable")
	ErrEntities         = errors.New("analyze entities")
	ErrRootComponent    = errors.New("analyze root component")
	ErrScopeRebase      = errors.New("scope rebase")
)

var h src.Helper

type Analyzer struct {
	Resolver TypeSystem
	Compator Compator
}

type (
	TypeSystem interface {
		Resolve(ts.Expr, ts.Scope) (ts.Expr, error)
	}
	Compator interface {
		Check(ts.Expr, ts.Expr, ts.TerminatorParams) error
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
			return src.Prog{}, fmt.Errorf("%w: %v: %v", ErrPkg, pkgName, err)
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

	resolvedEntities, allUsedEntities, err := a.analyzeEntities(pkg, imports, pkgs)
	if err != nil {
		return src.Pkg{}, fmt.Errorf("%w: %v", ErrEntities, err)
	}

	if err := a.analyzeUsed(pkg, allUsedEntities); err != nil {
		return src.Pkg{}, fmt.Errorf("%w: %v", ErrUsed, err)
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

func (a Analyzer) analyzeEntities(pkg src.Pkg, imports, pkgs map[string]src.Pkg) (map[string]src.Entity, map[src.EntityRef]struct{}, error) {
	resolvedPkgEntities := make(map[string]src.Entity, len(pkg.Entities))
	allUsedEntities := map[src.EntityRef]struct{}{} // both local and imported

	for entityName, entity := range pkg.Entities {
		if entity.Exported || entityName == pkg.RootComponent {
			allUsedEntities[src.EntityRef{Name: entityName}] = struct{}{} // normalize?
		}

		resolvedEntity, entitiesUsedByEntity, err := a.analyzeEntity(entityName, pkg.Entities, imports, pkgs)
		if err != nil {
			return nil, nil, fmt.Errorf("%w: %v: %v", ErrEntity, entityName, err)
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
		return fmt.Errorf("%w: %v", ErrRootComponent, err)
	}

	return nil
}

func (a Analyzer) analyzeEntity(
	name string,
	entities map[string]src.Entity,
	imports, pkgs map[string]src.Pkg,
) (
	src.Entity,
	map[src.EntityRef]struct{},
	error,
) { //nolint:unparam
	entity := entities[name]
	scope := Scope{
		pkgs:     pkgs,
		imports:  imports,
		local:    entities,
		builtins: a.builtinEntities(),
		visited:  map[src.EntityRef]struct{}{},
	}

	switch entity.Kind { // https://github.com/emil14/neva/issues/186
	case src.TypeEntity:
		resolvedDef, usedTypeEntities, err := a.analyzeType(name, scope)
		if err != nil {
			return src.Entity{}, nil, err
		}
		return src.Entity{
			Type:     resolvedDef,
			Kind:     src.TypeEntity,
			Exported: entity.Exported,
		}, usedTypeEntities, nil
	case src.MsgEntity:
		resolvedMsg, usedEntities, err := a.analyzeMsg(entity.Msg, scope, nil)
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
		return src.Entity{}, nil, errors.New("unknown entity type")
	}

	return src.Entity{}, map[src.EntityRef]struct{}{}, nil
}

func (Analyzer) builtinEntities() map[string]src.Entity {
	return map[string]src.Entity{
		"int": h.BaseTypeEntity(),
		"vec": h.BaseTypeEntity(h.ParamWithNoConstr("t")),
	}
}

func (a Analyzer) analyzeType(name string, scope Scope) (ts.Def, map[src.EntityRef]struct{}, error) {
	def, err := scope.GetType(name) // FIXME it's not enough to create scope once, it' has dynamic nature
	if err != nil {
		panic(err)
	}

	testExpr := ts.Expr{
		Inst: ts.InstExpr{
			Ref:  name,
			Args: a.getTestExprArgs(def),
		},
	}

	// TODO return simplified defs (type t1 pkg1.t0<t0> // t1<int> -> vec<int>)
	if _, err = a.Resolver.Resolve(testExpr, scope); err != nil {
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

func (a Analyzer) analyzeMsg(
	msg src.Msg,
	scope Scope,
	resolvedConstr *ts.Expr,
) (src.Msg, map[src.EntityRef]struct{}, error) {
	if msg.Ref != nil {
		subMsg, err := scope.getMsg(*msg.Ref)
		if err != nil {
			return src.Msg{}, nil, fmt.Errorf("%w: %v", ErrReferencedMsg, err)
		}
		if msg.Ref.Pkg != "" { // rebase needed
			scope, err = scope.rebase(msg.Ref.Pkg)
			if err != nil {
				return src.Msg{}, nil, fmt.Errorf("%w: %v", ErrScopeRebase, err)
			}
		}
		resolvedSubMsg, used, err := a.analyzeMsg(subMsg, scope, resolvedConstr)
		if err != nil {
			return src.Msg{}, nil, fmt.Errorf("%w: %v, %v", ErrNestedMsg, err, msg.Ref)
		}
		used[*msg.Ref] = struct{}{}
		return resolvedSubMsg, used, nil // TODO do we really want unpacking here?
	}

	resolvedType, err := a.Resolver.Resolve(msg.Value.Type, scope)
	if err != nil {
		return src.Msg{}, nil, fmt.Errorf("%w: ", err)
	}

	if resolvedConstr != nil {
		if err := a.Compator.Check(resolvedType, *resolvedConstr, ts.TerminatorParams{}); err != nil {
			return src.Msg{}, nil, fmt.Errorf("%w: %v", ErrInassignMsg, err)
		}
	}

	var msgToReturn src.Msg
	switch resolvedType.Inst.Ref { // TODO literals
	case "int":
		if msg.Value.Vec != nil {
			return src.Msg{}, nil, fmt.Errorf("%w: %v", ErrUnwantedMsgField, msg.Value.Vec)
		}
		msgToReturn = msg
	case "vec":
		if msg.Value.Int != 0 {
			return src.Msg{}, nil, fmt.Errorf("%w: %v", ErrUnwantedMsgField, msg.Value.Vec)
		}
		if msg.Value.Vec == nil {
			return src.Msg{}, nil, fmt.Errorf("%w: %v", ErrMissingMsgField, msg.Value.Vec)
		}
		vecType := resolvedType.Inst.Args[0]
		for i, el := range msg.Value.Vec {
			analyzedEl, _, err := a.analyzeMsg(el, scope, &vecType)
			if err != nil {
				return src.Msg{}, nil, fmt.Errorf("%w: #%d, err %v", ErrVecEl, i, err)
			}
			msg.Value.Vec[i] = analyzedEl
		}
	default:
		return src.Msg{}, nil, fmt.Errorf("%w: %v", ErrUnknownMsgType, resolvedType.Inst.Ref)
	}

	return msgToReturn, scope.visited, nil
}

func (a Analyzer) analyzeComponent(component src.Component) (map[string]struct{}, error) {
	if err := a.analyzeTypeParameters(component.TypeParams); err != nil {
		panic(err)
	}
	if err := a.analyzeIO(component.IO); err != nil { // TODO pass type params?
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
