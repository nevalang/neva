package connector

import (
	"context"
	"errors"

	"github.com/emil14/neva/internal/core"
	"github.com/emil14/neva/internal/pkg/utils"
	"github.com/emil14/neva/internal/runtime"
	"github.com/emil14/neva/internal/runtime/src"
)

type (
	Interceptor interface {
		AfterSending(src.Connection, core.Msg) core.Msg
		BeforeReceiving(saddr src.AbsolutePortAddr, point src.ReceiverConnectionPoint, msg core.Msg) core.Msg
		AfterReceiving(saddr src.AbsolutePortAddr, point src.ReceiverConnectionPoint, msg core.Msg)
	}
)

var (
	ErrMapper          = errors.New("mapper")
	ErrDictKeyNotFound = errors.New("dict key not found")
)

type Connector struct {
	interceptor Interceptor
}

func (c Connector) Connect(ctx context.Context, conns []runtime.Connection) error {
	return nil
}

func MustNew(i Interceptor) Connector {
	utils.NilPanic(i)
	return Connector{i}
}
