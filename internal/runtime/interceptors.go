package runtime

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

// traceEventVersion tracks JSONL schema version for emitted runtime events.
const traceEventVersion = 2

// EventKind is a runtime transport event kind.
type EventKind string

const (
	EventSent EventKind = "sent"
	EventRecv EventKind = "recv"
)

// SentEvent is emitted when runtime sends a message through an outport.
type SentEvent struct {
	Msg          Msg          `json:"message"`
	PortSlotAddr PortSlotAddr `json:"port"`
	EventKind    EventKind    `json:"event"`
	CauseIndexes []uint64     `json:"causeIndexes"`
	Version      int          `json:"v"`
	Index        uint64       `json:"index"`
}

// RecvEvent is emitted when runtime receives a message from an inport.
type RecvEvent struct {
	Msg          Msg          `json:"message"`
	PortSlotAddr PortSlotAddr `json:"port"`
	Event        EventKind    `json:"event"`
	Version      int          `json:"v"`
	Index        uint64       `json:"index"`
}

// NoEffectInterceptor exist to be used as default interceptor
// when no actual interception is needed. Just to satisfy compiler.
type NoEffectInterceptor struct{}

func (NoEffectInterceptor) Prepare() error { return nil }

//nolint:ireturn // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
func (p NoEffectInterceptor) Sent(
	_ context.Context,
	_ PortSlotAddr,
	ordered OrderedMsg,
	_ TraceHop,
) {
}

//nolint:ireturn // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
func (p NoEffectInterceptor) Received(_ context.Context, _ PortSlotAddr, ordered OrderedMsg) OrderedMsg {
	return ordered
}

// JSONLTraceFileWriter is an interceptor
// that writes each tracing event into a file in a JSONL format.
//
//nolint:recvcheck // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
type JSONLTraceFileWriter struct {
	file    *os.File
	comment string
}

// Open safely opens the file for writing and returns its close func.
func (d *JSONLTraceFileWriter) Open(filepath string) (func() error, error) {
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

// Sent implements Interceptor interface.
// The only thing this interceptor really does - is writes JSONL line to a file.
//
//nolint:ireturn // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
func (d *JSONLTraceFileWriter) Sent(
	_ context.Context,
	sender PortSlotAddr,
	ordered OrderedMsg,
	hop TraceHop,
) {
	evt := SentEvent{
		Version:      traceEventVersion,
		EventKind:    EventSent,
		Index:        ordered.index,
		CauseIndexes: hop.CauseIndexes,
		PortSlotAddr: d.normalizePortSlotAddrPath(sender),
		Msg:          ordered.Msg,
	}
	d.writeTraceEvent(evt)
}

// Received implements Interceptor interface.
// The only thing this interceptor really does - is writes JSONL line to a file.
//
//nolint:ireturn // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
func (d *JSONLTraceFileWriter) Received(_ context.Context, receiver PortSlotAddr, ordered OrderedMsg) OrderedMsg {
	evt := RecvEvent{
		Version:      traceEventVersion,
		Event:        EventRecv,
		Index:        ordered.index,
		PortSlotAddr: d.normalizePortSlotAddrPath(receiver),
		Msg:          ordered.Msg,
	}
	d.writeTraceEvent(evt)
	return ordered
}

func (d JSONLTraceFileWriter) writeTraceEvent(evt any) {
	encoded, err := json.Marshal(evt)
	if err != nil {
		panic(err)
	}
	if _, err := fmt.Fprintln(d.file, string(encoded)); err != nil {
		panic(err)
	}
}

func (d JSONLTraceFileWriter) normalizePortSlotAddrPath(slotAddr PortSlotAddr) PortSlotAddr {
	slotAddr.Path = d.normalizePortPath(slotAddr.Path)
	return slotAddr
}

func (d JSONLTraceFileWriter) normalizePortPath(path string) string {
	parts := strings.Split(path, "/")
	if len(parts) == 0 {
		return path
	}

	lastPart := parts[len(parts)-1]
	if lastPart == "in" || lastPart == "out" {
		parts = parts[:len(parts)-1]
	}

	return strings.Join(parts, "/")
}

func NewDebugInterceptor(comment string) *JSONLTraceFileWriter {
	return &JSONLTraceFileWriter{
		comment: comment,
	}
}

// NewInterceptor always enables in-memory tracing.
// When tracePath is empty, it skips JSONL emission and returns the production interceptor.
//
//nolint:ireturn // Interceptor is the stable runtime transport contract for generated programs.
func NewInterceptor(tracePath, comment string) (Interceptor, func() error, error) {
	if tracePath == "" {
		return NoEffectInterceptor{}, func() error { return nil }, nil
	}

	interceptor := &JSONLTraceFileWriter{
		comment: comment,
	}
	closeFn, err := interceptor.Open(tracePath)
	if err != nil {
		return nil, nil, err
	}
	return interceptor, closeFn, nil
}
