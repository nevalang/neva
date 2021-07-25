package core

type NativeModule struct {
	in   Inport
	out  Outports
	impl func(NodeIO)
}

func (a NativeModule) Interface() Interface {
	return Interface{
		In:  a.in,
		Out: a.out,
	}
}

func NewNativeModule(
	in Inport,
	out Outports,
	impl func(NodeIO),
) NativeModule {
	return NativeModule{
		in:   in,
		out:  out,
		impl: impl,
	}
}
