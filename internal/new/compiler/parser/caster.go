package parser

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/emil14/neva/internal/new/compiler"
)

var (
	ErrNet           = errors.New("net")
	ErrPortAddr      = errors.New("port addr")
	ErrPortAddrParts = errors.New("string must consist of node and port separated by dot")
)

type caster struct{}

func (c caster) Cast(mod Module) (compiler.Module, error) {
	cnet, err := c.castNet(mod.Net)
	if err != nil {
		return compiler.Module{}, fmt.Errorf("%w: %v", ErrNet, err)
	}

	return compiler.Module{
		IO:     c.castIO(mod.IO),
		DepsIO: c.castDeps(mod.Deps),
		Nodes:  c.castNodes(mod.Nodes),
		Net:    cnet,
		Meta: compiler.ModuleMeta{
			WantCompilerVersion: mod.Meta.Compiler,
		},
	}, nil
}

func (c caster) castDeps(deps map[string]IO) map[string]compiler.IO {
	cdeps := make(map[string]compiler.IO, len(deps))
	for name, io := range deps {
		cdeps[name] = c.castIO(io)
	}
	return cdeps
}

func (c caster) castIO(io IO) compiler.IO {
	return compiler.IO{
		In:  c.castPorts(io.In),
		Out: c.castPorts(io.Out),
	}
}

func (c caster) castPorts(ports map[string]string) map[string]compiler.Port {
	cports := make(map[string]compiler.Port, len(ports))
	for name, typ := range ports {
		cports[name] = compiler.Port{
			Type:    c.castPortType(name),
			MsgType: c.castMsgType(typ),
		}
	}
	return cports
}

func (c caster) castPortType(typ string) compiler.PortType {
	if strings.HasSuffix(typ, "[]") {
		return compiler.ArrPortType
	}
	return compiler.NormPortType
}

func (c caster) castMsgType(typ string) compiler.MsgType {
	switch typ {
	case "int":
		return compiler.IntMsgType
	case "str":
		return compiler.StrMsgType
	case "bool":
		return compiler.BoolMsgType
	case "sig":
		return compiler.SigMsgType
	}
	return compiler.UnknownMsgType
}

func (c caster) castNodes(nodes Nodes) compiler.ModuleNodes {
	return compiler.ModuleNodes{
		Const:   c.castConstNode(nodes.Const),
		Workers: nodes.Workers,
	}
}

func (c caster) castConstNode(constOutPorts map[string]ConstOutPort) map[string]compiler.Msg {
	constNode := make(map[string]compiler.Msg, len(constOutPorts))
	for outPortName, msg := range constOutPorts {
		constNode[outPortName] = c.castMsg(msg)
	}
	return constNode
}

func (c caster) castMsg(msg ConstOutPort) compiler.Msg {
	switch msg.Type {
	case "str":
		return compiler.Msg{
			Type:     compiler.StrMsgType,
			StrValue: msg.Str,
		}
	case "int":
		return compiler.Msg{
			Type:     compiler.IntMsgType,
			IntValue: msg.Int,
		}
	case "bool":
		return compiler.Msg{
			Type:      compiler.IntMsgType,
			BoolValue: msg.Bool,
		}
	case "sig":
		return compiler.Msg{
			Type: compiler.IntMsgType,
		}
	}
	return compiler.Msg{
		Type: compiler.UnknownMsgType,
	}
}

func (c caster) castNet(net map[string][]string) ([]compiler.Connection, error) {
	cnet := make([]compiler.Connection, 0, len(net))

	for from, to := range net {
		cfrom, err := c.castPortAddr(from)
		if err != nil {
			return nil, fmt.Errorf("%w: from: %v", ErrPortAddr, err)
		}

		cto, err := c.castPortAddrs(to)
		if err != nil {
			return nil, fmt.Errorf("to: %w", err)
		}

		cnet = append(cnet, compiler.Connection{
			From: cfrom,
			To:   cto,
		})
	}

	return cnet, nil
}

func (c caster) castPortAddrs(addrs []string) ([]compiler.PortAddr, error) {
	caddrs := make([]compiler.PortAddr, 0, len(addrs))

	for _, addr := range addrs {
		caddr, err := c.castPortAddr(addr)
		if err != nil {
			return nil, fmt.Errorf("%w: %v", ErrPortAddr, err)
		}

		caddrs = append(caddrs, caddr)
	}

	return caddrs, nil
}

func (c caster) castPortAddr(addr string) (compiler.PortAddr, error) {
	parts := strings.Split(addr, ".")
	if len(parts) != 2 {
		return compiler.PortAddr{}, ErrPortAddrParts
	}

	node := parts[0]
	port := parts[1]

	bracketStart := strings.Index(port, "[")
	if bracketStart == -1 {
		return compiler.PortAddr{
			Node: node,
			Port: port,
		}, nil
	}

	bracketEnd := strings.Index(port, "]")
	if bracketEnd == -1 {
		return compiler.PortAddr{
			Node: node,
			Port: port,
		}, nil
	}

	strIdx := port[bracketStart+1 : bracketEnd]
	if strIdx == ":" {
		return compiler.PortAddr{
			Type: compiler.ArrByPassPortAddr,
			Node: node,
			Port: port[:bracketStart],
		}, nil
	}

	numIdx, err := strconv.ParseUint(port[bracketStart+1:bracketEnd], 10, 64)
	if err != nil {
		return compiler.PortAddr{
			Node: node,
			Port: port,
		}, nil
	}

	return compiler.PortAddr{
		Type: compiler.NormPortAddr,
		Node: node,
		Port: port[:bracketStart],
		Idx:  uint8(numIdx),
	}, nil
}
