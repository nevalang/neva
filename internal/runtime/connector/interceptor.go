package connector

import (
	"fmt"
	"log"
	"strings"

	"github.com/emil14/neva/internal/core"
	"github.com/emil14/neva/internal/runtime"
)

type DefaultInterceptor struct{}

func (d DefaultInterceptor) AfterSending(connection runtime.Connection, msg core.Msg) core.Msg {
	log.Println("AfterSending", d.formatConnection(connection), msg)
	return msg
}

func (d DefaultInterceptor) BeforeReceiving(from, to runtime.AbsolutePortAddr, msg core.Msg) core.Msg {
	// log.Println("BeforeReceiving", from, to, msg)
	return msg
}

func (d DefaultInterceptor) AfterReceiving(from, to runtime.AbsolutePortAddr, msg core.Msg) {
	log.Println("AfterReceiving", from, to, msg)
}

func (d DefaultInterceptor) formatConnection(connection runtime.Connection) string {
	to := []string{}
	for _, receiver := range connection.ReceiversConnectionPoints {
		to = append(to, receiver.PortAddr.String())
	}

	return fmt.Sprintf(
		"%s -> %s",
		connection.SenderPortAddr,
		strings.Join(to, ", "),
	)
}
