package analyzer

import (
	"errors"
	"fmt"

	src "github.com/nevalang/neva/internal/compiler/sourcecode"
	ts "github.com/nevalang/neva/pkg/typesystem"
)

var ErrInterfaceTypeParams = errors.New("could not resolve interface type parameters")

type analyzeInterfaceParams struct {
	allowEmptyInports  bool
	allowEmptyOutports bool
}

func (a Analyzer) analyzeInterface(
	def src.Interface,
	scope src.Scope,
	params analyzeInterfaceParams,
) (src.Interface, error) {
	resolvedParams, err := a.analyzeTypeParams(def.TypeParams.Params, scope)
	if err != nil {
		return src.Interface{}, fmt.Errorf("%w: %v", ErrInterfaceTypeParams, def.TypeParams)
	}

	resolvedIO, err := a.analyzeIO(resolvedParams, def.IO, scope, params)
	if err != nil {
		return src.Interface{}, fmt.Errorf("analyze IO: %w", err)
	}

	typeParams := src.TypeParams{
		Params: resolvedParams,
		Meta:   def.TypeParams.Meta,
	}

	return src.Interface{
		TypeParams: typeParams,
		IO:         resolvedIO,
	}, nil
}

var (
	ErrEmptyInports  = errors.New("IO cannot have empty inports")
	ErrEmptyOutports = errors.New("IO cannot have empty outports")
)

func (a Analyzer) analyzeIO(
	typeParams []ts.Param,
	io src.IO,
	scope src.Scope,
	params analyzeInterfaceParams,
) (src.IO, error) {
	if !params.allowEmptyInports && len(io.In) == 0 {
		return src.IO{}, ErrEmptyInports
	}

	if !params.allowEmptyOutports && len(io.Out) == 0 {
		return src.IO{}, ErrEmptyOutports
	}

	resolvedIn, err := a.analyzePorts(typeParams, io.In, scope)
	if err != nil {
		return src.IO{}, fmt.Errorf("analyze inports: %w: %v", err, io.In)
	}

	resolvedOit, err := a.analyzePorts(typeParams, io.Out, scope)
	if err != nil {
		return src.IO{}, fmt.Errorf("analyze outports: %w: %v", err, io.In)
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
) (map[string]src.Port, error) {
	resolvedPorts := make(map[string]src.Port, len(ports))
	for name, port := range ports {
		resolvedPort, err := a.analyzePort(params, port, scope)
		if err != nil {
			return nil, fmt.Errorf("analyze port: %v: %w", name, err)
		}
		resolvedPorts[name] = resolvedPort
	}
	return resolvedPorts, nil
}

func (a Analyzer) analyzePort(params []ts.Param, port src.Port, scope src.Scope) (src.Port, error) {
	resolvedDef, err := a.analyzeTypeDef(ts.Def{
		Params:   params,
		BodyExpr: &port.TypeExpr,
	}, scope, analyzeTypeDefParams{allowEmptyBody: false})
	if err != nil {
		return src.Port{}, fmt.Errorf("analyze type def: %w", err)
	}
	return src.Port{
		TypeExpr: *resolvedDef.BodyExpr,
		IsArray:  port.IsArray,
	}, nil
}
