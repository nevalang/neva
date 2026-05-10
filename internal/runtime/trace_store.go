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
	hop.Message = fmt.Sprint(msg)
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

// FormatDataflowTrace renders panic-focused Dataflow Trace in a readable flow format.
func FormatDataflowTrace(msg Msg) string {
	path := TracePath(msg)
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
	builder.WriteString(fmt.Sprintf("panic sink: %s\n", panicReceiver))
	if panicComponent != "" {
		builder.WriteString(fmt.Sprintf("panic component: %s\n", panicComponent))
	}
	builder.WriteString("events:\n")

	for i := len(path) - 1; i >= 0; i-- {
		hop := path[i]
		builder.WriteString(fmt.Sprintf("  %d. %s\n", len(path)-i, formatHopFlow(hop)))
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

func resetTraceStoreForTests() {
	globalTraceStore.mu.Lock()
	defer globalTraceStore.mu.Unlock()
	globalTraceStore.hops = make(map[uint64]TraceHop)
}
