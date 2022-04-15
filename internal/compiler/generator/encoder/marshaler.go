package encoder

import (
	"errors"
	"fmt"

	"github.com/emil14/neva/pkg/runtimesdk"
	"github.com/gogo/protobuf/proto"
)

type marshaler struct{}

var ErrProto = errors.New("proto marshal")

func (m marshaler) Marshal(sdkProg *runtimesdk.Program) ([]byte, error) {
	bb, err := proto.Marshal(sdkProg)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrProto, err)
	}

	return bb, nil
}

func NewMarshaler() marshaler {
	return marshaler{}
}
