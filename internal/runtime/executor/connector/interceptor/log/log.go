package log

import (
	"fmt"
	"strings"

	"github.com/emil14/neva/internal/core"
	"github.com/emil14/neva/internal/pkg/tools"
	"github.com/emil14/neva/internal/runtime/src"
)

type Logger interface {
	Printf(format string, v ...any)
}

type Interceptor struct {
	logger Logger
}

func (i Interceptor) AfterSending(conn src.Connection, msg core.Msg) core.Msg {
	i.logger.Printf("sent: %s", i.formatConnection(conn, msg))
	return msg
}

func (i Interceptor) BeforeReceiving(
	saddr src.Ports,
	rpoint src.ConnectionSide,
	msg core.Msg,
) core.Msg {
	i.logger.Printf(
		"prepare: %s <- %s <- %s",
		i.formatPortAddr(rpoint.PortAddr), msg, i.formatPortAddr(saddr),
	)
	return msg
}

func (i Interceptor) AfterReceiving(
	saddr src.Ports,
	rpoint src.ConnectionSide,
	msg core.Msg,
) {
	i.logger.Printf(
		"received: %s <- %s <- %s",
		i.formatPortAddr(rpoint.PortAddr), msg, i.formatPortAddr(saddr),
	)
}

func (i Interceptor) formatConnection(conn src.Connection, msg core.Msg) string {
	to := []string{}
	for _, r := range conn.ReceiverSides {
		s := i.formatPortAddr(r.PortAddr)
		if r.Action == src.ReadDict {
			s = "." + strings.Join(r.Payload.ReadDict, ".")
		}
		to = append(to, s)
	}

	return fmt.Sprintf(
		"%s -> %s -> %s",
		i.formatPortAddr(conn.SenderSide.PortAddr),
		msg,
		strings.Join(to, ", "),
	)
}

func (i Interceptor) formatPortAddr(addr src.Ports) string {
	return fmt.Sprintf("%s.%s[%d]", addr.Path, addr.Port, addr.Idx)
}

func MustNew(l Logger) Interceptor {
	tools.PanicWithNil(l)
	return Interceptor{l}
}
