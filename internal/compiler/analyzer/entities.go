package analyzer

import (
	"errors"
	"fmt"

	"github.com/nevalang/neva/internal/compiler"
)

var (
	ErrEntity              = errors.New("analyze entity")
	ErrScopeGetLocalEntity = errors.New("scope get local entity")
	ErrType                = errors.New("analyze type")
	ErrUnknownEntityKind   = errors.New("unknown entity kind")
	ErrMsg                 = errors.New("analyze msg")
	ErrComponent           = errors.New("analyze component")
	ErrInterface           = errors.New("analyze interface")
)

func (a Analyzer) analyzeEntities(
	pkgName string,
	pkg compiler.Pkg,
	scope Scope,
) (
	map[string]compiler.Entity,
	map[compiler.EntityRef]struct{},
	error,
) {
	resolvedEntities := make(map[string]compiler.Entity, len(pkg.Entities))
	used := map[compiler.EntityRef]struct{}{} // both local and imported

	for entityName, entity := range pkg.Entities {
		if entityName == "main" && pkgName != "main" {
			panic("main entity inside not main package")
		}

		if entity.Exported || entityName == "main" {
			used[compiler.EntityRef{
				Pkg:  pkgName,
				Name: entityName,
			}] = struct{}{}
		}

		resolvedEntity, entitiesUsedByEntity, err := a.analyzeEntity(entityName, scope)
		if err != nil {
			return nil, nil, fmt.Errorf("%w: %v", errors.Join(ErrEntity, err), entityName)
		}

		for entityRef := range entitiesUsedByEntity {
			used[entityRef] = struct{}{}
		}

		resolvedEntities[entityName] = resolvedEntity
	}

	return resolvedEntities, used, nil
}

func (a Analyzer) analyzeEntity(name string, scope Scope) (compiler.Entity, map[compiler.EntityRef]struct{}, error) {
	entity, err := scope.getLocalEntity(name)
	if err != nil {
		return compiler.Entity{}, nil, errors.Join(ErrScopeGetLocalEntity, err)
	}

	switch entity.Kind { // https://github.com/nevalang/neva/issues/186
	case compiler.TypeEntity:
		resolvedDef, usedTypeEntities, err := a.analyzeType(name, scope)
		if err != nil {
			return compiler.Entity{}, nil, errors.Join(ErrType, err)
		}
		return compiler.Entity{
			Type:     resolvedDef,
			Kind:     compiler.TypeEntity,
			Exported: entity.Exported,
		}, usedTypeEntities, nil
	case compiler.MsgEntity:
		resolvedMsg, usedEntities, err := a.analyzeMsg(entity.Msg, scope, nil)
		if err != nil {
			return compiler.Entity{}, nil, errors.Join(ErrMsg, err)
		}
		return compiler.Entity{
			Msg:      resolvedMsg,
			Kind:     compiler.MsgEntity,
			Exported: entity.Exported,
		}, usedEntities, nil
	case compiler.InterfaceEntity:
		resolvedInterface, used, err := a.analyzeInterface(entity.Interface, scope)
		if err != nil {
			return compiler.Entity{}, nil, errors.Join(ErrInterface, err)
		}
		return compiler.Entity{
			Exported:  entity.Exported,
			Kind:      compiler.InterfaceEntity,
			Interface: resolvedInterface,
		}, used, nil
	case compiler.ComponentEntity:
		resolvedComponent, used, err := a.analyzeCmp(entity.Component, scope)
		if err != nil {
			return compiler.Entity{}, nil, errors.Join(ErrComponent, err)
		}
		return compiler.Entity{
			Exported:  entity.Exported,
			Kind:      compiler.ComponentEntity,
			Component: resolvedComponent,
		}, used, nil
	default:
		return compiler.Entity{}, nil, ErrUnknownEntityKind
	}
}
