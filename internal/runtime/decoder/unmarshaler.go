package decoder

import (
	"errors"
	"fmt"

	"github.com/emil14/neva/pkg/runtimesdk"
	"github.com/gogo/protobuf/proto"
)

var ErrProto = errors.New("proto")

type unmarshaler struct{}

func (u unmarshaler) Unmarshal(bb []byte, prog *runtimesdk.Program) error {
	if err := proto.Unmarshal(bb, prog); err != nil {
		return fmt.Errorf("%w: %v", ErrProto, err)
	}

	return nil
}

func NewUnmarshaler() Unmarshaler {
	return unmarshaler{}
}
