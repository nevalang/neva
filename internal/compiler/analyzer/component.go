package analyzer

import (
	"github.com/nevalang/neva/internal/compiler"
	src "github.com/nevalang/neva/pkg/ast"
)

func (a Analyzer) analyzeComponent(
	componentName string,
	//nolint:gocritic // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
	component src.Component,
	//nolint:gocritic // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
	scope src.Scope,
) (src.Component, *compiler.Error) {
	if err := validateDuplicateDirectives(component.Directives); err != nil {
		return src.Component{}, err
	}

	_, hasExtern := component.Directives.Find(src.DirectiveKindExtern)

	resolvedIface, err := a.analyzeInterface(
		component.Interface,
		scope,
		analyzeInterfaceParams{
			allowEmptyInports:  hasExtern,
			allowEmptyOutports: hasExtern,
		},
	)
	if err != nil {
		return src.Component{}, compiler.Error{
			Meta: &component.Meta,
		}.Wrap(err)
	}

	if err := a.validateComponentDirectives(component, resolvedIface, scope, hasExtern); err != nil {
		return src.Component{}, err
	}

	if hasExtern {
		if len(component.Nodes) != 0 || len(component.Net) != 0 {
			return src.Component{}, &compiler.Error{
				Message: "Component with nodes or network cannot use #extern directive",
				Meta:    &component.Meta,
			}
		}
		return component, nil
	}

	resolvedNodes, nodesIfaces, hasGuard, err := a.analyzeNodes(
		componentName,
		resolvedIface,
		component.Nodes,
		component.Net,
		scope,
	)
	if err != nil {
		return src.Component{}, compiler.Error{
			Meta: &component.Meta,
		}.Wrap(err)
	}

	if len(component.Net) == 0 {
		return src.Component{}, &compiler.Error{
			Message: "Component must have network",
			Meta:    &component.Meta,
		}
	}

	analyzedNet, err := a.analyzeNetwork(
		component.Net,
		resolvedIface,
		hasGuard,
		resolvedNodes,
		nodesIfaces,
		scope,
	)
	if err != nil {
		return src.Component{}, compiler.Error{
			Meta: &component.Meta,
		}.Wrap(err)
	}

	return src.Component{
		Directives: component.Directives,
		Interface:  resolvedIface,
		Nodes:      resolvedNodes,
		Net:        analyzedNet,
		Meta:       component.Meta,
	}, nil
}

// validateComponentDirectives enforces declaration-level directive contracts
// after interface type parameters are resolved. The raw #bind_type expression
// remains in the AST so IR generation can substitute a call site's arguments.
//
//nolint:gocritic // Analyzer pipeline values are passed by value consistently.
func (a Analyzer) validateComponentDirectives(
	component src.Component,
	iface src.Interface,
	scope src.Scope,
	hasExtern bool,
) *compiler.Error {
	if component.Directives.Has(src.DirectiveKindBind) {
		return &compiler.Error{
			Message: "#bind directive is only valid on a component node",
			Meta:    &component.Meta,
		}
	}

	bindType, hasBindType := component.Directives.Find(src.DirectiveKindBindType)
	if !hasBindType {
		return nil
	}
	if !hasExtern {
		return &compiler.Error{
			Message: "#bind_type directive requires #extern",
			Meta:    &bindType.Meta,
		}
	}

	if _, err := a.resolver.ResolveExprWithFrame(
		bindType.BindType.TypeExpr,
		iface.TypeParams.ToFrame(),
		scope,
	); err != nil {
		return &compiler.Error{
			Message: err.Error(),
			Meta:    &bindType.Meta,
		}
	}

	return nil
}
