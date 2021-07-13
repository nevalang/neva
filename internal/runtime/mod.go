package runtime

import "github.com/emil14/refactored-garbanzo/internal/types"

// AbstractModule represents atomic and complex modules.
type AbstractModule interface {
	Run(in map[string]<-chan Msg, out map[string]chan<- Msg)
	Ports() (InPorts, OutPorts)
}

type InPorts Ports

type OutPorts Ports

type Ports map[string]types.Type

// AtomicModule is a module with the native implementation.
type AtomicModule struct {
	in   InPorts
	out  OutPorts
	impl func(
		in map[string]<-chan Msg,
		out map[string]chan<- Msg,
	)
}

func (a AtomicModule) Ports() (InPorts, OutPorts) {
	return a.in, a.out
}

func (a AtomicModule) Run(in map[string]<-chan Msg, out map[string]chan<- Msg) {
	a.impl(in, out)
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
	net []ChanRel
}

func (cm ComplexModule) Ports() (InPorts, OutPorts) {
	return cm.in, cm.out
}

func (m ComplexModule) Run(in map[string]chan Msg, out map[string]chan Msg) {
	ConnectAll(m.net)
}

type Env map[string]AbstractModule

// ChanRel represents one-to-many relation between sender and receiver channels.
type ChanRel struct {
	Sender    <-chan Msg
	Receivers []chan Msg
}

func NewComplexModule(
	in InPorts,
	out OutPorts,
	net []ChanRel,
) ComplexModule {
	return ComplexModule{
		in:  in,
		out: out,
		net: net,
	}
}
