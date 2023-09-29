package analyzer

import (
	"errors"
	"fmt"

	"github.com/nevalang/neva/internal/compiler/src"
	"github.com/nevalang/neva/pkg/typesystem"
)

type Analyzer struct {
	resolver typesystem.Resolver
}

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
	switch entity.Kind {
	case src.TypeEntity:
		if err := a.analyzeTypeDef(entity.Type); err != nil {
			return fmt.Errorf("resolve type: %w", err)
		}
	case src.ConstEntity:
		if err := a.analyzeConst(entity.Const); err != nil {
			return fmt.Errorf("analyze const: %w", err)
		}
	case src.InterfaceEntity:
		if err := a.analyzeInterface(entity.Interface); err != nil {
			return fmt.Errorf("analyze interface: %w", err)
		}
	case src.ComponentEntity:
		if err := a.analyzeComponent(entity.Component); err != nil {
			return fmt.Errorf("analyze component: %w", err)
		}
	default:
		return fmt.Errorf("%w: %v", ErrUnknownEntityKind, entity.Kind)
	}
	return nil
}

var ErrCustomBaseType = errors.New("custom type must have body expression and cannot be used for recursive definitions")

func (a Analyzer) analyzeTypeDef(def typesystem.Def) error {
	// We check this here because it's ok for type-system to face base types.
	if def.BodyExpr == nil || def.CanBeUsedForRecursiveDefinitions { // FIXME will conflict with actual base types
		return ErrCustomBaseType
	}
	return nil
}

func (a Analyzer) resolveTypeExpr(expr typesystem.Expr) (typesystem.Expr, error) {
	return expr, nil // TODO
}

var (
	ErrEmptyConst                  = errors.New("const must have value or reference to another const")
	ErrResolveConstType            = errors.New("cannot resolve constant type")
	ErrConstValuesOfDifferentTypes = errors.New("constant cannot have values of different types at once")
)

func (a Analyzer) analyzeConst(constant src.Const) error {
	if constant.Value == nil && constant.Ref == nil {
		return ErrEmptyConst
	}

	if constant.Value == nil {
		panic("references for constants not implemented yet")
	}

	resolvedType, err := a.resolveTypeExpr(constant.Value.TypeExpr)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrResolveConstType, err)
	}

	switch resolvedType.Inst.Ref {
	case "bool":
		if constant.Value.Int != 0 || constant.Value.Float != 0 || constant.Value.Str != "" {
			return fmt.Errorf("%w: %v", ErrConstValuesOfDifferentTypes, constant.Value)
		}
	case "int":
		if constant.Value.Bool != false || constant.Value.Float != 0 || constant.Value.Str != "" {
			return fmt.Errorf("%w: %v", ErrConstValuesOfDifferentTypes, constant.Value)
		}
	case "float":
		if constant.Value.Bool != false || constant.Value.Int != 0 || constant.Value.Str != "" {
			return fmt.Errorf("%w: %v", ErrConstValuesOfDifferentTypes, constant.Value)
		}
	case "str":
		if constant.Value.Bool != false || constant.Value.Int != 0 || constant.Value.Float != 0 {
			return fmt.Errorf("%w: %v", ErrConstValuesOfDifferentTypes, constant.Value)
		}
	}

	return nil
}

func (a Analyzer) analyzeInterface(def src.Interface) error {
	return nil
}

func (a Analyzer) analyzeComponent(def src.Component) error {
	return nil
}
