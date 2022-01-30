package parser

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/emil14/neva/internal/new/compiler"
)

var (
	ErrNet               = errors.New("net")
	ErrPortAddr          = errors.New("port addr")
	ErrPortAddrParts     = errors.New("string must consist of node and port separated by dot")
	ErrNodeNotFound      = errors.New("node not found")
	ErrPortNotFound      = errors.New("port not found")
	ErrComponentNotFound = errors.New("component not found")
	ErrPortType          = errors.New("port type")
)

type caster struct{}

func (c caster) Cast(mod Module) (compiler.Module, error) {
	io := c.castIO(mod.IO)
	deps := c.castDeps(mod.Deps)
	nodes := c.castNodes(mod.Nodes)

	net, err := c.castNet(io, deps, nodes, mod.Net)
	if err != nil {
		return compiler.Module{}, fmt.Errorf("%w: %v", ErrNet, err)
	}

	return compiler.Module{
		IO:     io,
		DepsIO: deps,
		Nodes:  nodes,
		Net:    net,
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
		return compiler.ArrPort
	}
	return compiler.NormPort
}

func (c caster) castMsgType(typ string) compiler.MsgType {
	switch typ {
	case "int":
		return compiler.IntMsg
	case "str":
		return compiler.StrMsg
	case "bool":
		return compiler.BoolMsg
	case "sig":
		return compiler.SigMsg
	}
	return compiler.UnknownMsg
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
			Type:     compiler.StrMsg,
			StrValue: msg.Str,
		}
	case "int":
		return compiler.Msg{
			Type:     compiler.IntMsg,
			IntValue: msg.Int,
		}
	case "bool":
		return compiler.Msg{
			Type:      compiler.IntMsg,
			BoolValue: msg.Bool,
		}
	case "sig":
		return compiler.Msg{
			Type: compiler.IntMsg,
		}
	}
	return compiler.Msg{
		Type: compiler.UnknownMsg,
	}
}

func (c caster) castNet(
	io compiler.IO,
	deps map[string]compiler.IO,
	nodes compiler.ModuleNodes,
	net map[string][]string,
) ([]compiler.Connection, error) {
	cnet := make([]compiler.Connection, 0, len(net))

	for from, to := range net {
		cfrom, err := c.castPortAddr(true, io, deps, nodes, from)
		if err != nil {
			return nil, fmt.Errorf("%w: from: %v", ErrPortAddr, err)
		}

		cto, err := c.castPortAddrs(false, io, deps, nodes, to)
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

func (c caster) castPortAddrs(
	isOutGoing bool,
	io compiler.IO,
	deps map[string]compiler.IO,
	nodes compiler.ModuleNodes,
	addrs []string,
) ([]compiler.AbsPortAddr, error) {
	caddrs := make([]compiler.AbsPortAddr, 0, len(addrs))

	for _, addr := range addrs {
		caddr, err := c.castPortAddr(isOutGoing, io, deps, nodes, addr)
		if err != nil {
			return nil, fmt.Errorf("%w: %v", ErrPortAddr, err)
		}

		caddrs = append(caddrs, caddr)
	}

	return caddrs, nil
}

func (c caster) castPortAddr(
	isOutGoing bool,
	io compiler.IO,
	deps map[string]compiler.IO,
	nodes compiler.ModuleNodes,
	addr string,
) (compiler.AbsPortAddr, error) {
	parts := strings.Split(addr, ".")
	if len(parts) != 2 {
		return compiler.AbsPortAddr{}, fmt.Errorf("%w: %s", ErrPortAddrParts, addr)
	}

	node := parts[0]
	port := parts[1]

	portType, err := c.portType(isOutGoing, io, deps, nodes, node, port)
	if err != nil {
		return compiler.AbsPortAddr{}, fmt.Errorf("%w: %v", ErrPortType, addr)
	}

	typ := compiler.NormPortAddr
	if portType == compiler.ArrPort {
		typ = compiler.ArrByPassPortAddr
	}

	portName, portIdx, ok := c.splitPort(port)
	if !ok {
		return compiler.AbsPortAddr{
			Type: typ,
			Node: node,
			Port: port,
		}, nil
	}

	return compiler.AbsPortAddr{
		Type: compiler.NormPortAddr,
		Node: node,
		Port: portName,
		Idx:  portIdx,
	}, nil
}

func (c caster) portType(
	isOutGoing bool,
	io compiler.IO,
	deps map[string]compiler.IO,
	nodes compiler.ModuleNodes,
	nodeName, portName string,
) (compiler.PortType, error) {
	if nodeName == "in" {
		port, ok := io.In[portName]
		if !ok {
			return 0, fmt.Errorf("%w: %s.%s", ErrPortNotFound, nodeName, portName)
		}
		return port.Type, nil
	}

	if nodeName == "out" {
		port, ok := io.Out[portName]
		if !ok {
			return 0, fmt.Errorf("%w: %s.%s", ErrPortNotFound, nodeName, portName)
		}
		return port.Type, nil
	}

	if nodeName == "const" {
		return compiler.NormPort, nil
	}

	worker, ok := nodes.Workers[nodeName]
	if !ok {
		return 0, fmt.Errorf("%w: %s.%s", ErrNodeNotFound, nodeName, portName)
	}

	dep, ok := deps[worker]
	if !ok {
		return 0, fmt.Errorf("%w: %s", ErrComponentNotFound, worker)
	}

	if isOutGoing {
		port, ok := dep.Out[portName]
		if !ok {
			return 0, fmt.Errorf("%w: %s.%s", ErrPortNotFound, nodeName, portName)
		}
		return port.Type, nil
	}

	port, ok := dep.Out[portName]
	if !ok {
		return 0, fmt.Errorf("%w: %s.%s", ErrPortNotFound, nodeName, portName)
	}

	return port.Type, nil
}

func (c caster) splitPort(port string) (string, uint8, bool) {
	bracketStart := strings.Index(port, "[")
	if bracketStart == -1 {
		return "", 0, false
	}

	bracketEnd := strings.Index(port, "]")
	if bracketEnd == -1 {
		return "", 0, false
	}

	idx, err := strconv.ParseUint(port[bracketStart+1:bracketEnd], 10, 64)
	if err != nil {
		return "", 0, false
	}

	return port[:bracketStart], uint8(idx), true
}
