package protocoder

import (
	"github.com/golang/protobuf/proto"

	"github.com/emil14/neva/internal/runtime/program"
	"github.com/emil14/neva/pkg/runtimesdk"
)

type coder struct{}

func (c coder) Code(prog program.Program) ([]byte, error) {
	p := c.cast(prog)
	return proto.Marshal(&p)
}

func (c coder) cast(prog program.Program) runtimesdk.Program {
	return runtimesdk.Program{}
}
