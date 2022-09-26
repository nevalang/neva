package log

import (
	"fmt"
	"log"
	"strings"

	"github.com/emil14/neva/internal/core"
	"github.com/emil14/neva/internal/runtime/src"
)

type LoggingInterceptor struct{}

func (l LoggingInterceptor) AfterSending(conn src.Connection, msg core.Msg) core.Msg {
	log.Printf("sent: %s", l.fmtConn(conn, msg))
	return msg
}

func (l LoggingInterceptor) BeforeReceiving(
	saddr src.AbsolutePortAddr,
	rpoint src.ReceiverConnectionPoint,
	msg core.Msg,
) core.Msg {
	log.Printf("prepare: %s <- %s <- %s", l.fmtPortAddr(rpoint.PortAddr), msg, l.fmtPortAddr(saddr))
	return msg
}

func (l LoggingInterceptor) AfterReceiving(
	saddr src.AbsolutePortAddr,
	rpoint src.ReceiverConnectionPoint,
	msg core.Msg,
) {
	log.Printf("received: %s <- %s <- %s", l.fmtPortAddr(rpoint.PortAddr), msg, l.fmtPortAddr(saddr))
}

func (l LoggingInterceptor) fmtConn(conn src.Connection, msg core.Msg) string {
	to := []string{}
	for _, r := range conn.ReceiversConnectionPoints {
		s := l.fmtPortAddr(r.PortAddr)
		if r.Type == src.DictReading {
			s = "." + strings.Join(r.DictReadingPath, ".")
		}
		to = append(to, s)
	}

	return fmt.Sprintf(
		"%s -> %s -> %s",
		l.fmtPortAddr(conn.SenderPortAddr),
		msg,
		strings.Join(to, ", "),
	)
}

func (l LoggingInterceptor) fmtPortAddr(addr src.AbsolutePortAddr) string {
	return fmt.Sprintf("%s.%s[%d]", addr.Path, addr.Port, addr.Idx)
}
