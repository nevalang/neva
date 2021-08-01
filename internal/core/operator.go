package core

type operator struct {
	in   InportsInterface
	out  OutportsInterface
	impl func(NodeIO) error // should return error if io invalid
}

func (a operator) Interface() Interface {
	return Interface{
		In:  a.in,
		Out: a.out,
	}
}

func (op operator) startStream(io NodeIO) error {
	// check io
	// return err if needed
	// run go impl
	// return nil
	return nil
}

func NewOperator(
	in InportsInterface,
	out OutportsInterface,
	impl func(NodeIO) error,
) operator {
	return operator{
		in:   in,
		out:  out,
		impl: impl,
	}
}
