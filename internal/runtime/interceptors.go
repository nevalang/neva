package runtime

import (
	"fmt"
	"os"
	"strings"
)

type ProdInterceptor struct{}

func (ProdInterceptor) Sent(sender PortSlotAddr, msg Msg) Msg {
	return msg
}

func (ProdInterceptor) Received(receiver PortSlotAddr, msg Msg) Msg {
	return msg
}

type DebugInterceptor struct {
	logger Logger
}

type Logger interface {
	Printf(format string, v ...any) error
}

func NewDebugInterceptor(logger Logger) DebugInterceptor {
	return DebugInterceptor{logger: logger}
}

type FileLogger struct {
	filepath string
}

func (f FileLogger) Printf(format string, v ...any) error {
	file, err := os.OpenFile(f.filepath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = fmt.Fprintf(file, format, v...)
	return err
}

func (d DebugInterceptor) Sent(
	sender PortSlotAddr,
	msg Msg,
) Msg {
	d.logger.Printf(
		"sent | %v | %v\n",
		d.formatPortSlotAddr(sender), d.formatMsg(msg),
	)
	return msg
}

func (d DebugInterceptor) Received(
	receiver PortSlotAddr,
	msg Msg,
) Msg {
	d.logger.Printf(
		"recv | %v | %v\n",
		d.formatPortSlotAddr(receiver),
		d.formatMsg(msg),
	)

	return msg
}

func (d DebugInterceptor) formatMsg(msg Msg) string {
	if s, ok := msg.(StrMsg); ok {
		return fmt.Sprintf(`"%s"`, s)
	}
	return fmt.Sprint(msg)
}

func (d DebugInterceptor) formatPortSlotAddr(slotAddr PortSlotAddr) string {
	parts := strings.Split(slotAddr.Path, "/")
	lastPart := parts[len(parts)-1]
	if lastPart == "in" || lastPart == "out" {
		parts = parts[:len(parts)-1]
	}
	slotAddr.Path = strings.Join(parts, "/")

	s := fmt.Sprintf("%v:%v", slotAddr.Path, slotAddr.Port)
	if slotAddr.Index != nil {
		s = fmt.Sprintf("%v[%v]", s, *slotAddr.Index)
	}

	return s
}
