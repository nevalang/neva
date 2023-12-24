package analyzer

import (
	"errors"
	"fmt"

	src "github.com/nevalang/neva/pkg/sourcecode"
)

var (
	ErrMainEntityNotFound       = errors.New("Main entity is not found")
	ErrMainEntityIsNotComponent = errors.New("Main entity is not a component")
	ErrMainEntityExported       = errors.New("Main entity cannot be exported")
	ErrMainPkgExports           = errors.New("Main package must cannot have exported entities")
)

func (a Analyzer) mainSpecificPkgValidation(mainPkgName string, mod src.Module, scope src.Scope) *Error {
	mainPkg := mod.Packages[mainPkgName]

	location := &src.Location{
		ModRef:  scope.Location.ModRef,
		PkgName: mainPkgName,
	}

	entityMain, filename, ok := mainPkg.Entity("Main")
	if !ok {
		return &Error{
			Err:      ErrMainEntityNotFound,
			Location: location,
		}
	}

	location.FileName = filename

	if entityMain.Kind != src.ComponentEntity {
		return &Error{
			Err:      ErrMainEntityIsNotComponent,
			Location: location,
		}
	}

	if entityMain.IsPublic {
		return &Error{
			Err:      ErrMainEntityExported,
			Location: location,
			Meta:     &entityMain.Component.Meta,
		}
	}

	scope = scope.WithLocation(*location)

	if err := a.analyzeMainComponent(entityMain.Component, mainPkg, scope); err != nil {
		return Error{
			Location: location,
			Meta:     &entityMain.Component.Meta,
		}.Merge(err)
	}

	if err := mainPkg.Entities(func(entity src.Entity, entityName, fileName string) error {
		if entity.IsPublic {
			return &Error{
				Err:      fmt.Errorf("%w: exported entity %v", ErrMainPkgExports, entityName),
				Meta:     entity.Meta(),
				Location: location,
			}
		}
		return nil
	}); err != nil {
		return err.(*Error) //nolint:forcetypeassert
	}

	return nil
}
