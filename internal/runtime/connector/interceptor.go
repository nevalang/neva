package connector

import (
	"fmt"
	"log"
	"strings"

	"github.com/emil14/neva/internal/core"
	"github.com/emil14/neva/internal/runtime"
)

type LoggingInterceptor struct{}

func (l LoggingInterceptor) AfterSending(connection runtime.Connection, msg core.Msg) core.Msg {
	log.Printf("after sending: %s, %s", l.formatConnection(connection), msg)
	return msg
}

func (l LoggingInterceptor) BeforeReceiving(
	sender, receiver runtime.AbsolutePortAddr,
	point runtime.ReceiverConnectionPoint,
	msg core.Msg,
) core.Msg {
	log.Printf("before receiving: %s <- %s, %s", l.formatPortAddr(receiver), l.formatPortAddr(sender), msg)
	return msg
}

func (l LoggingInterceptor) AfterReceiving(
	sender, receiver runtime.AbsolutePortAddr,
	point runtime.ReceiverConnectionPoint,
	msg core.Msg,
) {
	log.Printf("after receiving: %s <- %s, %s", l.formatPortAddr(receiver), l.formatPortAddr(sender), msg)
}

func (l LoggingInterceptor) formatConnection(connection runtime.Connection) string {
	to := []string{}
	for _, receiver := range connection.ReceiversConnectionPoints {
		s := l.formatPortAddr(receiver.PortAddr)
		if receiver.Type == runtime.DictKeyReading {
			s = "." + strings.Join(receiver.DictReadingPath, ".")
		}
		to = append(to, s)
	}

	return fmt.Sprintf(
		"%s -> %s",
		l.formatPortAddr(connection.SenderPortAddr),
		strings.Join(to, ", "),
	)
}

func (l LoggingInterceptor) formatPortAddr(addr runtime.AbsolutePortAddr) string {
	return fmt.Sprintf("%s.%s[%d]", addr.Path, addr.Port, addr.Idx)
}
