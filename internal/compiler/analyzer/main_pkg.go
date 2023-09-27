package analyzer

import (
	"fmt"

	"github.com/nevalang/neva/internal/compiler/src"
)

func (a Analyzer) mainSpecificPkgValidation(pkg src.Package, pkgs map[string]src.Package) error {
	entityMain, ok := pkg.Entity("main")
	if !ok {
		panic("analyzer: no main entity")
	}

	if entityMain.Kind != src.ComponentEntity {
		panic("analyzer: main entity is not a component")
	}

	if entityMain.Exported {
		panic("analyzer: main entity is exported")
	}

	if err := a.analyzeMainComponent(entityMain.Component, pkg, pkgs); err != nil {
		panic(fmt.Errorf("analyzer: %w", err))
	}

	return nil
}
