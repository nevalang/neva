package main

import "github.com/emil14/neva/pkg/runtimesdk"

func prog(
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

func conns(cc ...*runtimesdk.Connection) []*runtimesdk.Connection {
	return cc
}

func conn(sender *runtimesdk.PortAddr, receivers []*runtimesdk.ConnectionPoint) *runtimesdk.Connection {
	return &runtimesdk.Connection{
		SenderOutPortAddr:        sender,
		ReceiverConnectionPoints: receivers,
	}
}

func points(pp ...*runtimesdk.ConnectionPoint) []*runtimesdk.ConnectionPoint {
	return pp
}

func point(in *runtimesdk.PortAddr) *runtimesdk.ConnectionPoint {
	return &runtimesdk.ConnectionPoint{
		InPortAddr:      in,
		Type:            0,
		StructFieldPath: []string{},
	}
}

func consts(cc ...*runtimesdk.Constant) []*runtimesdk.Constant {
	return cc
}

func triggers(tt ...*runtimesdk.Trigger) []*runtimesdk.Trigger {
	return tt
}

func cnst(out *runtimesdk.PortAddr, msg *runtimesdk.Msg) *runtimesdk.Constant {
	return &runtimesdk.Constant{
		OutPortAddr: out,
		Msg:         msg,
	}
}

func trigger(in, out *runtimesdk.PortAddr, msg *runtimesdk.Msg) *runtimesdk.Trigger {
	return &runtimesdk.Trigger{
		InPortAddr:  in,
		OutPortAddr: out,
		Msg:         msg,
	}
}

func strMsg(s string) *runtimesdk.Msg {
	return &runtimesdk.Msg{
		Str:  "hello world!\n",
		Type: runtimesdk.MsgType_VALUE_TYPE_STR, //nolint
	}
}

func ops(oo ...*runtimesdk.Operator) []*runtimesdk.Operator {
	return oo
}

func op(ref *runtimesdk.OperatorRef, in, out []*runtimesdk.PortAddr) *runtimesdk.Operator {
	return &runtimesdk.Operator{
		Ref:          ref,
		InPortAddrs:  in,
		OutPortAddrs: out,
	}
}

func opref(pkg, name string) *runtimesdk.OperatorRef {
	return &runtimesdk.OperatorRef{
		Pkg: pkg, Name: name,
	}
}

func ports(pp ...*runtimesdk.Port) []*runtimesdk.Port {
	return pp
}

func portsAddrs(pp ...*runtimesdk.PortAddr) []*runtimesdk.PortAddr {
	return pp
}

func port(path, name string, buf uint32) *runtimesdk.Port {
	return &runtimesdk.Port{
		Addr:    portAddr(path, name),
		BufSize: buf,
	}
}

func portAddr(path, name string) *runtimesdk.PortAddr {
	return slot(path, name, 0)
}

func slot(path, name string, idx uint32) *runtimesdk.PortAddr {
	return &runtimesdk.PortAddr{
		Path: path, Port: name, Idx: idx,
	}
}

func effects(
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
