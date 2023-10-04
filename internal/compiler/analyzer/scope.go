package analyzer

import (
	"errors"
	"fmt"
	"strings"

	"github.com/nevalang/neva/internal/compiler/src"
	ts "github.com/nevalang/neva/pkg/typesystem"
)

type Scope struct {
	loc  Location
	prog src.Program
}

// Location is used by scope to resolve references.
type Location struct {
	pkg  string
	file string
}

func (s Scope) GetType(ref fmt.Stringer) (ts.Def, ts.Scope, error) {
	parsedRef, ok := ref.(src.EntityRef)
	if !ok {
		panic("not ok parse ref")
	}

	def, location, err := s.getType(parsedRef)
	if err != nil {
		return ts.Def{}, Scope{}, fmt.Errorf("get type: %w", err)
	}

	if parsedRef.Pkg == "" {
		return def, s, nil
	}

	return def, Scope{
		loc:  location,
		prog: s.prog,
	}, nil
}

// parse refer assumes refs in <pkg_name>.<entity_name≥ or just <entity_name> format
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
	entity, found, err := s.Entity(ref)
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

// getEntity returns entity by passed reference.
// If entity is local (ref has no pkg) the current location.pkg is used
// Otherwise we use current file imports to resolve external ref.
func (s Scope) Entity(entityRef src.EntityRef) (src.Entity, Location, error) {
	if entityRef.Pkg == "" {
		entity, filename, ok := s.prog[s.loc.pkg].Entity(entityRef.Name)
		if !ok {
			panic("")
		}
		return entity, Location{
			pkg:  s.loc.pkg,
			file: filename,
		}, nil
	}

	realImportPkgName, ok := s.prog[s.loc.pkg][s.loc.file].Imports[entityRef.Pkg]
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
		pkg:  realImportPkgName,
		file: fileName,
	}, nil
}
