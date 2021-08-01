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

// TODO use?
func (nmod nativeModule) startStream(io NodeIO) error {
	// check io
	// return err if needed
	// run go impl
	// return nil
	return nil
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
