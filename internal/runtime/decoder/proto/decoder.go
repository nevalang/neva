package decoder

import (
	"errors"
	"fmt"

	"github.com/emil14/neva/internal/pkg/tools"
	"github.com/emil14/neva/internal/runtime/src"
	"github.com/emil14/neva/pkg/runtimesdk"
)

type Proto struct {
	unmarshaler Unmarshaler
	caster      caster
}

type (
	Unmarshaler interface {
		Unmarshal([]byte, *runtimesdk.Program) error
	}
)

var (
	ErrCast      = errors.New("cast")
	ErrUnmarshal = errors.New("unmarshal")
)

func (p Proto) Decode(bb []byte) (src.Program, error) {
	var sdkProg runtimesdk.Program
	if err := p.unmarshaler.Unmarshal(bb, &sdkProg); err != nil {
		return src.Program{}, fmt.Errorf("%w: %v", ErrUnmarshal, err)
	}

	prog, err := p.caster.Cast(&sdkProg)
	if err != nil {
		return src.Program{}, fmt.Errorf("%w: %v", ErrCast, err)
	}

	return prog, nil
}

func MustNewProto(unmarshaler Unmarshaler) Proto {
	tools.PanicWithNil(unmarshaler)

	return Proto{
		caster:      caster{},
		unmarshaler: unmarshaler,
	}
}
