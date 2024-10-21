package analyzer

import (
	"errors"
	"fmt"

	"github.com/nevalang/neva/internal/compiler"
	src "github.com/nevalang/neva/internal/compiler/sourcecode"
)

var (
	ErrMainEntityNotFound  = errors.New("Main entity is not found")
	ErrMainEntityIsNotFlow = errors.New("Main entity is not a component")
	ErrMainEntityExported  = errors.New("Main entity cannot be exported")
	ErrMainPkgExports      = errors.New("Main package must cannot have exported entities")
)

func (a Analyzer) mainSpecificPkgValidation(mainPkgName string, mod src.Module, scope src.Scope) *compiler.Error {
	mainPkg := mod.Packages[mainPkgName]

	location := &src.Location{
		ModRef:  scope.Location.ModRef,
		PkgName: mainPkgName,
	}

	entityMain, filename, ok := mainPkg.Entity("Main")
	if !ok {
		return &compiler.Error{
			Err:      ErrMainEntityNotFound,
			Location: location,
		}
	}

	location.FileName = filename

	if entityMain.Kind != src.ComponentEntity {
		return &compiler.Error{
			Err:      ErrMainEntityIsNotFlow,
			Location: location,
		}
	}

	if entityMain.IsPublic {
		return &compiler.Error{
			Err:      ErrMainEntityExported,
			Location: location,
			Range:    &entityMain.Component.Meta,
		}
	}

	scope = scope.WithLocation(*location)

	if err := a.analyzeMainComponent(entityMain.Component, scope); err != nil {
		return compiler.Error{
			Location: location,
			Range:    &entityMain.Component.Meta,
		}.Wrap(err)
	}

	if err := mainPkg.Entities(func(entity src.Entity, entityName, _ string) error {
		if entity.IsPublic {
			return &compiler.Error{
				Err:      fmt.Errorf("%w: exported entity %v", ErrMainPkgExports, entityName),
				Range:    entity.Meta(),
				Location: location,
			}
		}
		return nil
	}); err != nil {
		return err.(*compiler.Error) //nolint:forcetypeassert
	}

	return nil
}
