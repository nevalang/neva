package decoder

import "github.com/emil14/neva/internal/new/runtime"

type Proto struct{}

func (p Proto) Decode([]byte) (runtime.Program, error) {
	return runtime.Program{}, nil
}
