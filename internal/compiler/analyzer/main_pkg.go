package analyzer

import (
	"errors"
	"fmt"

	src "github.com/nevalang/neva/internal/compiler/sourcecode"
)

var (
	ErrMainEntityNotFound       = errors.New("Main entity is not found")
	ErrMainEntityIsNotComponent = errors.New("Main entity is not a component")
	ErrMainEntityExported       = errors.New("Main entity cannot be exported")
	ErrMainPkgExports           = errors.New("Main package must cannot have exported entities")
)

//nolint:funlen
func (a Analyzer) mainSpecificPkgValidation(mainPkgName string, mod src.Module) *Error {
	mainPkg := mod.Packages[mainPkgName]

	fallbackLocation := &src.Location{
		ModuleName: "entry",
		PkgName:    mainPkgName,
	}

	entityMain, filename, ok := mainPkg.Entity("Main")
	if !ok {
		return &Error{
			Err:      ErrMainEntityNotFound,
			Location: fallbackLocation,
		}
	}

	fallbackLocation.FileName = filename

	if entityMain.Kind != src.ComponentEntity {
		return &Error{
			Err:      ErrMainEntityIsNotComponent,
			Location: fallbackLocation,
		}
	}

	if entityMain.Exported {
		return &Error{
			Err:      ErrMainEntityExported,
			Location: fallbackLocation,
			Meta:     &entityMain.Component.Meta,
		}
	}

	scope := src.Scope{
		Module: mod,
		Location: src.Location{
			PkgName:  mainPkgName,
			FileName: filename,
		},
	}

	if err := a.analyzeMainComponent(entityMain.Component, mainPkg, scope); err != nil {
		return Error{
			Location: fallbackLocation,
			Meta:     &entityMain.Component.Meta,
		}.Merge(err)
	}

	if err := mainPkg.Entities(func(entity src.Entity, entityName, fileName string) error {
		if entity.Exported {
			var meta src.Meta
			if m, err := entity.Meta(); err != nil {
				meta = m
			}
			return &Error{
				Err:  fmt.Errorf("%w: exported entity %v", ErrMainPkgExports, entityName),
				Meta: &meta,
				Location: &src.Location{
					ModuleName: "entry",
					PkgName:    mainPkgName,
					FileName:   filename,
				},
			}
		}
		return nil
	}); err != nil {
		return err.(*Error) //nolint:forcetypeassert
	}

	return nil
}
