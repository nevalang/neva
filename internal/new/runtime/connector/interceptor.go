package connector

import (
	"github.com/emil14/neva/internal/new/core"
	"github.com/emil14/neva/internal/new/runtime"
)

type DefaultInterceptor struct{}

func (c DefaultInterceptor) AfterSend(_ runtime.Connection, msg core.Msg) core.Msg {
	return msg
}

func (c DefaultInterceptor) BeforeReceive(from, to runtime.FullPortAddr, msg core.Msg) core.Msg {
	return msg
}
