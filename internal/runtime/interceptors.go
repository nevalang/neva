package runtime

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

type traceSentEvent struct {
	Message       string `json:"message"`
	Event         string `json:"event"`
	Port          string `json:"port"`
	ParentTraceID uint64 `json:"parentTraceId"`
	TraceID       uint64 `json:"traceId"`
	Version       int    `json:"v"`
}

type traceRecvEvent struct {
	Message string `json:"message"`
	Event   string `json:"event"`
	Port    string `json:"port"`
	TraceID uint64 `json:"traceId"`
	Version int    `json:"v"`
}

type ProdInterceptor struct{}

func (ProdInterceptor) Prepare() error { return nil }

//nolint:ireturn // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
func (ProdInterceptor) Sent(sender PortSlotAddr, msg Msg) Msg { return msg }

//nolint:ireturn // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
func (ProdInterceptor) Received(receiver PortSlotAddr, msg Msg) Msg { return msg }

//nolint:recvcheck // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
type DebugInterceptor struct {
	file    *os.File
	comment string
}

func (d *DebugInterceptor) Open(filepath string) (func() error, error) {
	file, err := os.OpenFile(filepath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC|os.O_APPEND, 0644)
	if err != nil {
		//nolint:wrapcheck // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
		return nil, err
	}
	d.file = file
	if _, err := fmt.Fprintln(d.file, d.comment); err != nil {
		//nolint:wrapcheck // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
		return nil, err
	}
	return file.Close, nil
}

//nolint:ireturn // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
func (d *DebugInterceptor) Sent(sender PortSlotAddr, msg Msg) Msg {
	evt := traceSentEvent{
		Version:       1,
		Event:         "sent",
		TraceID:       mustTraceIDFromMsg(msg),
		ParentTraceID: parentTraceIDFromMsg(msg),
		Port:          d.formatPortSlotAddr(sender),
		Message:       d.formatMsg(msg),
	}
	writeTraceEvent(d.file, evt)
	return msg
}

//nolint:ireturn // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
func (d *DebugInterceptor) Received(receiver PortSlotAddr, msg Msg) Msg {
	evt := traceRecvEvent{
		Version: 1,
		Event:   "recv",
		TraceID: mustTraceIDFromMsg(msg),
		Port:    d.formatPortSlotAddr(receiver),
		Message: d.formatMsg(msg),
	}
	writeTraceEvent(d.file, evt)
	return msg
}

func (d DebugInterceptor) formatMsg(msg Msg) string {
	if strMsg, ok := msg.(StringMsg); ok {
		return fmt.Sprintf("%q", strMsg.Str())
	}
	if bytesMsg, ok := msg.(BytesMsg); ok {
		return fmt.Sprintf("%q", bytesMsg.Bytes())
	}
	return fmt.Sprint(msg)
}

func writeTraceEvent(file *os.File, evt any) {
	encoded, err := json.Marshal(evt)
	if err != nil {
		panic(err)
	}
	if _, err := fmt.Fprintln(file, string(encoded)); err != nil {
		panic(err)
	}
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

func NewDebugInterceptor(comment string) *DebugInterceptor {
	return &DebugInterceptor{comment: comment}
}
