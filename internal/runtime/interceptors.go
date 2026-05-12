package runtime

import (
	"encoding/json"
	"fmt"
	"os"
)

// traceEventVersion tracks JSONL schema version for emitted runtime events.
const traceEventVersion = 1

// EventKind is a runtime transport event kind.
type EventKind string

const (
	EventSent EventKind = "sent"
	EventRecv EventKind = "recv"
)

// EventPort identifies a concrete runtime endpoint.
type EventPort struct {
	Index *uint8 `json:"index,omitempty"`
	Path  string `json:"path"`
	Name  string `json:"name"`
}

// SentEvent is emitted when runtime sends a message through an outport.
type SentEvent struct {
	Port          EventPort `json:"port"`
	Event         EventKind `json:"event"`
	Message       string    `json:"message"`
	Version       int       `json:"v"`
	TraceID       uint64    `json:"traceId"`
	ParentTraceID uint64    `json:"parentTraceId"`
}

// RecvEvent is emitted when runtime receives a message from an inport.
type RecvEvent struct {
	Port    EventPort `json:"port"`
	Event   EventKind `json:"event"`
	Message string    `json:"message"`
	Version int       `json:"v"`
	TraceID uint64    `json:"traceId"`
}

type ProdInterceptor struct{}

func (ProdInterceptor) Prepare() error { return nil }

//nolint:ireturn // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
func (ProdInterceptor) Sent(sender PortSlotAddr, ordered OrderedMsg) OrderedMsg {
	recordOrderedSent(sender, ordered)
	return ordered
}

//nolint:ireturn // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
func (ProdInterceptor) Received(receiver PortSlotAddr, ordered OrderedMsg) OrderedMsg {
	recordOrderedReceived(receiver, ordered)
	return ordered
}

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
	recordOrderedSent(sender, ordered)
	evt := SentEvent{
		Version:       traceEventVersion,
		Event:         EventSent,
		TraceID:       ordered.index,
		ParentTraceID: parentTraceIDFromMsg(ordered.Msg),
		Port:          eventPortFromSlot(sender),
		Message:       d.formatMsg(ordered.Msg),
	}
	writeTraceEvent(d.file, evt)
	return ordered
}

//nolint:ireturn // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
func (d *DebugInterceptor) Received(receiver PortSlotAddr, ordered OrderedMsg) OrderedMsg {
	recordOrderedReceived(receiver, ordered)
	evt := RecvEvent{
		Version: traceEventVersion,
		Event:   EventRecv,
		TraceID: ordered.index,
		Port:    eventPortFromSlot(receiver),
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

func eventPortFromSlot(slotAddr PortSlotAddr) EventPort {
	return EventPort{
		Path:  normalizePortPath(slotAddr.Path),
		Name:  slotAddr.Port,
		Index: slotAddr.Index,
	}
}

func NewDebugInterceptor(comment string) *DebugInterceptor {
	return &DebugInterceptor{comment: comment}
}
