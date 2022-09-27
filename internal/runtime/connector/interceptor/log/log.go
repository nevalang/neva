package log

import (
	"fmt"
	"strings"

	"github.com/emil14/neva/internal/core"
	"github.com/emil14/neva/internal/pkg/utils"
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
	saddr src.AbsolutePortAddr,
	rpoint src.ReceiverConnectionPoint,
	msg core.Msg,
) core.Msg {
	i.logger.Printf(
		"prepare: %s <- %s <- %s",
		i.formatPortAddr(rpoint.PortAddr), msg, i.formatPortAddr(saddr),
	)
	return msg
}

func (i Interceptor) AfterReceiving(
	saddr src.AbsolutePortAddr,
	rpoint src.ReceiverConnectionPoint,
	msg core.Msg,
) {
	i.logger.Printf(
		"received: %s <- %s <- %s",
		i.formatPortAddr(rpoint.PortAddr), msg, i.formatPortAddr(saddr),
	)
}

func (i Interceptor) formatConnection(conn src.Connection, msg core.Msg) string {
	to := []string{}
	for _, r := range conn.ReceiversConnectionPoints {
		s := i.formatPortAddr(r.PortAddr)
		if r.Type == src.DictReading {
			s = "." + strings.Join(r.DictReadingPath, ".")
		}
		to = append(to, s)
	}

	return fmt.Sprintf(
		"%s -> %s -> %s",
		i.formatPortAddr(conn.SenderPortAddr),
		msg,
		strings.Join(to, ", "),
	)
}

func (i Interceptor) formatPortAddr(addr src.AbsolutePortAddr) string {
	return fmt.Sprintf("%s.%s[%d]", addr.Path, addr.Port, addr.Idx)
}

func MustNew(l Logger) Interceptor {
	utils.NilPanic(l)
	return Interceptor{l}
}
