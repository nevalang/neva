package runtime

import (
	"fmt"
	"strings"
	"sync"
)

type TraceHop struct {
	Sender        *PortSlotAddr
	Receiver      *PortSlotAddr
	Message       string
	TraceID       uint64
	ParentTraceID uint64
}

type traceStore struct {
	hops map[uint64]TraceHop
	mu   sync.RWMutex
}

//nolint:gochecknoglobals // runtime-wide trace store must be shared across port events.
var globalTraceStore = traceStore{
	hops: make(map[uint64]TraceHop),
}

func recordSent(traceID, parentTraceID uint64, sender PortSlotAddr, msg Msg) {
	globalTraceStore.mu.Lock()
	defer globalTraceStore.mu.Unlock()

	hop := globalTraceStore.hops[traceID]
	hop.TraceID = traceID
	hop.ParentTraceID = parentTraceID
	senderCopy := sender
	hop.Sender = &senderCopy
	hop.Message = fmt.Sprint(UnwrapTraceMsg(msg))
	globalTraceStore.hops[traceID] = hop
}

func recordReceived(traceID uint64, receiver PortSlotAddr) {
	globalTraceStore.mu.Lock()
	defer globalTraceStore.mu.Unlock()

	hop := globalTraceStore.hops[traceID]
	hop.TraceID = traceID
	receiverCopy := receiver
	hop.Receiver = &receiverCopy
	globalTraceStore.hops[traceID] = hop
}

func traceHopByID(traceID uint64) (TraceHop, bool) {
	globalTraceStore.mu.RLock()
	defer globalTraceStore.mu.RUnlock()

	hop, ok := globalTraceStore.hops[traceID]
	return hop, ok
}

// TracePathByID reconstructs message ancestry from newest to oldest.
func TracePathByID(traceID uint64) []TraceHop {
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

		hop, ok := traceHopByID(cur)
		if !ok {
			break
		}
		path = append(path, hop)
		cur = hop.ParentTraceID
	}

	return path
}

// TracePath reconstructs message ancestry from newest to oldest.
func TracePath(msg Msg) []TraceHop {
	traceID, ok := TraceIDFromMsg(msg)
	if !ok {
		return nil
	}
	return TracePathByID(traceID)
}

// FormatDataflowTrace renders Dataflow Trace in graph terms from newest to oldest.
func FormatDataflowTrace(msg Msg) string {
	path := TracePath(msg)
	if len(path) == 0 {
		return ""
	}

	var builder strings.Builder
	builder.WriteString("panic cause dataflow trace (newest -> oldest):\n")

	for i := range path {
		hop := path[i]
		indent := strings.Repeat("  ", i)
		builder.WriteString(indent)
		builder.WriteString(formatHop(hop))
		builder.WriteByte('\n')
	}

	return strings.TrimRight(builder.String(), "\n")
}

func formatHop(hop TraceHop) string {
	recv := "<?>"
	send := "<?>"
	if hop.Receiver != nil {
		recv = formatPortSlotAddr(*hop.Receiver)
	}
	if hop.Sender != nil {
		send = formatPortSlotAddr(*hop.Sender)
	}
	return fmt.Sprintf("%s <- %s", recv, send)
}

func formatPortSlotAddr(slot PortSlotAddr) string {
	s := fmt.Sprintf("%s:%s", slot.Path, slot.Port)
	if slot.Index != nil {
		s = fmt.Sprintf("%s[%d]", s, *slot.Index)
	}
	return s
}

func resetTraceStoreForTests() {
	globalTraceStore.mu.Lock()
	defer globalTraceStore.mu.Unlock()
	globalTraceStore.hops = make(map[uint64]TraceHop)
}
