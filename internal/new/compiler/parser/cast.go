package parser

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/emil14/neva/internal/new/compiler"
)

type caster struct{}

func (c caster) From(mod compiler.Module) module {
	in, out := c.fromIO(mod.IO)
	deps := c.fromDeps(mod.Deps)
	cnst := c.fromConst(mod.Const)
	net := c.fromNet(mod.Net)
	return module{
		Deps:    deps,
		In:      in,
		Out:     out,
		Const:   cnst,
		Workers: mod.Workers,
		Net:     net,
	}
}

func (c caster) To(mod module) compiler.Module {
	return compiler.Module{
		IO:      c.toIO(mod.In, mod.Out),
		Deps:    c.toDeps(mod.Deps),
		Const:   c.toConst(mod.Const),
		Workers: map[string]string(mod.Workers),
		Net:     c.toNet(mod.Net),
	}
}

func (c caster) fromConst(map[string]compiler.Const) map[string]Const {
	return map[string]Const{} // TODO
}

func (c caster) fromIO(io compiler.IO) (inports, outports) {
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

func (c caster) fromDeps(deps map[string]compiler.IO) moduleDeps {
	result := moduleDeps{}
	for k, v := range deps {
		in, out := c.fromIO(v)
		result[k] = io{
			In:  in,
			Out: out,
		}
	}
	return result
}

func (c caster) fromNet(net compiler.Net) net {
	return nil // TODO
}

func (c caster) toIO(in inports, out outports) compiler.IO {
	return compiler.IO{
		In:  c.castPorts(ports(in)),
		Out: c.castPorts(ports(out)),
	}
}

func (c caster) castPorts(from ports) compiler.Ports {
	to := compiler.Ports{}

	for port, t := range from {
		portType := compiler.PortType{Type: compiler.TypeByName(t)}
		if strings.HasSuffix(port, "[]") {
			portType.Arr = true
			port = strings.TrimSuffix(port, "[]")
		}

		to[port] = portType
	}

	return to
}

func (c caster) toDeps(from moduleDeps) map[string]compiler.IO {
	to := map[string]compiler.IO{}

	for name, pio := range from {
		io := c.toIO(pio.In, pio.Out)

		to[name] = compiler.IO{
			In:  io.In,
			Out: io.Out,
		}
	}

	return to
}

func (c caster) toConst(from map[string]Const) map[string]compiler.Const {
	res := map[string]compiler.Const{}

	for name, cnst := range from {
		switch cnst.Type {
		case compiler.IntType.String():
			res[name] = compiler.Const{
				Type:     compiler.IntType,
				IntValue: cnst.IntValue,
			}
		}
	}

	return res
}

func (c caster) toNet(from net) compiler.Net {
	to := compiler.Net{}

	for senderNode, outgoingConnections := range from {
		for outport, nodesToInports := range outgoingConnections {
			senderPortPoint := c.castPortPoint(senderNode, outport)
			receivers := map[compiler.PortAddr]struct{}{}

			for receiver, receiverInports := range nodesToInports {
				for _, inport := range receiverInports {
					receiverPortPoint := c.castPortPoint(receiver, inport)
					receivers[receiverPortPoint] = struct{}{}
				}
			}

			to[senderPortPoint] = receivers
		}
	}

	return to
}

func (c caster) castPortPoint(node string, port string) compiler.PortAddr {
	bracketStart := strings.Index(port, "[")
	if bracketStart == -1 {
		return compiler.PortAddr{
			Node: node,
			Port: port,
		}
	}

	bracketEnd := strings.Index(port, "]")
	if bracketEnd == -1 {
		return compiler.PortAddr{
			Node: node,
			Port: port,
		}
	}

	idx, err := strconv.ParseUint(port[bracketStart+1:bracketEnd], 10, 64)
	if err != nil {
		return compiler.PortAddr{
			Node: node,
			Port: port,
		}
	}

	return compiler.PortAddr{
		Node: node,
		Port: port[:bracketStart],
		Slot: uint8(idx),
	}
}
