package main

import (
	"fmt"

	cprog "github.com/emil14/respect/internal/compiler/program"
	"github.com/emil14/respect/pkg/sdk"
)

type Caster interface {
	toSDK(cprog.Program) (sdk.Program, error)
	// fromSDK(sdk.Program) (cprog.Program, error)
}

type caster struct{}

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
	if _, ok := from.(cprog.Operator); ok {
		return sdk.Component{
			Io: c.castIO(from.Interface()),
		}, nil
	}

	mod, ok := from.(cprog.Module)
	if !ok {
		return sdk.Component{}, fmt.Errorf("casterr: unknown component type")
	}

	return sdk.Component{
		Io:      c.castIO(mod.Interface()),
		Workers: mod.Workers,
		Const:   c.castConst(mod.Const),
		Deps:    c.castDeps(mod.Deps),
		Net:     c.castNet(mod.Net),
	}, nil
}

func (c caster) castConst(from map[string]cprog.Const) map[string]sdk.Const {
	to := make(map[string]sdk.Const, len(from))
	for k, v := range from {
		to[k] = sdk.Const{
			Type:  v.Type.String(),
			Value: v.IntValue, // TMP
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

func (c caster) castNet(net cprog.Net) []sdk.Connection {
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
