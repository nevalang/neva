package decoder

type decoder struct{}

func MustNew() decoder {
	return decoder{}
}
