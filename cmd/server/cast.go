package main

import (
	"fmt"

	"github.com/emil14/neva/internal/compiler/program"
	cprog "github.com/emil14/neva/internal/compiler/program"
	"github.com/emil14/neva/pkg/sdk"
)

type Caster interface {
	toSDK(cprog.Program) (sdk.Program, error)
}

type caster struct{}

func (c caster) toOperators(from map[string]cprog.Operator) map[string]sdk.Operator {
	to := make(map[string]sdk.Operator, len(from))
	for k, v := range from {
		to[k] = sdk.Operator{
			Io: c.castIO(v.IO),
		}
	}
	return to
}

func (c caster) toSDK(from cprog.Program) (sdk.Program, error) {
	cc, err := c.castComponents(from.Scope)
	if err != nil {
		return sdk.Program{}, err
	}
	return sdk.Program{
		Scope: cc,
		Root:  from.Root,
	}, nil
}

func (c caster) castComponents(from map[string]cprog.Component) (map[string]sdk.Component, error) {
	r := map[string]sdk.Component{}
	for k, v := range from {
		cmpnt, err := c.castComponent(v)
		if err != nil {
			return nil, err
		}
		r[k] = cmpnt
	}
	return r, nil
}

func (c caster) castComponent(from cprog.Component) (sdk.Component, error) {
	if from.Type == program.OperatorComponent {
		return sdk.Component{
			Io: c.castIO(from.IO()),
		}, nil
	}

	if from.Type != program.ModuleComponent {
		return sdk.Component{}, fmt.Errorf("unknown component type %d", from.Type)
	}

	return sdk.Component{
		Io:      c.castIO(from.Module.Interface()),
		Workers: from.Module.Workers,
		Const:   c.castConst(from.Module.Const),
		Deps:    c.castDeps(from.Module.DepsIO),
		Net:     c.castNet(from.Module.Net),
	}, nil
}

func (c caster) castConst(from map[string]cprog.Const) map[string]sdk.Const {
	to := make(map[string]sdk.Const, len(from))
	for k, v := range from {
		to[k] = sdk.Const{
			// Type:  v.Type,
			Value: v.Int(), // TMP
		}
	}
	return to
}

func (c caster) castDeps(from map[string]cprog.IO) map[string]sdk.Io {
	r := map[string]sdk.Io{}
	for k, v := range from {
		r[k] = c.castIO(v)
	}
	return r
}

func (c caster) castNet(net cprog.Connections) []sdk.Connection {
	r := make([]sdk.Connection, 0, len(net))
	for from, to := range net {
		for rcvr := range to {
			r = append(r, c.sdkConnection(from, rcvr))
		}
	}
	return r
}

func (c caster) sdkConnection(from, to cprog.PortAddr) sdk.Connection {
	return sdk.Connection{
		From: c.sdkPortAddr(from),
		To:   c.sdkPortAddr(to),
	}
}

func (c caster) sdkPortAddr(from cprog.PortAddr) sdk.PortAddr {
	return sdk.PortAddr{
		Node: from.Node,
		Idx:  int32(from.Idx),
		Port: from.Port,
	}
}

func (c caster) castIO(from cprog.IO) sdk.Io {
	return sdk.Io{
		In:  c.castPorts(from.In),
		Out: c.castPorts(from.Out),
	}
}

func (c caster) castPorts(from cprog.Ports) map[string]string {
	to := make(map[string]string, len(from))
	for name, typ := range from {
		to[name] = typ.String()
	}
	return to
}
