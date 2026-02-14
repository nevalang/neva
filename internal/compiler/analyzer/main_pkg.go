package analyzer

import (
	"fmt"

	"github.com/nevalang/neva/internal/compiler"
	src "github.com/nevalang/neva/pkg/ast"
	"github.com/nevalang/neva/pkg/core"
)

func (a Analyzer) mainSpecificPkgValidation(mainPkgName string, mod src.Module, scope src.Scope) *compiler.Error {
	mainPkg := mod.Packages[mainPkgName]

	location := core.Location{
		ModRef:  scope.Location().ModRef,
		Package: mainPkgName,
	}

	entityMain, filename, ok := mainPkg.Entity("Main")
	if !ok {
		return &compiler.Error{
			Message: "Main entity is not found",
			Meta: &core.Meta{
				Location: location,
			},
		}
	}

	location.Filename = filename

	if entityMain.Kind != src.ComponentEntity {
		return &compiler.Error{
			Message: "Main entity must be a component",
			Meta: &core.Meta{
				Location: location,
			},
		}
	}

	if entityMain.IsPublic {
		return &compiler.Error{
			Message: "Main entity cannot be exported",
			Meta: &core.Meta{
				Location: location,
			},
		}
	}

	scope = scope.Relocate(location)

	if len(entityMain.Component) != 1 {
		return &compiler.Error{
			Message: "Main entity must have non-overloaded component",
			Meta: &core.Meta{
				Location: location,
			},
		}
	}

	if err := a.analyzeMainComponent(entityMain.Component[0], scope); err != nil {
		return compiler.Error{
			Meta: &entityMain.Component[0].Meta,
		}.Wrap(err)
	}

	for result := range mainPkg.Entities() {
		if result.Entity.IsPublic {
			return &compiler.Error{
				Message: fmt.Sprintf("Unexpected public entity in main package: %v", result.EntityName),
				Meta:    result.Entity.Meta(),
			}
		}
	}

	return nil
}
