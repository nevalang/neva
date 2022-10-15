package runtimehelpers

import "github.com/emil14/neva/pkg/runtimesdk"

func Prog(
	start *runtimesdk.PortAddr,
	ports []*runtimesdk.Port,
	effects *runtimesdk.Effects,
	conns []*runtimesdk.Connection,
) *runtimesdk.Program {
	return &runtimesdk.Program{
		StartPort:   start,
		Ports:       ports,
		Effects:     effects,
		Connections: conns,
	}
}

func Connections(cc ...*runtimesdk.Connection) []*runtimesdk.Connection {
	return cc
}

func Connection(sender *runtimesdk.PortAddr, receivers []*runtimesdk.ConnectionPoint) *runtimesdk.Connection {
	return &runtimesdk.Connection{
		SenderOutPortAddr:        sender,
		ReceiverConnectionPoints: receivers,
	}
}

func Points(pp ...*runtimesdk.ConnectionPoint) []*runtimesdk.ConnectionPoint {
	return pp
}

func Point(in *runtimesdk.PortAddr) *runtimesdk.ConnectionPoint {
	return &runtimesdk.ConnectionPoint{
		InPortAddr:      in,
		Type:            0,
		StructFieldPath: []string{},
	}
}

func Constants(cc ...*runtimesdk.Constant) []*runtimesdk.Constant {
	return cc
}

func Triggers(tt ...*runtimesdk.Trigger) []*runtimesdk.Trigger {
	return tt
}

func Constant(out *runtimesdk.PortAddr, msg *runtimesdk.Msg) *runtimesdk.Constant {
	return &runtimesdk.Constant{
		OutPortAddr: out,
		Msg:         msg,
	}
}

func Trigger(in, out *runtimesdk.PortAddr, msg *runtimesdk.Msg) *runtimesdk.Trigger {
	return &runtimesdk.Trigger{
		InPortAddr:  in,
		OutPortAddr: out,
		Msg:         msg,
	}
}

func StrMsg(s string) *runtimesdk.Msg {
	return &runtimesdk.Msg{
		Str:  "hello world!\n",
		Type: runtimesdk.MsgType_VALUE_TYPE_STR, //nolint
	}
}

func Operators(oo ...*runtimesdk.Operator) []*runtimesdk.Operator {
	return oo
}

func Operator(ref *runtimesdk.OperatorRef, in, out []*runtimesdk.PortAddr) *runtimesdk.Operator {
	return &runtimesdk.Operator{
		Ref:          ref,
		InPortAddrs:  in,
		OutPortAddrs: out,
	}
}

func OperatorRef(pkg, name string) *runtimesdk.OperatorRef {
	return &runtimesdk.OperatorRef{
		Pkg: pkg, Name: name,
	}
}

func Ports(pp ...*runtimesdk.Port) []*runtimesdk.Port {
	return pp
}

func PortAddrs(pp ...*runtimesdk.PortAddr) []*runtimesdk.PortAddr {
	return pp
}

func Port(path, name string, buf uint32) *runtimesdk.Port {
	return &runtimesdk.Port{
		Addr:    PortAddr(path, name),
		BufSize: buf,
	}
}

func PortAddr(path, name string) *runtimesdk.PortAddr {
	return Slot(path, name, 0)
}

func Slot(path, name string, idx uint32) *runtimesdk.PortAddr {
	return &runtimesdk.PortAddr{
		Path: path, Port: name, Idx: idx,
	}
}

func Effects(
	Operators []*runtimesdk.Operator,
	Constants []*runtimesdk.Constant,
	Triggers []*runtimesdk.Trigger,
) *runtimesdk.Effects {
	return &runtimesdk.Effects{
		Operators: Operators,
		Constants: Constants,
		Triggers:  Triggers,
	}
}
