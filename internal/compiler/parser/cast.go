package parser

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/emil14/respect/internal/compiler/program"
)

type caster struct{}

func (c caster) From(mod program.Module) module {
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

func (c caster) To(mod module) program.Module {
	return program.Module{
		IO:      c.toIO(mod.In, mod.Out),
		Deps:    c.toDeps(mod.Deps),
		Const:   c.toConst(mod.Const),
		Workers: map[string]string(mod.Workers),
		Net:     c.toNet(mod.Net),
	}
}

func (c caster) fromConst(map[string]program.Const) map[string]Const {
	return map[string]Const{} // TODO
}

func (c caster) fromIO(io program.IO) (inports, outports) {
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

func (c caster) fromDeps(deps map[string]program.IO) moduleDeps {
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

func (c caster) fromNet(net program.Net) net {
	return nil // TODO
}

func (c caster) toIO(in inports, out outports) program.IO {
	return program.IO{
		In:  c.castPorts(ports(in)),
		Out: c.castPorts(ports(out)),
	}
}

func (c caster) castPorts(from ports) program.Ports {
	to := program.Ports{}

	for port, t := range from {
		portType := program.PortType{Type: program.TypeByName(t)}
		if strings.HasSuffix(port, "[]") {
			portType.Arr = true
			port = strings.TrimSuffix(port, "[]")
		}

		to[port] = portType
	}

	return to
}

func (c caster) toDeps(from moduleDeps) map[string]program.IO {
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

func (c caster) toConst(from map[string]Const) map[string]program.Const {
	res := map[string]program.Const{}

	for name, cnst := range from {
		switch cnst.Type {
		case program.IntType.String():
			res[name] = program.Const{
				Type:     program.IntType,
				IntValue: cnst.IntValue,
			}
		}
	}

	return res
}

func (c caster) toNet(from net) program.Net {
	to := program.Net{}

	for senderNode, outgoingConnections := range from {
		for outport, nodesToInports := range outgoingConnections {
			senderPortPoint := c.castPortPoint(senderNode, outport)
			receivers := map[program.PortAddr]struct{}{}

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

func (c caster) castPortPoint(node string, port string) program.PortAddr {
	bracketStart := strings.Index(port, "[")
	if bracketStart == -1 {
		return program.PortAddr{
			Node: node,
			Port: port,
		}
	}

	bracketEnd := strings.Index(port, "]")
	if bracketEnd == -1 {
		return program.PortAddr{
			Node: node,
			Port: port,
		}
	}

	idx, err := strconv.ParseUint(port[bracketStart+1:bracketEnd], 10, 64)
	if err != nil {
		return program.PortAddr{
			Node: node,
			Port: port,
		}
	}

	return program.PortAddr{
		Node: node,
		Port: port[:bracketStart],
		Slot:  uint8(idx),
	}
}
