package encoder

import (
	"errors"
	"fmt"

	"github.com/emil14/neva/internal/new/runtime"
	"github.com/emil14/neva/internal/pkg/utils"
	"github.com/emil14/neva/pkg/runtimesdk"
)

type (
	Marshaler interface {
		Marshal(*runtimesdk.Program) ([]byte, error)
	}
	Caster interface {
		Cast(runtime.Program) (runtimesdk.Program, error)
	}
)

var (
	ErrCast    = errors.New("caster")
	ErrMarshal = errors.New("marshaller")
)

type Proto struct {
	marshaler Marshaler
	caster    Caster
}

func (p Proto) Encode(prog runtime.Program) ([]byte, error) {
	sdkProg, err := p.caster.Cast(prog)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrCast, err)
	}

	bb, err := p.marshaler.Marshal(&sdkProg)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrMarshal, err)
	}

	return bb, nil
}

func MustNew(marshaler Marshaler, caster Caster) Proto {
	utils.NilArgsFatal(marshaler)

	return Proto{marshaler, caster}
}
