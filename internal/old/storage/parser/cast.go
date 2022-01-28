package parser

import (
	"fmt"
	"strconv"
	"strings"

	program "github.com/emil14/neva/internal/new/compiler"
)

type yamlCaster struct{}

func (c yamlCaster) from(mod program.Module) module {
	in, out := c.fromIO(mod.IO)
	return module{
		IO:      io{In: in, Out: out},
		Deps:    c.fromDeps(mod.Deps),
		Const:   c.fromConst(mod.Nodes.Const),
		Workers: mod.Nodes.Workers,
		Start:   false,
		Net:     c.fromNet(mod.Net),
	}
}

func (c yamlCaster) fromConst(from map[string]program.Const) map[string]constant {
	to := map[string]constant{}

	for k, v := range from {
		to[k] = constant{
			Type:     string(v.Type()),
			IntValue: v.Int(),
		}
	}

	return to
}

func (c yamlCaster) fromIO(io program.IO) (inports, outports) {
	in := inports{}
	for k, v := range io.In {
		in[k] = fmt.Sprintf("%v", v)
	}

	out := outports{}
	for k, v := range io.Out {
		out[k] = fmt.Sprintf("%v", v)
	}

	return in, out
}

func (c yamlCaster) fromDeps(from map[string]program.IO) deps {
	to := deps{}

	for dep, depIO := range from {
		in, out := c.fromIO(depIO)
		to[dep] = io{
			Params: []string{}, // todo
			In:     in,
			Out:    out,
		}
	}

	return to
}

func (c yamlCaster) fromNet(net program.ModuleNet) net {
	return nil // TODO
}

func (c yamlCaster) to(mod module) (program.Module, error) {
	connections, err := c.toNet(mod.Net)
	if err != nil {
		return program.Module{}, err
	}

	return program.Module{
		IO:   c.toIO(mod.IO.In, mod.IO.Out),
		Deps: c.toDeps(mod.Deps),
		Nodes: program.ModuleNodes{
			Const:   map[string]program.Const{},
			Workers: mod.Workers,
		},
		Net: connections,
	}, nil
}

func (c yamlCaster) toIO(in map[string]string, out map[string]string) program.IO {
	return program.IO{
		In:  c.toPorts(ports(in)),
		Out: c.toPorts(ports(out)),
	}
}

func (c yamlCaster) toPorts(from map[string]string) program.Ports {
	to := program.Ports{}

	for port, t := range from {
		typ, err := toTypeName(t)
		if err != nil {
			return program.Ports{}
		}
		portType := program.PortType{DataType: typ}
		if strings.HasSuffix(port, "[]") {
			portType.IsArr = true
			port = strings.TrimSuffix(port, "[]")
		}

		to[port] = portType
	}

	return to
}

func (c yamlCaster) toDeps(from deps) map[string]program.IO {
	to := map[string]program.IO{}

	for name, pio := range from {
		io := c.toIO(pio.In, pio.Out)
		to[name] = program.IO{
			In:  io.In,
			Out: io.Out,
		}
	}

	return to
}

func (c yamlCaster) toConst(from map[string]constant) map[string]program.Const {
	to := map[string]program.Const{}

	for name, cnst := range from {
		switch cnst.Type {
		case program.DataTypeInt.String():
			to[name] = program.NewIntConst(cnst.IntValue)
		case program.TypeStr.String():
			to[name] = program.NewStrConst(cnst.StrValue)
		case program.TypeBool.String():
			to[name] = program.NewBoolConst(cnst.BoolValue)
		}
	}

	return to
}

func (c yamlCaster) toNet(from map[string]string) (program.ModuleNet, error) {
	to := make(program.ModuleNet, len(from))

	for sender, receivers := range from {
		rsvrs := make(map[program.ConnectionPoint]struct{}, len(receivers))
		for _, receiver := range receivers {
			receiverAddr, err := c.toPortAddr(receiver)
			if err != nil {
				return nil, err
			}
			rsvrs[receiverAddr] = struct{}{}
		}

		senderAddr, err := c.toPortAddr(sender)
		if err != nil {
			return nil, err
		}

		to[senderAddr] = rsvrs
	}

	return to, nil
}

func (c yamlCaster) toPortAddr(s string) (program.ConnectionPoint, error) {
	parts := strings.Split(s, ".")
	if len(parts) != 2 {
		return program.ConnectionPoint{}, fmt.Errorf("invalid sender %s", s)
	}

	var (
		node = parts[0]
		port = parts[1]
	)

	bracketStart := strings.Index(port, "[")
	if bracketStart == -1 {
		return program.ConnectionPoint{Node: node, Port: port}, nil
	}

	bracketEnd := strings.Index(port, "]")
	if bracketEnd == -1 {
		return program.ConnectionPoint{Node: node, Port: port}, nil
	}

	idx, err := strconv.ParseUint(port[bracketStart+1:bracketEnd], 10, 64)
	if err != nil {
		return program.ConnectionPoint{Node: node, Port: port}, nil
	}

	return program.ConnectionPoint{
		Node: node,
		Port: port[:bracketStart],
		Idx:  uint8(idx),
	}, nil
}

type typeName string

const (
	intType  typeName = "int"
	strType  typeName = "str"
	boolType typeName = "bool"
	sigType  typeName = "sig"
)

func (p parser) toTypeName(name string) (program.MsgType, error) {
	switch typeName(name) {
	case intType:
		return program.TypeInt, nil
	case strType:
		return program.TypeStr, nil
	case boolType:
		return program.TypeBool, nil
	}

	return 0, fmt.Errorf("unknown type %s", name)
}
