package analyzer

import (
	"errors"
	"fmt"
	"strings"

	"github.com/nevalang/neva/internal/compiler/src"
	ts "github.com/nevalang/neva/pkg/typesystem"
)

type Scope struct {
	location Location
	prog     src.Program
}

type Location struct {
	pkg  src.Package
	file src.File
}

func (s Scope) GetType(ref string) (ts.Def, ts.Scope, error) {
	parsedRef := s.parseRef(ref)

	def, found, err := s.getType(parsedRef)
	if err != nil {
		return ts.Def{}, Scope{}, fmt.Errorf("get type: %w", err)
	}

	if parsedRef.Pkg == "" {
		return def, s, nil
	}

	return def, Scope{
		location: found,
		prog:     s.prog,
	}, nil
}

func (s Scope) parseRef(ref string) src.EntityRef {
	var entityRef src.EntityRef

	parts := strings.Split(ref, ".")
	if len(parts) == 2 {
		entityRef.Pkg = parts[0]
		entityRef.Name = parts[1]
	} else {
		entityRef.Name = ref
	}

	return entityRef
}

var ErrEntityNotType = errors.New("entity not type")

func (s Scope) getType(ref src.EntityRef) (ts.Def, Location, error) {
	entity, found, err := s.getEntity(ref)
	if err != nil {
		return ts.Def{}, Location{}, err
	}

	if entity.Kind != src.TypeEntity {
		return ts.Def{}, Location{}, fmt.Errorf("%w: want %v, got %v", ErrEntityNotType, src.TypeEntity, entity.Kind)
	}

	return entity.Type, found, nil
}

var (
	ErrNoImport          = errors.New("no import found")
	ErrEntityNotExported = errors.New("entity is not exported")
)

func (s Scope) getEntity(entityRef src.EntityRef) (src.Entity, Location, error) {
	if entityRef.Pkg == "" {
		entity, filename, ok := s.location.pkg.Entity(entityRef.Name)
		if !ok {
			panic("")
		}
		return entity, Location{
			pkg:  s.location.pkg,
			file: s.location.pkg[filename],
		}, nil
	}

	realImportPkgName, ok := s.location.file.Imports[entityRef.Pkg]
	if !ok {
		return src.Entity{}, Location{}, fmt.Errorf("%w: %v", ErrNoImport, entityRef.Pkg)
	}

	entity, fileName, err := s.prog.Entity(src.EntityRef{
		Pkg:  realImportPkgName,
		Name: entityRef.Name,
	})
	if err != nil {
		panic(err)
	}

	if !entity.Exported {
		return src.Entity{}, Location{}, fmt.Errorf("%w: %v", ErrEntityNotExported, entityRef.Name)
	}

	return entity, Location{
		pkg:  s.prog[realImportPkgName],
		file: s.prog[realImportPkgName][fileName],
	}, nil
}
