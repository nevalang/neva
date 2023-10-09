package analyzer

import (
	"errors"
	"fmt"
	"strings"

	"github.com/nevalang/neva/internal/compiler/src"
	ts "github.com/nevalang/neva/pkg/typesystem"
)

type Analyzer struct {
	resolver ts.Resolver
}

// Analyze method formats error from a.analyze so end-user can easily understand what's wrong.
func (a Analyzer) Analyze(prog src.Program) error {
	return a.analyze(prog)
}

var (
	ErrEmptyProgram      = errors.New("empty program")
	ErrMainPkgNotFound   = errors.New("main package not found")
	ErrEmptyPkg          = errors.New("package must not be empty")
	ErrUnknownEntityKind = errors.New("unknown entity kind")
)

// analyze returns error if program is invalid. It also modifies program by resolving types.
func (a Analyzer) analyze(prog src.Program) error {
	if len(prog) == 0 {
		return ErrEmptyProgram
	}

	mainPkg, ok := prog["main"]
	if !ok {
		return ErrMainPkgNotFound
	}

	if err := a.mainSpecificPkgValidation(mainPkg, prog); err != nil {
		return fmt.Errorf("main specific pkg validation: %w", err)
	}

	for pkgName := range prog {
		resolvedPkg, err := a.analyzePkg(pkgName, prog)
		if err != nil {
			return fmt.Errorf("analyze pkg: %v: %w", pkgName, err)
		}
		prog[pkgName] = resolvedPkg
	}

	return nil
}

// TODO check that there's no 2 entities with the same name
// and that there's no unused entities.
func (a Analyzer) analyzePkg(pkgName string, prog src.Program) (src.Package, error) {
	if len(pkgName) == 0 {
		return nil, ErrEmptyPkg
	}

	resolvedPkg := make(map[string]src.File, len(pkgName))
	for fileName, file := range prog[pkgName] {
		resolvedPkg[fileName] = src.File{
			Imports:  file.Imports,
			Entities: make(map[string]src.Entity, len(file.Entities)),
		}
	}

	if err := prog[pkgName].Entities(func(entity src.Entity, entityName, fileName string) error {
		scope := src.Scope{
			Prog: prog,
			Loc: src.ScopeLocation{
				PkgName:  pkgName,
				FileName: fileName,
			},
		}
		resolvedEntity, err := a.analyzeEntity(entity, scope)
		if err != nil {
			return fmt.Errorf("analyze entity: %v: %v: %w", entityName, fileName, err)
		}
		resolvedPkg[fileName].Entities[entityName] = resolvedEntity
		return nil
	}); err != nil {
		return nil, fmt.Errorf("entities: %w", err)
	}

	return resolvedPkg, nil
}

func (a Analyzer) analyzeEntity(entity src.Entity, scope src.Scope) (src.Entity, error) {
	resolvedEntity := src.Entity{
		Exported: entity.Exported,
		Kind:     entity.Kind,
	}
	isStd := strings.HasPrefix(scope.Loc.PkgName, "std/")

	switch entity.Kind {
	case src.TypeEntity:
		resolvedTypeDef, err := a.analyzeTypeDef(entity.Type, scope, analyzeTypeDefParams{isStd})
		if err != nil {
			return src.Entity{}, fmt.Errorf("resolve type: %w", err)
		}
		resolvedEntity.Type = resolvedTypeDef
	case src.ConstEntity:
		resolvedConst, err := a.analyzeConst(entity.Const, scope)
		if err != nil {
			return src.Entity{}, fmt.Errorf("analyze const: %w", err)
		}
		resolvedEntity.Const = resolvedConst
	case src.InterfaceEntity:
		resolvedInterface, err := a.analyzeInterface(entity.Interface, scope, analyzeInterfaceParams{
			allowEmptyInports:  false,
			allowEmptyOutports: isStd,
		})
		if err != nil {
			return src.Entity{}, fmt.Errorf("analyze interface: %w", err)
		}
		resolvedEntity.Interface = resolvedInterface
	case src.ComponentEntity:
		resolvedComp, err := a.analyzeComponent(entity.Component, scope, analyzeComponentParams{
			iface: analyzeInterfaceParams{
				allowEmptyInports:  isStd, // e.g. `Const` component has no inports
				allowEmptyOutports: false,
			},
		})
		if err != nil {
			return src.Entity{}, fmt.Errorf("analyze component: %w", err)
		}
		resolvedEntity.Component = resolvedComp
	default:
		return src.Entity{}, fmt.Errorf("%w: %v", ErrUnknownEntityKind, entity.Kind)
	}

	return resolvedEntity, nil
}

func MustNew(resolver ts.Resolver) Analyzer {
	return Analyzer{
		resolver: resolver,
	}
}
