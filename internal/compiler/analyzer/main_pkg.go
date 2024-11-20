package analyzer

import (
	"fmt"

	"github.com/nevalang/neva/internal/compiler"
	src "github.com/nevalang/neva/internal/compiler/sourcecode"
)

func (a Analyzer) mainSpecificPkgValidation(mainPkgName string, mod src.Module, scope src.Scope) *compiler.Error {
	mainPkg := mod.Packages[mainPkgName]

	location := src.Location{
		Module:  scope.Location().Module,
		Package: mainPkgName,
	}

	entityMain, filename, ok := mainPkg.Entity("Main")
	if !ok {
		return &compiler.Error{
			Message:  "Main entity is not found",
			Location: &location,
		}
	}

	location.Filename = filename

	if entityMain.Kind != src.ComponentEntity {
		return &compiler.Error{
			Message:  "Main entity must be a component",
			Location: &location,
		}
	}

	if entityMain.IsPublic {
		return &compiler.Error{
			Message:  "Main entity cannot be exported",
			Location: &location,
			Meta:     &entityMain.Component.Meta,
		}
	}

	scope = scope.Relocate(location)

	if err := a.analyzeMainComponent(entityMain.Component, scope); err != nil {
		return compiler.Error{
			Location: &location,
			Meta:     &entityMain.Component.Meta,
		}.Wrap(err)
	}

	for result := range mainPkg.Entities() {
		if result.Entity.IsPublic {
			return &compiler.Error{
				Message: fmt.Sprintf("Unexpected public entity in main package: %v", result.EntityName),
				Meta:    result.Entity.Meta(),
				Location: &src.Location{
					Module:   scope.Location().Module,
					Package:  mainPkgName,
					Filename: result.FileName,
				},
			}
		}
	}

	return nil
}
