package analyzer

import (
	"errors"
	"fmt"

	"github.com/nevalang/neva/internal/compiler/src"
)

var (
	ErrMainEntityNotFound       = errors.New("entity main is not found")
	ErrMainEntityIsNotComponent = errors.New("main entity is not a component")
	ErrMainEntityExported       = errors.New("main entity is exported")
	ErrMainPkgExports           = errors.New("main pkg must not have exported entities")
)

func (a Analyzer) mainSpecificPkgValidation(pkg src.Package, pkgs map[string]src.Package) error {
	entityMain, _, ok := pkg.Entity("main")
	if !ok {
		return ErrMainEntityNotFound
	}

	if entityMain.Kind != src.ComponentEntity {
		return ErrMainEntityIsNotComponent
	}

	if entityMain.Exported {
		return ErrMainEntityExported
	}

	if err := a.analyzeMainComponent(entityMain.Component, pkg, pkgs); err != nil {
		return fmt.Errorf("analyze main component: %w", err)
	}

	// FIXME will conflict with general validation
	if err := pkg.Entities(func(entity src.Entity, entityName, fileName string) error {
		if entity.Exported {
			return fmt.Errorf("%w: file %v, entity %v", ErrMainPkgExports, fileName, entityName)
		}
		return nil
	}); err != nil {
		return fmt.Errorf("entities: %w", err)
	}

	return nil
}
