package analyzer

import (
	"errors"
	"fmt"

	"github.com/nevalang/neva/internal/compiler/src"
	"github.com/nevalang/neva/pkg/typesystem"
	ts "github.com/nevalang/neva/pkg/typesystem"
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

// AnalyzeAndResolve returns error if program is invalid. It also modifies program by resolving types.
func (a Analyzer) AnalyzeAndResolve(prog src.Program) error {
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
		resolvedEntity, err := a.analyzeEntity(entityName, entity, pkg[fileName])
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

func (a Analyzer) analyzeEntity(entityName string, entity src.Entity, file src.File) (src.Entity, error) {
	resolvedEntity := src.Entity{
		Exported: entity.Exported,
		Kind:     entity.Kind,
	}

	switch entity.Kind {
	case src.TypeEntity:
		resolvedTypeDef, err := a.analyzeTypeDef(entity.Type)
		if err != nil {
			return src.Entity{}, fmt.Errorf("resolve type: %w", err)
		}
		resolvedEntity.Type = resolvedTypeDef
	case src.ConstEntity:
		resolvedConst, err := a.analyzeConst(entity.Const)
		if err != nil {
			return src.Entity{}, fmt.Errorf("analyze const: %w", err)
		}
		resolvedEntity.Const = resolvedConst
	case src.InterfaceEntity:
		resolvedInterface, err := a.analyzeInterface(entity.Interface)
		if err != nil {
			return src.Entity{}, fmt.Errorf("analyze interface: %w", err)
		}
		resolvedEntity.Interface = resolvedInterface
	case src.ComponentEntity:
		if err := a.analyzeComponent(entity.Component); err != nil {
			return src.Entity{}, fmt.Errorf("analyze component: %w", err)
		}
	default:
		return src.Entity{}, fmt.Errorf("%w: %v", ErrUnknownEntityKind, entity.Kind)
	}

	return resolvedEntity, nil
}

var ErrCustomBaseType = errors.New("custom type must have body expression and cannot be used for recursive definitions")

func (a Analyzer) analyzeTypeDef(def typesystem.Def) (typesystem.Def, error) {
	return def, nil // TODO
}

// TODO unused
func (Analyzer) buildTestExprArgs(params []ts.Param) []ts.Expr {
	args := make([]ts.Expr, 0, len(params))
	for _, param := range params {
		if param.Constr == nil {
			args = append(args, ts.Expr{
				Inst: &ts.InstExpr{Ref: "any"},
			})
		} else {
			args = append(args, *param.Constr)
		}
	}
	return args
}

// FIXME constr_refereing_type_parameter_(generics_inside_generics)" t<int, vec<int>> {t<a, b vec<a>>, vec<t>, int}
func (a Analyzer) resolveTypeParams(params []ts.Param) ([]ts.Param, error) {
	resolvedParams := make([]ts.Param, 0, len(params))
	for _, param := range params {
		if param.Constr == nil {
			resolvedParams = append(resolvedParams, param)
			continue
		}
		resolvedParam, err := a.resolveTypeExpr(*param.Constr)
		if err != nil {
			return nil, fmt.Errorf("analyze type expr: %w", err)
		}
		resolvedParams = append(resolvedParams, ts.Param{
			Name:   param.Name,
			Constr: &resolvedParam,
		})
	}
	return resolvedParams, nil
}

func (a Analyzer) resolveTypeExpr(expr typesystem.Expr) (typesystem.Expr, error) {
	a.resolver.Resolve(expr, nil)
	return expr, nil
}

func (a Analyzer) analyzeComponent(def src.Component) error {
	return nil
}
