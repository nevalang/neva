package analyze

import (
	"errors"
	"fmt"
	"strings"

	"github.com/emil14/neva/internal/compiler/src"
	ts "github.com/emil14/neva/pkg/types"
)

var (
	ErrGetEntity              = errors.New("can't get entity")
	ErrEntityKind             = errors.New("wrong entity kind")
	ErrNoImport               = errors.New("entity refers to not imported package")
	ErrLocalOrBuiltinNotFound = errors.New("local entity not found")
	ErrImports                = errors.New("can't build imports")
	ErrImportNotFound         = errors.New("imported package not found")
	ErrPkgNotFound            = errors.New("package not found")
	ErrNotImported            = errors.New("pkg not imported")
	ErrImportedEntityNotFound = errors.New("entity not found in imported package")
	ErrEntityNotExported      = errors.New("imported entity not exported")
	ErrRebase                 = errors.New("rebase")
)

// Scope implements types.Scope interface and some additional methods for analyzer
type Scope struct {
	base     string                     // Base must always refer to existing package in pkgs
	pkgs     map[string]src.Pkg         // Pkgs maps real names to all packages
	builtins map[string]src.Entity      // Second location for lookup for local entities
	visited  map[src.EntityRef]struct{} // Set of all visited entities
}

// Update parses ref and, if it has pkg, calls rebase
func (s Scope) Update(ref string) (ts.Scope, error) {
	pkgAlias := s.parseRef(ref).Pkg
	if pkgAlias == "" {
		return s, nil
	}

	scope, err := s.rebase(pkgAlias)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", errors.Join(ErrRebase, err), pkgAlias)
	}

	return scope, nil
}

func (s Scope) GetType(ref string) (ts.Def, error) {
	return s.getType(s.parseRef(ref))
}

func (s Scope) getLocalType(name string) (ts.Def, error) {
	return s.getType(src.EntityRef{Name: name})
}

func (s Scope) getType(ref src.EntityRef) (ts.Def, error) {
	entity, err := s.getEntity(ref)
	if err != nil {
		return ts.Def{}, err
	}

	if entity.Kind != src.TypeEntity {
		return ts.Def{}, fmt.Errorf("%w: want %v, got %v", ErrEntityKind, src.TypeEntity, entity.Kind)
	}

	return entity.Type, nil
}

func (s Scope) rebase(pkgAlias string) (Scope, error) {
	imports := s.pkgs[s.base].Imports // we assume s.base is valid

	pkgName, ok := imports[pkgAlias]
	if !ok {
		return Scope{}, fmt.Errorf("%w: %v", ErrImportNotFound, pkgAlias)
	}

	if _, ok = s.pkgs[pkgName]; !ok {
		return Scope{}, fmt.Errorf("%w: %v", ErrPkgNotFound, pkgAlias)
	}

	s.base = pkgName

	return s, nil
}

func (s Scope) getMsg(ref src.EntityRef) (src.Msg, error) {
	entity, err := s.getEntity(ref)
	if err != nil {
		return src.Msg{}, fmt.Errorf("%w: %v", ErrGetEntity, err)
	}

	if entity.Kind != src.MsgEntity {
		return src.Msg{}, fmt.Errorf("%w: want %v, got %v", ErrEntityKind, src.MsgEntity, entity.Kind)
	}

	return entity.Msg, nil
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

func (s Scope) getEntityByString(ref string) (src.Entity, error) {
	return s.getEntity(s.parseRef(ref))
}

func (s Scope) getLocalEntity(name string) (src.Entity, error) {
	return s.getEntity(src.EntityRef{Name: name})
}

func (s Scope) getEntity(entityRef src.EntityRef) (src.Entity, error) {
	basePkg := s.pkgs[s.base]

	if entityRef.Pkg == "" {
		localDef, ok := basePkg.Entities[entityRef.Name]
		if ok {
			s.visited[src.EntityRef{
				Pkg:  s.base,
				Name: entityRef.Name,
			}] = struct{}{}
			return localDef, nil
		}

		builtinDef, ok := s.builtins[entityRef.Name]
		if !ok {
			return src.Entity{}, fmt.Errorf("%w: %v", ErrLocalOrBuiltinNotFound, entityRef.Name)
		}
		s.visited[entityRef] = struct{}{}

		return builtinDef, nil
	}

	realImportedPkgName, ok := basePkg.Imports[entityRef.Pkg]
	if !ok {
		return src.Entity{}, fmt.Errorf("%w: %v", ErrNoImport, entityRef.Pkg)
	}

	importedPkg, ok := s.pkgs[realImportedPkgName]
	if !ok {
		return src.Entity{}, fmt.Errorf("%w: %v", ErrImportNotFound, realImportedPkgName)
	}

	importedEntity, ok := importedPkg.Entities[entityRef.Name]
	if !ok {
		return src.Entity{}, fmt.Errorf("%w: %v", ErrImportedEntityNotFound, entityRef.Name)
	}

	if !importedEntity.Exported {
		return src.Entity{}, fmt.Errorf("%w: %v", ErrEntityNotExported, entityRef.Name)
	}

	s.visited[entityRef] = struct{}{}

	return importedEntity, nil
}

func (s Scope) getComponent(entityRef src.EntityRef) (src.Component, error) {
	entity, err := s.getEntity(entityRef)
	if err != nil {
		return src.Component{}, nil
	}
	if entity.Kind != src.ComponentEntity {
		return src.Component{}, fmt.Errorf("%w: want %v, got %v", ErrEntityKind, src.ComponentEntity, entity.Kind)
	}
	return entity.Component, nil
}

func (s Scope) getInterface(entityRef src.EntityRef) (src.Interface, error) {
	entity, err := s.getEntity(entityRef)
	if err != nil {
		return src.Interface{}, nil
	}
	if entity.Kind != src.InterfaceEntity {
		return src.Interface{}, fmt.Errorf("%w: want %v, got %v", ErrEntityKind, src.InterfaceEntity, entity.Kind)
	}
	return entity.Interface, nil
}

// FIXME:
// pkg1 {
//     import pkg3 // <- unused import
//     E1
// }
// pkg2 {
//     e1 -> pkg3.e1 // makes pkg3 used import
// }
