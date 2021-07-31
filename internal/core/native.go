package core

type NativeModule struct {
	in   InportsInterface
	out  OutportsInterface
	impl func(NodeIO) error
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
	impl func(NodeIO) error,
) NativeModule {
	return NativeModule{
		in:   in,
		out:  out,
		impl: impl,
	}
}
