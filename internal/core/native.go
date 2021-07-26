package core

type NativeModule struct {
	in   InportsInterface
	out  OutportsInterface
	impl func(NodeIO)
}

func (a NativeModule) Interface() Interface {
	return Interface{
		In:  a.in,
		Out: a.out,
	}
}

func NewNativeModule(
	in InportsInterface,
	out OutportsInterface,
	impl func(NodeIO),
) NativeModule {
	return NativeModule{
		in:   in,
		out:  out,
		impl: impl,
	}
}
