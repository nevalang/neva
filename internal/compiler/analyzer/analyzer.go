package analyzer

import (
	"errors"
	"fmt"

	"github.com/nevalang/neva/internal/compiler/src"
	ts "github.com/nevalang/neva/pkg/typesystem"
)

type Analyzer struct {
	prog     src.Program
	resolver ts.Resolver
}

var (
	ErrEmptyProgram      = errors.New("empty program")
	ErrMainPkgNotFound   = errors.New("main package not found")
	ErrEmptyPkg          = errors.New("package must not be empty")
	ErrUnknownEntityKind = errors.New("unknown entity kind")
)

// Analyze returns error if program is invalid. It also modifies program by resolving types.
func (a Analyzer) Analyze(prog src.Program) error {
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

	for name, pkg := range prog {
		resolvedPkg, err := a.analyzePkg(pkg, prog)
		if err != nil {
			return fmt.Errorf("analyze pkg: %v: %w", name, err)
		}
		prog[name] = resolvedPkg
	}

	return nil
}

func (a Analyzer) analyzePkg(pkg src.Package, prog src.Program) (src.Package, error) {
	if len(pkg) == 0 {
		return nil, ErrEmptyPkg
	}

	resolvedPkg := make(map[string]src.File, len(pkg))
	for fileName, file := range pkg {
		resolvedPkg[fileName] = src.File{
			Imports:  file.Imports,
			Entities: make(map[string]src.Entity, len(file.Entities)),
		}
	}

	if err := pkg.Entities(func(entity src.Entity, entityName, fileName string) error {
		resolvedEntity, err := a.analyzeEntity(entity, prog)
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

func (a Analyzer) analyzeEntity(entity src.Entity, prog src.Program) (src.Entity, error) {
	resolvedEntity := src.Entity{
		Exported: entity.Exported,
		Kind:     entity.Kind,
	}

	switch entity.Kind {
	case src.TypeEntity:
		resolvedTypeDef, err := a.analyzeTypeDef(entity.Type, Scope{prog: prog})
		if err != nil {
			return src.Entity{}, fmt.Errorf("resolve type: %w", err)
		}
		resolvedEntity.Type = resolvedTypeDef
	case src.ConstEntity:
		resolvedConst, err := a.analyzeConst(entity.Const, prog)
		if err != nil {
			return src.Entity{}, fmt.Errorf("analyze const: %w", err)
		}
		resolvedEntity.Const = resolvedConst
	case src.InterfaceEntity:
		resolvedInterface, err := a.analyzeInterface(entity.Interface, prog)
		if err != nil {
			return src.Entity{}, fmt.Errorf("analyze interface: %w", err)
		}
		resolvedEntity.Interface = resolvedInterface
	case src.ComponentEntity:
		resolvedComp, err := a.analyzeComponent(entity.Component, prog)
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
