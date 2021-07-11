package runtime

import "github.com/emil14/refactored-garbanzo/internal/types"

// AbstractModule represents atomic and complex modules.
type AbstractModule interface {
	Run(in map[string]<-chan Msg, out map[string]chan<- Msg)
	Ports() (InPorts, OutPorts)
}

// AtomicModule is a module with the native implementation.
type AtomicModule struct {
	in   InPorts
	out  OutPorts
	impl func(
		in map[string]<-chan Msg,
		out map[string]chan<- Msg,
	)
}

func (a AtomicModule) Run(in map[string]<-chan Msg, out map[string]chan<- Msg) {
	a.impl(in, out)
}

func (a AtomicModule) Ports() (InPorts, OutPorts) {
	return a.in, a.out
}

func NewAtomicModule(
	in InPorts,
	out OutPorts,
	impl func(
		in map[string]<-chan Msg,
		out map[string]chan<- Msg,
	),
) AtomicModule {
	return AtomicModule{
		in:   in,
		out:  out,
		impl: impl,
	}
}

// ComplexModule is a composition of other modules.
type ComplexModule struct {
	in  InPorts
	out OutPorts
	net []Conn
}

func (cm ComplexModule) Ports() (InPorts, OutPorts) {
	return cm.in, cm.out
}

func (m ComplexModule) Run(in map[string]chan Msg, out map[string]chan Msg) {
	ConnectAll(m.net)
}

type InPorts Ports

type OutPorts Ports

type Env map[string]AbstractModule

type Conn struct {
	Sender    <-chan Msg
	Receivers []chan Msg
}

type Ports map[string]types.Type

func NewComplexModule(
	in InPorts,
	out OutPorts,
	net []Conn,
) ComplexModule {
	return ComplexModule{
		in:  in,
		out: out,
		net: net,
	}
}
