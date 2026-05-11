package runtime

import (
	"encoding/json"
	"fmt"
	"os"
)

// traceEventVersion tracks JSONL schema version for trace events.
const traceEventVersion = 1

// traceEventPort identifies a concrete runtime port endpoint.
type traceEventPort struct {
	Index *uint8 `json:"index,omitempty"`
	Path  string `json:"path"`
	Name  string `json:"name"`
}

// traceSentEvent is emitted when runtime sends a message through outport.
type traceSentEvent struct {
	Port          traceEventPort `json:"port"`
	Event         string         `json:"event"`
	Message       string         `json:"message"`
	Version       int            `json:"v"`
	TraceID       uint64         `json:"traceId"`
	ParentTraceID uint64         `json:"parentTraceId"`
}

// traceRecvEvent is emitted when runtime receives a message from inport.
type traceRecvEvent struct {
	Port    traceEventPort `json:"port"`
	Event   string         `json:"event"`
	Message string         `json:"message"`
	Version int            `json:"v"`
	TraceID uint64         `json:"traceId"`
}

type ProdInterceptor struct{}

func (ProdInterceptor) Prepare() error { return nil }

//nolint:ireturn // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
func (ProdInterceptor) Sent(sender PortSlotAddr, ordered OrderedMsg) OrderedMsg { return ordered }

//nolint:ireturn // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
func (ProdInterceptor) Received(receiver PortSlotAddr, ordered OrderedMsg) OrderedMsg { return ordered }

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
func (d *DebugInterceptor) Sent(sender PortSlotAddr, ordered OrderedMsg) OrderedMsg {
	evt := traceSentEvent{
		Version:       traceEventVersion,
		Event:         "sent",
		TraceID:       ordered.index,
		ParentTraceID: parentTraceIDFromMsg(ordered.Msg),
		Port:          traceEventPortFromSlot(sender),
		Message:       d.formatMsg(ordered.Msg),
	}
	writeTraceEvent(d.file, evt)
	return ordered
}

//nolint:ireturn // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
func (d *DebugInterceptor) Received(receiver PortSlotAddr, ordered OrderedMsg) OrderedMsg {
	evt := traceRecvEvent{
		Version: traceEventVersion,
		Event:   "recv",
		TraceID: ordered.index,
		Port:    traceEventPortFromSlot(receiver),
		Message: d.formatMsg(ordered.Msg),
	}
	writeTraceEvent(d.file, evt)
	return ordered
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

func traceEventPortFromSlot(slotAddr PortSlotAddr) traceEventPort {
	return traceEventPort{
		Path:  normalizePortPath(slotAddr.Path),
		Name:  slotAddr.Port,
		Index: slotAddr.Index,
	}
}

func NewDebugInterceptor(comment string) *DebugInterceptor {
	return &DebugInterceptor{comment: comment}
}
