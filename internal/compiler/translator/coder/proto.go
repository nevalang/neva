package coder

import (
	"github.com/emil14/neva/internal/runtime/program"
	"github.com/emil14/neva/pkg/runtimesdk"
	"github.com/golang/protobuf/proto"
)

type protobuf struct{}

func (pb protobuf) Code(prog program.Program) ([]byte, error) {
	p := pb.cast(prog)
	return proto.Marshal(&p)
}

func (pb protobuf) cast(prog program.Program) runtimesdk.Program {
	return runtimesdk.Program{}
}
