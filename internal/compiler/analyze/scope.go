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
	ErrNotImported            = errors.New("pkg not imported")
	ErrImportedEntityNotFound = errors.New("entity not found in imported package")
	ErrEntityNotExported      = errors.New("imported entity not exported")
)

type Scope struct {
	pkgs, imports   map[string]src.Pkg
	local, builtins map[string]src.Entity
	visited         map[src.EntityRef]struct{}
}

// Update will parse ref and, if it has pkg, call rebase with that pkg.
func (s Scope) Update(ref string) (ts.Scope, error) {
	pkg := s.parseRef(ref).Pkg
	if pkg == "" {
		return s, nil
	}
	return s.rebase(pkg)
}

func (s Scope) GetType(ref string) (ts.Def, error) {
	entity, err := s.getEntityByString(ref)
	if err != nil {
		return ts.Def{}, fmt.Errorf("%w: %v", ErrGetEntity, err)
	}

	if entity.Kind != src.TypeEntity {
		return ts.Def{}, fmt.Errorf("%w: want %v, got %v", ErrEntityKind, src.TypeEntity, entity.Kind)
	}

	return entity.Type, nil
}

func (s Scope) rebase(pkgName string) (Scope, error) {
	newBase, ok := s.imports[pkgName]
	if !ok {
		return Scope{}, fmt.Errorf("%w: %v", ErrNotImported, pkgName)
	}

	newImports, err := s.getImports(newBase.Imports)
	if err != nil {
		return Scope{}, fmt.Errorf("%w: %v", ErrImports, err)
	}

	s.imports = newImports
	s.local = newBase.Entities

	return s, nil
}

func (s Scope) getImports(pkgImports map[string]string) (map[string]src.Pkg, error) {
	imports := make(map[string]src.Pkg, len(pkgImports))
	for alias, pkgRef := range pkgImports {
		importedPkg, ok := s.pkgs[pkgRef]
		if !ok {
			return nil, fmt.Errorf("%w: %v", ErrImportNotFound, pkgRef)
		}
		imports[alias] = importedPkg
	}
	return imports, nil
}

func (s Scope) getMsg(ref src.EntityRef) (src.Msg, error) {
	entity, err := s.getEntity(ref)
	if err != nil {
		return src.Msg{}, fmt.Errorf("%w: %v", ErrGetEntity, err)
	}

	if entity.Kind != src.MsgEntity {
		return src.Msg{}, fmt.Errorf("%w: want %v, got %v", ErrEntityKind, src.TypeEntity, entity.Kind)
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

func (s Scope) getEntity(entityRef src.EntityRef) (src.Entity, error) {
	if entityRef.Pkg == "" {
		localDef, ok := s.local[entityRef.Name]
		if ok {
			return localDef, nil
		}

		builtinDef, ok := s.builtins[entityRef.Name]
		if !ok {
			return src.Entity{}, fmt.Errorf("%w: %v", ErrLocalOrBuiltinNotFound, entityRef.Name)
		}

		return builtinDef, nil
	}

	importedPkg, ok := s.imports[entityRef.Pkg]
	if !ok {
		return src.Entity{}, fmt.Errorf("%w: %v", ErrNoImport, entityRef.Pkg)
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
