package runtime

import "fbp/internal/types"

// AbsModule represents abstract module.
type AbsModule interface {
	Run(in map[string]<-chan Msg, out map[string]chan<- Msg)
	Ports() (InPorts, OutPorts)
}

// Module represents used refined module.
type Module struct {
	in  InPorts
	out OutPorts
	wm  WorkerMap
	net []Conn
}

func (m Module) Ports() (InPorts, OutPorts) {
	return m.in, m.out
}

type InPorts PortMap

type OutPorts PortMap

type WorkerMap map[string]AbsModule

type Conn struct {
	sender    <-chan Msg   // outPort
	receivers []chan<- Msg // inPorts
}

type PortMap map[string]types.Type

func NewModule(
	in InPorts,
	out OutPorts,
	wm WorkerMap,
	net []Conn,
) Module {
	return Module{
		in:  in,
		out: out,
		wm:  wm,
		net: net,
	}
}

// NativeModule represents native module implementation.
type NativeModule struct {
	in   InPorts
	out  OutPorts
	impl func(
		in map[string]<-chan Msg,
		out map[string]chan<- Msg,
	)
}

func (nm NativeModule) Run(in map[string]<-chan Msg, out map[string]chan<- Msg) {
	nm.impl(in, out)
}

func (nm NativeModule) Ports() (InPorts, OutPorts) {
	return nm.in, nm.out
}

func NewNativeModule(
	in InPorts,
	out OutPorts,
	impl func(
		in map[string]<-chan Msg,
		out map[string]chan<- Msg,
	),
) NativeModule {
	return NativeModule{
		in:   in,
		out:  out,
		impl: impl,
	}
}
