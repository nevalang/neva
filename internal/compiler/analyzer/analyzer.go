package analyzer

import (
	"errors"
	"fmt"

	"github.com/nevalang/neva/internal/compiler/src"
	"github.com/nevalang/neva/pkg/typesystem"
)

type Analyzer struct{}

var (
	ErrEmptyProgram      = errors.New("empty program")
	ErrMainPkgNotFound   = errors.New("main package not found")
	ErrEmptyPkg          = errors.New("package must not be empty")
	ErrUnknownEntityKind = errors.New("unknown entity kind")
)

func (a Analyzer) Analyze(prog src.Program) (src.Program, error) {
	if len(prog) == 0 {
		return nil, ErrEmptyProgram
	}

	mainPkg, ok := prog["main"]
	if !ok {
		return nil, ErrMainPkgNotFound
	}

	if err := a.mainSpecificPkgValidation(mainPkg, prog); err != nil {
		return nil, fmt.Errorf("main specific pkg validation: %w", err)
	}

	for name, pkg := range prog {
		if err := a.analyzePkg(pkg, prog); err != nil {
			return nil, fmt.Errorf("analyze pkg: %v: %w", name, err)
		}
	}

	return nil, nil
}

func (a Analyzer) analyzePkg(pkg src.Package, prog src.Program) error {
	if len(pkg) == 0 {
		return ErrEmptyPkg
	}

	if err := pkg.Entities(func(entity src.Entity, entityName, fileName string) error {
		if err := a.analyzeEntity(entityName, entity, pkg[fileName]); err != nil {
			return fmt.Errorf("analyze entity: %v: %v: %w", entityName, fileName, err)
		}
		return nil
	}); err != nil {
		return err
	}

	return nil
}

func (a Analyzer) analyzeEntity(entityName string, entity src.Entity, file src.File) error {
	// analyzedEntity := src.Entity{
	// 	Exported: entity.Exported,
	// 	Kind:     entity.Kind,
	// }

	switch entity.Kind {
	case src.TypeEntity:
		if err := a.analyzeTypeDef(entity.Type); err != nil {
			return fmt.Errorf("resolve type: %w", err)
		}
	// case src.ConstEntity:
	// 	analyzedConst, err := a.analyzeConst(entity.Const)
	// 	if err != nil {
	// 		return fmt.Errorf("analyze const: %w", err)
	// 	}
	// 	analyzedEntity.Const = analyzedConst
	// case src.InterfaceEntity:
	// 	analyzedEntity, err := a.analyzeInterface(entity.Interface)
	// 	if err != nil {
	// 		return fmt.Errorf("analyze interface: %w", err)
	// 	}
	// 	analyzedEntity.Interface = analyzedEntity
	// case src.ComponentEntity:
	// 	analyzedComponent, err := a.analyzeComponent(entity.Component)
	// 	if err != nil {
	// 		return fmt.Errorf("analyze component: %w", err)
	// 	}
	// 	analyzedEntity.Component = analyzedComponent
	default:
		return fmt.Errorf("%w: %v", ErrUnknownEntityKind, entity.Kind)
	}

	// file.Entities[entityName] = analyzedEntity

	return nil
}

var (
	ErrTypeDefWithoutBody = errors.New("type def without body")
)

func (a Analyzer) analyzeTypeDef(def typesystem.Def) error {
	if def.BodyExpr == nil {
		return ErrTypeDefWithoutBody
	}

	return nil
}
