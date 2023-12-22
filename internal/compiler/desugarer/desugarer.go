package desugarer

import (
	"errors"
	"fmt"
	"maps"

	"github.com/nevalang/neva/internal/compiler/analyzer"
	src "github.com/nevalang/neva/internal/compiler/sourcecode"
	ts "github.com/nevalang/neva/pkg/typesystem"
)

var ErrConstSenderEntityKind = errors.New("Entity that is used as a const reference in component's network must be of kind constant") //nolint:lll

type Desugarer struct{}

func (d Desugarer) Desugar(mod src.Module) (src.Module, error) {
	desugaredPkgs := make(map[string]src.Package, len(mod.Packages))

	for pkgName, pkg := range mod.Packages {
		desugaredPkg, err := d.desugarPkg(pkg, pkgName, mod)
		if err != nil {
			return src.Module{}, err
		}
		desugaredPkgs[pkgName] = desugaredPkg
	}

	return src.Module{
		Manifest: mod.Manifest,
		Packages: desugaredPkgs,
	}, nil
}

func (d Desugarer) desugarPkg(pkg src.Package, pkgName string, mod src.Module) (src.Package, error) {
	desugaredPkgs := make(src.Package, len(pkg))

	for fileName, file := range pkg {
		scope := src.Scope{
			Location: src.Location{
				ModuleName: "entry",
				PkgName:    pkgName,
				FileName:   fileName,
			},
			Module: mod,
		}

		desugaredFile, err := d.desugarFile(file, scope)
		if err != nil {
			return nil, err
		}

		desugaredPkgs[fileName] = desugaredFile
	}

	return desugaredPkgs, nil
}

func (d Desugarer) desugarFile(file src.File, scope src.Scope) (src.File, error) {
	desugaredEntities := make(map[string]src.Entity, len(file.Entities))

	for entityName, entity := range file.Entities {
		desugaredEntity, err := d.desugarEntity(entity, scope)
		if err != nil {
			return src.File{}, err
		}
		desugaredEntities[entityName] = desugaredEntity
	}

	return src.File{
		Imports:  file.Imports,
		Entities: desugaredEntities,
	}, nil
}

func (d Desugarer) desugarEntity(entity src.Entity, scope src.Scope) (src.Entity, error) {
	if entity.Kind != src.ComponentEntity {
		return entity, nil
	}

	desugarComponent, err := d.desugarComponent(entity.Component, scope)
	if err != nil {
		return src.Entity{}, err
	}

	return src.Entity{
		Exported:  entity.Exported,
		Kind:      entity.Kind,
		Component: desugarComponent,
	}, nil
}

func (d Desugarer) desugarComponent(component src.Component, scope src.Scope) (src.Component, error) {
	if len(component.Net) == 0 {
		return component, nil
	}

	desugaredNodes := make(map[string]src.Node, len(component.Nodes))
	maps.Copy(desugaredNodes, component.Nodes)
	desugaredNet := make([]src.Connection, 0, len(component.Net))

	for _, conn := range component.Net {
		if conn.SenderSide.ConstRef == nil {
			desugaredNet = append(desugaredNet, conn)
			continue
		}

		constTypeExpr, err := d.getConstType(*conn.SenderSide.ConstRef, scope)
		if err != nil {
			return src.Component{}, err
		}

		constRefStr := conn.SenderSide.ConstRef.String()

		desugaredNodes[constRefStr] = src.Node{
			Directives: map[src.Directive][]string{
				"runtime_func_msg": {constRefStr},
			},
			EntityRef: src.EntityRef{
				Pkg:  "std/builtin",
				Name: "Const",
			},
			TypeArgs: []ts.Expr{constTypeExpr},
		}

		constNodeOutportAddr := src.PortAddr{
			Node: constRefStr,
			Port: "out",
		}

		desugaredNet = append(desugaredNet, src.Connection{
			SenderSide: src.SenderConnectionSide{
				PortAddr:  &constNodeOutportAddr,
				Selectors: conn.SenderSide.Selectors,
				Meta:      conn.SenderSide.Meta,
			},
			ReceiverSides: conn.ReceiverSides,
			Meta:          conn.Meta,
		})
	}

	return src.Component{
		Directives: map[src.Directive][]string{},
		Interface:  src.Interface{},
		Nodes:      desugaredNodes,
		Net:        desugaredNet,
		Meta:       src.Meta{},
	}, nil
}

func (d Desugarer) getConstType(ref src.EntityRef, scope src.Scope) (ts.Expr, *analyzer.Error) {
	entity, _, err := scope.Entity(ref)
	if err != nil {
		return ts.Expr{}, &analyzer.Error{
			Err:      err,
			Location: &scope.Location,
			Meta:     &ref.Meta,
		}
	}

	if entity.Kind != src.ConstEntity {
		return ts.Expr{}, &analyzer.Error{
			Err:      fmt.Errorf("%w: %v", ErrConstSenderEntityKind, entity.Kind),
			Location: &scope.Location,
			Meta:     entity.Meta(),
		}
	}

	if entity.Const.Ref != nil {
		expr, err := d.getConstType(*entity.Const.Ref, scope)
		if err != nil {
			return ts.Expr{}, analyzer.Error{
				Location: &scope.Location,
				Meta:     entity.Meta(),
			}.Merge(err)
		}

		return expr, nil
	}

	return entity.Const.Value.TypeExpr, nil
}
