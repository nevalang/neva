package analyze

import (
	"errors"
	"fmt"
	"strings"

	"github.com/emil14/neva/internal/compiler/src"
	ts "github.com/emil14/neva/pkg/types"
)

var (
	ErrGetEntity  = errors.New("get entity")
	ErrEntityKind = errors.New("wrong entity kind")
)

type Scope struct {
	imports         map[string]src.Pkg
	local, builtins map[string]src.Entity
	visited         map[src.EntityRef]struct{}
}

// GetType implements types.Scope interface
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

func (s Scope) getEntityByString(ref string) (src.Entity, error) {
	var entityRef src.EntityRef

	parts := strings.Split(ref, ".")
	if len(parts) == 2 {
		entityRef.Pkg = parts[0]
		entityRef.Name = parts[1]
	} else {
		entityRef.Name = ref
	}

	return s.getEntity(entityRef)
}

func (s Scope) getEntity(entityRef src.EntityRef) (src.Entity, error) {
	if entityRef.Pkg == "" {
		localDef, ok := s.local[entityRef.Name]
		if ok {
			return localDef, nil
		}

		builtinDef, ok := s.builtins[entityRef.Name]
		if !ok {
			return src.Entity{}, errors.New("ref without package not found in local and builtin")
		}

		return builtinDef, nil
	}

	importedPkg, ok := s.imports[entityRef.Pkg]
	if !ok {
		return src.Entity{}, errors.New("referenced package is not found among imports")
	}

	importedEntity, ok := importedPkg.Entities[entityRef.Name]
	if !ok {
		return src.Entity{}, errors.New("referenced entity not found in imported package")
	}

	if !importedEntity.Exported {
		return src.Entity{}, errors.New("referenced entity not exported")
	}

	s.visited[entityRef] = struct{}{}

	return importedEntity, nil
}
