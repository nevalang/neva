package core

type operator struct {
	in   InportsInterface
	out  OutportsInterface
	impl func(NodeIO) error
}

func (a operator) Interface() Interface {
	return Interface{
		In:  a.in,
		Out: a.out,
	}
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
