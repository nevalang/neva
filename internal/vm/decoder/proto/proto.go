package proto

import (
	"google.golang.org/protobuf/proto"

	"github.com/nevalang/neva/internal/runtime"
	"github.com/nevalang/neva/pkg/ir"
)

type Decoder struct {
	adapter Adapter
}

func (t Decoder) Decode(file []byte) (runtime.Program, error) {
	var irProg ir.Program
	if err := proto.Unmarshal(file, &irProg); err != nil {
		return runtime.Program{}, err
	}

	runtimeProg, err := t.adapter.Adapt(&irProg)
	if err != nil {
		return runtime.Program{}, err
	}

	return runtimeProg, nil
}

func NewDecoder(adapter Adapter) Decoder {
	return Decoder{
		adapter: adapter,
	}
}
