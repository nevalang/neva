package core

type NativeModule struct {
	in   InportsInterface
	out  OutportsInterface
	deps Deps
	impl func(NodeIO)
}

func (a NativeModule) Interface() Interface {
	return Interface{
		In:  a.in,
		Out: a.out,
	}
}

func (n NativeModule) Deps() Deps {
	return n.deps
}

func NewNativeModule(
	in InportsInterface,
	out OutportsInterface,
	deps Deps,
	impl func(NodeIO),
) NativeModule {
	return NativeModule{
		in:   in,
		out:  out,
		deps: deps,
		impl: impl,
	}
}


