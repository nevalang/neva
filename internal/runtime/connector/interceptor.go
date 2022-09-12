package connector

import (
	"github.com/emil14/neva/internal/core"
	"github.com/emil14/neva/internal/runtime"
)

type DefaultInterceptor struct{}

func (c DefaultInterceptor) AfterSending(_ runtime.Connection, msg core.Msg) core.Msg {
	return msg
}

func (c DefaultInterceptor) BeforeReceiving(from, to runtime.PortAddr, msg core.Msg) core.Msg {
	return msg
}
