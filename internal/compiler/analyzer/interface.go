package analyzer

import (
	"errors"
	"fmt"

	"github.com/nevalang/neva/internal/compiler/src"
	ts "github.com/nevalang/neva/pkg/typesystem"
)

var ErrInterfaceTypeParams = errors.New("could not resolve interface type parameters")

func (a Analyzer) analyzeInterface(def src.Interface) (src.Interface, error) {
	resolvedParams, err := a.analyzeTypeParams(def.TypeParams)
	if err != nil {
		return src.Interface{}, fmt.Errorf("%w: %v", ErrInterfaceTypeParams, def.TypeParams)
	}

	resolvedIO, err := a.analyzeIO(resolvedParams, def.IO)
	if err != nil {
		return src.Interface{}, fmt.Errorf("analyze IO: %w", err)
	}

	return src.Interface{
		TypeParams: resolvedParams,
		IO:         resolvedIO,
	}, nil
}

func (a Analyzer) analyzeIO(params []ts.Param, io src.IO) (src.IO, error) {
	resolvedIn, err := a.analyzePorts(params, io.In)
	if err != nil {
		return src.IO{}, fmt.Errorf("analyze inports: %w: %v", err, io.In)
	}

	resolvedOit, err := a.analyzePorts(params, io.Out)
	if err != nil {
		return src.IO{}, fmt.Errorf("analyze outports: %w: %v", err, io.In)
	}

	return src.IO{
		In:  resolvedIn,
		Out: resolvedOit,
	}, nil
}

func (a Analyzer) analyzePorts(params []ts.Param, ports map[string]src.Port) (map[string]src.Port, error) {
	resolvedPorts := make(map[string]src.Port, len(ports))
	for name, port := range ports {
		resolvedPort, err := a.analyzePort(params, port)
		if err != nil {
			return nil, fmt.Errorf("analyze port: %v: %w", name, err)
		}
		resolvedPorts[name] = resolvedPort
	}
	return resolvedPorts, nil
}

func (a Analyzer) analyzePort(params []ts.Param, port src.Port) (src.Port, error) {
	if port.TypeExpr == nil {
		return port, nil
	}

	resolvedDef, err := a.analyzeTypeDef(ts.Def{
		Params:   params,
		BodyExpr: port.TypeExpr,
	})
	if err != nil {
		return src.Port{}, fmt.Errorf("analyze type def: %w", err)
	}

	return src.Port{
		TypeExpr: resolvedDef.BodyExpr,
		IsArray:  port.IsArray,
	}, nil
}
