package decoder

import (
	"errors"
	"fmt"

	"github.com/emil14/neva/internal/pkg/utils"
	"github.com/emil14/neva/internal/runtime"
	"github.com/emil14/neva/pkg/runtimesdk"
)

var (
	ErrCast      = errors.New("cast")
	ErrUnmarshal = errors.New("unmarshal")
)

type (
	Unmarshaler interface {
		Unmarshal([]byte, *runtimesdk.Program) error
	}
	Caster interface {
		Cast(*runtimesdk.Program) (runtime.Program, error)
	}
)

type Proto struct {
	unmarshaler Unmarshaler
	caster      Caster
}

func (p Proto) Decode(bb []byte) (runtime.Program, error) {
	var sdkProg runtimesdk.Program
	if err := p.unmarshaler.Unmarshal(bb, &sdkProg); err != nil {
		return runtime.Program{}, fmt.Errorf("%w: %v", ErrUnmarshal, err)
	}

	prog, err := p.caster.Cast(&sdkProg)
	if err != nil {
		return runtime.Program{}, fmt.Errorf("%w: %v", ErrCast, err)
	}

	return prog, nil
}

func MustNewProto(caster Caster, unmarshaler Unmarshaler) Proto {
	utils.PanicOnNil(caster, unmarshaler)

	return Proto{
		caster:      caster,
		unmarshaler: unmarshaler,
	}
}
