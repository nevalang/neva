package runtime

import (
	"fmt"
	"strings"
	"sync"
)

type traceCarrier interface {
	traceMeta() (uint64, uint64)
}

type TraceHop struct {
	Sender        *PortSlotAddr
	Receiver      *PortSlotAddr
	Message       string
	TraceID       uint64
	ParentTraceID uint64
}

type Tracer struct {
	store traceStore
}

type traceStore struct {
	hops map[uint64]TraceHop
	mu   sync.RWMutex
}

//nolint:gochecknoglobals // runtime-wide tracer must be shared across port events.
var globalTracer = NewTracer()

func NewTracer() *Tracer {
	return &Tracer{
		store: traceStore{
			hops: make(map[uint64]TraceHop),
		},
	}
}

// TraceIDFromMsg extracts runtime Dataflow Trace identity from message payload.
func TraceIDFromMsg(msg Msg) (uint64, bool) {
	carrier, ok := msg.(traceCarrier)
	if !ok {
		return 0, false
	}
	traceID, _ := carrier.traceMeta()
	return traceID, traceID != 0
}

func parentTraceIDFromMsg(msg Msg) uint64 {
	carrier, ok := msg.(traceCarrier)
	if !ok {
		return 0
	}
	_, parentTraceID := carrier.traceMeta()
	return parentTraceID
}

//nolint:ireturn // Msg is runtime contract type.
func withTrace(msg Msg, traceID, parentTraceID uint64) Msg {
	switch typedMsg := msg.(type) {
	case BoolMsg:
		typedMsg.traceID = traceID
		typedMsg.parentTraceID = parentTraceID
		return typedMsg
	case IntMsg:
		typedMsg.traceID = traceID
		typedMsg.parentTraceID = parentTraceID
		return typedMsg
	case FloatMsg:
		typedMsg.traceID = traceID
		typedMsg.parentTraceID = parentTraceID
		return typedMsg
	case StringMsg:
		typedMsg.traceID = traceID
		typedMsg.parentTraceID = parentTraceID
		return typedMsg
	case BytesMsg:
		typedMsg.traceID = traceID
		typedMsg.parentTraceID = parentTraceID
		return typedMsg
	case ListMsg:
		typedMsg.traceID = traceID
		typedMsg.parentTraceID = parentTraceID
		return typedMsg
	case DictMsg:
		typedMsg.traceID = traceID
		typedMsg.parentTraceID = parentTraceID
		return typedMsg
	case StructMsg:
		typedMsg.traceID = traceID
		typedMsg.parentTraceID = parentTraceID
		return typedMsg
	case UnionMsg:
		typedMsg.traceID = traceID
		typedMsg.parentTraceID = parentTraceID
		return typedMsg
	default:
		return msg
	}
}

func (t *Tracer) RecordSent(sender PortSlotAddr, ordered OrderedMsg) {
	t.store.mu.Lock()
	defer t.store.mu.Unlock()

	traceID := ordered.index
	parentTraceID := parentTraceIDFromMsg(ordered.Msg)
	hop := t.store.hops[traceID]
	hop.TraceID = traceID
	hop.ParentTraceID = parentTraceID
	senderCopy := sender
	hop.Sender = &senderCopy
	hop.Message = fmt.Sprint(ordered.Msg)
	t.store.hops[traceID] = hop
}

func (t *Tracer) RecordReceived(receiver PortSlotAddr, ordered OrderedMsg) {
	t.store.mu.Lock()
	defer t.store.mu.Unlock()

	hop := t.store.hops[ordered.index]
	hop.TraceID = ordered.index
	receiverCopy := receiver
	hop.Receiver = &receiverCopy
	t.store.hops[ordered.index] = hop
}

func (t *Tracer) traceHopByID(traceID uint64) (TraceHop, bool) {
	t.store.mu.RLock()
	defer t.store.mu.RUnlock()

	hop, ok := t.store.hops[traceID]
	return hop, ok
}

// TracePathByID reconstructs message ancestry from newest to oldest.
func (t *Tracer) TracePathByID(traceID uint64) []TraceHop {
	if traceID == 0 {
		return nil
	}

	path := make([]TraceHop, 0, 8)
	seen := make(map[uint64]struct{}, 8)
	cur := traceID

	for cur != 0 {
		if _, ok := seen[cur]; ok {
			break
		}
		seen[cur] = struct{}{}

		hop, ok := t.traceHopByID(cur)
		if !ok {
			break
		}
		path = append(path, hop)
		cur = hop.ParentTraceID
	}

	return path
}

// TracePath reconstructs message ancestry from newest to oldest.
func (t *Tracer) TracePath(msg Msg) []TraceHop {
	traceID, ok := TraceIDFromMsg(msg)
	if !ok {
		return nil
	}
	return t.TracePathByID(traceID)
}

// FormatDataflowTrace renders panic-focused Dataflow Trace in a readable flow format.
func (t *Tracer) FormatDataflowTrace(msg Msg) string {
	path := t.TracePath(msg)
	if len(path) == 0 {
		return ""
	}

	var builder strings.Builder
	panicReceiver := "<?>"
	if path[0].Receiver != nil {
		panicReceiver = formatPortSlotAddr(*path[0].Receiver)
	}
	panicComponent := componentNameFromReceiver(path[0].Receiver)

	builder.WriteString("panic cause dataflow trace\n")
	builder.WriteString("direction: oldest -> newest (left -> right)\n")
	_, _ = fmt.Fprintf(&builder, "panic sink: %s\n", panicReceiver)
	if panicComponent != "" {
		_, _ = fmt.Fprintf(&builder, "panic component: %s\n", panicComponent)
	}
	builder.WriteString("events:\n")

	for i := len(path) - 1; i >= 0; i-- {
		hop := path[i]
		_, _ = fmt.Fprintf(&builder, "  %d. %s\n", len(path)-i, formatHopFlow(hop))
		builder.WriteByte('\n')
	}

	return strings.TrimRight(builder.String(), "\n")
}

func formatHopFlow(hop TraceHop) string {
	recv := "<?>"
	send := "<?>"
	if hop.Receiver != nil {
		recv = formatPortSlotAddr(*hop.Receiver)
	}
	if hop.Sender != nil {
		send = formatPortSlotAddr(*hop.Sender)
	}
	return fmt.Sprintf("%s -> %s", send, recv)
}

func componentNameFromReceiver(receiver *PortSlotAddr) string {
	if receiver == nil {
		return ""
	}
	path := receiver.Path
	path = strings.TrimSuffix(path, "/in")
	path = strings.TrimSuffix(path, "/out")
	parts := strings.Split(path, "/")
	if len(parts) == 0 {
		return path
	}
	return parts[len(parts)-1]
}

func formatPortSlotAddr(slot PortSlotAddr) string {
	slot.Path = normalizePortPath(slot.Path)
	s := fmt.Sprintf("%s:%s", slot.Path, slot.Port)
	if slot.Index != nil {
		s = fmt.Sprintf("%s[%d]", s, *slot.Index)
	}
	return s
}

func normalizePortPath(path string) string {
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

func AsUnion(msg Msg) (UnionMsg, bool) {
	unionMsg, ok := msg.(UnionMsg)
	return unionMsg, ok
}

func TracePathByID(traceID uint64) []TraceHop {
	return globalTracer.TracePathByID(traceID)
}

func TracePath(msg Msg) []TraceHop {
	return globalTracer.TracePath(msg)
}

func FormatDataflowTrace(msg Msg) string {
	return globalTracer.FormatDataflowTrace(msg)
}

func resetTraceStoreForTests() {
	globalTracer.store.mu.Lock()
	defer globalTracer.store.mu.Unlock()
	globalTracer.store.hops = make(map[uint64]TraceHop)
}
