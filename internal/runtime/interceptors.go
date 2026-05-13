package runtime

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"slices"
)

// traceEventVersion tracks JSONL schema version for emitted runtime events.
const traceEventVersion = 2

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
	Port         EventPort `json:"port"`
	Event        EventKind `json:"event"`
	Message      string    `json:"message"`
	CauseIndexes []uint64  `json:"causeIndexes"`
	Version      int       `json:"v"`
	Index        uint64    `json:"index"`
}

// RecvEvent is emitted when runtime receives a message from an inport.
type RecvEvent struct {
	Port    EventPort `json:"port"`
	Event   EventKind `json:"event"`
	Message string    `json:"message"`
	Version int       `json:"v"`
	Index   uint64    `json:"index"`
}

type ProdInterceptor struct {
	tracer *Tracer
}

func (ProdInterceptor) Prepare() error { return nil }

func (p ProdInterceptor) getTracer() *Tracer {
	if p.tracer != nil {
		return p.tracer
	}
	return globalTracer
}

//nolint:ireturn // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
func (p ProdInterceptor) Sent(_ context.Context, sender PortSlotAddr, ordered OrderedMsg) OrderedMsg {
	p.getTracer().RecordSent(sender, ordered)
	return ordered
}

//nolint:ireturn // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
func (p ProdInterceptor) Received(ctx context.Context, receiver PortSlotAddr, ordered OrderedMsg) OrderedMsg {
	p.getTracer().RecordReceived(receiver, ordered)
	recordTraceReceive(ctx, ordered)
	return ordered
}

//nolint:recvcheck // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
type DebugInterceptor struct {
	tracer  *Tracer
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

func (d DebugInterceptor) getTracer() *Tracer {
	if d.tracer != nil {
		return d.tracer
	}
	return globalTracer
}

//nolint:ireturn // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
func (d *DebugInterceptor) Sent(_ context.Context, sender PortSlotAddr, ordered OrderedMsg) OrderedMsg {
	d.getTracer().RecordSent(sender, ordered)
	evt := SentEvent{
		Version:      traceEventVersion,
		Event:        EventSent,
		Index:        ordered.index,
		CauseIndexes: slices.Clone(ordered.causeIndexes),
		Port:         eventPortFromSlot(sender),
		Message:      d.formatMsg(ordered.Msg),
	}
	writeTraceEvent(d.file, evt)
	return ordered
}

//nolint:ireturn // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
func (d *DebugInterceptor) Received(ctx context.Context, receiver PortSlotAddr, ordered OrderedMsg) OrderedMsg {
	d.getTracer().RecordReceived(receiver, ordered)
	recordTraceReceive(ctx, ordered)
	evt := RecvEvent{
		Version: traceEventVersion,
		Event:   EventRecv,
		Index:   ordered.index,
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
	return &DebugInterceptor{
		tracer:  globalTracer,
		comment: comment,
	}
}

// NewInterceptor always enables in-memory tracing.
// When tracePath is empty, it skips JSONL emission and returns the production interceptor.
//
//nolint:ireturn // Interceptor is the stable runtime transport contract for generated programs.
func NewInterceptor(tracePath, comment string) (Interceptor, func() error, error) {
	tracer := globalTracer
	if tracePath == "" {
		return ProdInterceptor{tracer: tracer}, func() error { return nil }, nil
	}

	interceptor := &DebugInterceptor{
		tracer:  tracer,
		comment: comment,
	}
	closeFn, err := interceptor.Open(tracePath)
	if err != nil {
		return nil, nil, err
	}
	return interceptor, closeFn, nil
}
