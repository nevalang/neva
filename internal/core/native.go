package core

type nativeModule struct {
	in   InportsInterface
	out  OutportsInterface
	impl func(NodeIO) error // should return error if io invalid
}

func (a nativeModule) Interface() Interface {
	return Interface{
		In:  a.in,
		Out: a.out,
	}
}

func NewNativeModule(
	in InportsInterface,
	out OutportsInterface,
	impl func(NodeIO) error,
) nativeModule {
	return nativeModule{
		in:   in,
		out:  out,
		impl: impl,
	}
}
