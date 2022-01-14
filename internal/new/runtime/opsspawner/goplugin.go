package ops

import "github.com/emil14/neva/internal/new/runtime"

type Plugin struct{}

func (o Plugin) Operator(runtime.OpRef) (func(runtime.NodeIO) error, error) {
	return nil, nil
}
