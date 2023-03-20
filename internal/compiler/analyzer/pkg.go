package analyzer

import (
	"errors"
	"fmt"

	"github.com/emil14/neva/internal/compiler"
)

var (
	ErrEntities                     = errors.New("analyze entities")
	ErrUsed                         = errors.New("analyze used")
	ErrExecutablePkg                = errors.New("analyze package with root component")
	ErrUselessPkg                   = errors.New("package without root component must have exports")
	ErrUnusedImport                 = errors.New("unused import")
	ErrUnusedEntity                 = errors.New("unused entity")
	ErrRootComponent                = errors.New("analyze root component")
	ErrRootComponentNotFound        = errors.New("root component not found")
	ErrRootComponentWrongEntityKind = errors.New("entity with name of the root component is not component")
	ErrExportedRootComponent        = errors.New("root component must not be exported")
)

// analyzePkg checks that:
// If pkg has ref to root component then it satisfies the pkg-with-root-component-specific requirements;
// There's no imports of not found pkgs;
// There's no unused imports;
// All entities are analyzed and;
// Used (exported or referenced by exported entities or root component).
func (a Analyzer) analyzePkg(pkgName string, pkgs map[string]compiler.Pkg) (compiler.Pkg, error) { //nolint:unparam
	pkg := pkgs[pkgName]

	if pkgName != "main" && len(a.getExports(pkg.Entities)) == 0 {
		return compiler.Pkg{}, ErrUselessPkg
	}

	scope := Scope{
		base:     pkgName,
		pkgs:     pkgs,
		builtins: a.builtinEntities(),
		visited:  map[compiler.EntityRef]struct{}{},
	}

	resolvedEntities, used, err := a.analyzeEntities(pkgName, pkg, scope)
	if err != nil {
		return compiler.Pkg{}, errors.Join(ErrEntities, err)
	}

	if err := a.analyzeUsed(pkgName, pkg, used); err != nil {
		return compiler.Pkg{}, errors.Join(ErrUsed, err)
	}

	return compiler.Pkg{
		Entities: resolvedEntities,
		Imports:  pkg.Imports,
	}, nil
}

// analyzeExecutablePkg checks that:
// Entity referenced as root component exist;
// That entity is a component;
// It's not exported and;
// It satisfies root-component-specific requirements;
func (a Analyzer) analyzeMainPkg(pkg compiler.Pkg, pkgs map[string]compiler.Pkg) error {
	entity, ok := pkg.Entities["main"]
	if !ok {
		return ErrRootComponentNotFound
	}

	if entity.Kind != compiler.ComponentEntity {
		return fmt.Errorf("%w: %v", ErrRootComponentWrongEntityKind, entity.Kind)
	}

	if entity.Exported {
		return ErrExportedRootComponent
	}

	if err := a.analyzeMainComponent(entity.Component, pkg, pkgs); err != nil {
		return fmt.Errorf("%w: %v", ErrRootComponent, err)
	}

	return nil
}

// getExports returns only exported entities
func (a Analyzer) getExports(entities map[string]compiler.Entity) map[string]compiler.Entity {
	exports := make(map[string]compiler.Entity, len(entities))
	for name, entity := range entities {
		exports[name] = entity
	}
	return exports
}

func (Analyzer) builtinEntities() map[string]compiler.Entity {
	return map[string]compiler.Entity{
		"int": h.BaseTypeEntity(),
		"vec": h.BaseTypeEntity(h.ParamWithNoConstr("t")),
	}
}

// analyzeUsed returns error if there're unused imports or entities
func (Analyzer) analyzeUsed(pkgName string, pkg compiler.Pkg, usedEntities map[compiler.EntityRef]struct{}) error {
	usedImports := map[string]struct{}{}
	usedLocalEntities := map[string]struct{}{}

	for ref := range usedEntities {
		if ref.Pkg == pkgName {
			usedLocalEntities[ref.Name] = struct{}{}
		} else {
			usedImports[ref.Pkg] = struct{}{}
		}
	}

	for importAlias := range pkg.Imports {
		if _, ok := usedImports[importAlias]; !ok {
			return fmt.Errorf("%w: %v", ErrUnusedImport, importAlias)
		}
	}

	for entityName := range pkg.Entities {
		if _, ok := usedLocalEntities[entityName]; !ok {
			return fmt.Errorf("%w: %v", ErrUnusedEntity, entityName)
		}
	}

	return nil
}
