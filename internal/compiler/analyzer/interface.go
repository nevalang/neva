package analyzer

import (
	"errors"

	src "github.com/nevalang/neva/internal/compiler/sourcecode"
	ts "github.com/nevalang/neva/pkg/typesystem"
)

var (
	ErrInterfaceTypeParams = errors.New("Cannot resolve interface type parameters")
	ErrEmptyInports        = errors.New("Interface must have inports")
	ErrEmptyOutports       = errors.New("Interface must have outports")
	ErrInvalidInports      = errors.New("Inports are invalid")
	ErrInvalidOutports     = errors.New("Outports are invalid")
)

type analyzeInterfaceParams struct {
	allowEmptyInports  bool
	allowEmptyOutports bool
}

func (a Analyzer) analyzeInterface(
	iface src.Interface,
	scope src.Scope,
	params analyzeInterfaceParams,
) (src.Interface, *Error) {
	resolvedParams, err := a.analyzeTypeParams(iface.TypeParams.Params, scope)
	if err != nil {
		return src.Interface{}, Error{
			Err:      ErrInterfaceTypeParams,
			Location: &scope.Location,
			Meta:     &iface.Meta,
		}.Merge(err)
	}

	resolvedIO, err := a.analyzeIO(resolvedParams, iface.IO, scope, params)
	if err != nil {
		return src.Interface{}, Error{
			Err:      ErrInterfaceTypeParams,
			Location: &scope.Location,
			Meta:     &iface.Meta,
		}.Merge(err)
	}

	typeParams := src.TypeParams{
		Params: resolvedParams,
		Meta:   iface.TypeParams.Meta,
	}

	return src.Interface{
		TypeParams: typeParams,
		IO:         resolvedIO,
	}, nil
}

func (a Analyzer) analyzeIO(
	typeParams []ts.Param,
	io src.IO,
	scope src.Scope,
	params analyzeInterfaceParams,
) (src.IO, *Error) {
	if !params.allowEmptyInports && len(io.In) == 0 {
		return src.IO{}, &Error{
			Err:      ErrEmptyInports,
			Location: &scope.Location,
		}
	}

	if !params.allowEmptyOutports && len(io.Out) == 0 {
		return src.IO{}, &Error{
			Err:      ErrEmptyOutports,
			Location: &scope.Location,
		}
	}

	resolvedIn, err := a.analyzePorts(typeParams, io.In, scope)
	if err != nil {
		return src.IO{}, Error{
			Err:      ErrInvalidInports,
			Location: &scope.Location,
		}.Merge(err)
	}

	resolvedOit, err := a.analyzePorts(typeParams, io.Out, scope)
	if err != nil {
		return src.IO{}, Error{
			Err:      ErrInvalidOutports,
			Location: &scope.Location,
		}.Merge(err)
	}

	return src.IO{
		In:  resolvedIn,
		Out: resolvedOit,
	}, nil
}

func (a Analyzer) analyzePorts(
	params []ts.Param,
	ports map[string]src.Port,
	scope src.Scope,
) (map[string]src.Port, *Error) {
	resolvedPorts := make(map[string]src.Port, len(ports))
	for name, port := range ports {
		resolvedPort, err := a.analyzePort(params, port, scope)
		if err != nil {
			return nil, Error{
				Location: &scope.Location,
				Meta:     &port.Meta,
			}.Merge(err)
		}
		resolvedPorts[name] = resolvedPort
	}
	return resolvedPorts, nil
}

func (a Analyzer) analyzePort(params []ts.Param, port src.Port, scope src.Scope) (src.Port, *Error) {
	resolvedDef, err := a.analyzeTypeDef(
		ts.Def{
			Params:   params,
			BodyExpr: &port.TypeExpr,
		},
		scope, analyzeTypeDefParams{allowEmptyBody: false},
	)
	if err != nil {
		return src.Port{}, Error{
			Location: &scope.Location,
			Meta:     &port.Meta,
		}.Merge(err)
	}

	return src.Port{
		TypeExpr: *resolvedDef.BodyExpr,
		IsArray:  port.IsArray,
	}, nil
}
