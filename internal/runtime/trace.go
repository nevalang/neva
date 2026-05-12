package runtime

import (
	"context"
	"fmt"
	"slices"
	"strings"
	"sync"
)

type traceCarrier interface {
	traceMeta() (uint64, []uint64)
}

type TraceHop struct {
	Sender         *PortSlotAddr
	Receiver       *PortSlotAddr
	Message        string
	ParentTraceIDs []uint64
	TraceID        uint64
}

type TraceTree struct {
	Parents []TraceTree
	Hop     TraceHop
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

func parentTraceIDsFromMsg(msg Msg) []uint64 {
	carrier, ok := msg.(traceCarrier)
	if !ok {
		return nil
	}
	_, parentTraceIDs := carrier.traceMeta()
	return slices.Clone(parentTraceIDs)
}

func parentTraceIDsFromOrdered(causes []OrderedMsg) []uint64 {
	if len(causes) == 0 {
		return nil
	}

	parentTraceIDs := make([]uint64, 0, len(causes))
	seen := make(map[uint64]struct{}, len(causes))
	for _, cause := range causes {
		if cause.index == 0 {
			continue
		}
		if _, ok := seen[cause.index]; ok {
			continue
		}
		seen[cause.index] = struct{}{}
		parentTraceIDs = append(parentTraceIDs, cause.index)
	}
	slices.Sort(parentTraceIDs)
	return parentTraceIDs
}

//nolint:ireturn // Msg is runtime contract type.
func withTrace(msg Msg, traceID uint64, parentTraceIDs []uint64) Msg {
	switch typedMsg := msg.(type) {
	case BoolMsg:
		typedMsg.traceID = traceID
		typedMsg.parentTraceIDs = slices.Clone(parentTraceIDs)
		return typedMsg
	case IntMsg:
		typedMsg.traceID = traceID
		typedMsg.parentTraceIDs = slices.Clone(parentTraceIDs)
		return typedMsg
	case FloatMsg:
		typedMsg.traceID = traceID
		typedMsg.parentTraceIDs = slices.Clone(parentTraceIDs)
		return typedMsg
	case StringMsg:
		typedMsg.traceID = traceID
		typedMsg.parentTraceIDs = slices.Clone(parentTraceIDs)
		return typedMsg
	case BytesMsg:
		typedMsg.traceID = traceID
		typedMsg.parentTraceIDs = slices.Clone(parentTraceIDs)
		return typedMsg
	case ListMsg:
		typedMsg.traceID = traceID
		typedMsg.parentTraceIDs = slices.Clone(parentTraceIDs)
		return typedMsg
	case DictMsg:
		typedMsg.traceID = traceID
		typedMsg.parentTraceIDs = slices.Clone(parentTraceIDs)
		return typedMsg
	case StructMsg:
		typedMsg.traceID = traceID
		typedMsg.parentTraceIDs = slices.Clone(parentTraceIDs)
		return typedMsg
	case UnionMsg:
		typedMsg.traceID = traceID
		typedMsg.parentTraceIDs = slices.Clone(parentTraceIDs)
		return typedMsg
	default:
		return msg
	}
}

type traceActivationKey struct{}

type traceActivationState struct {
	parents map[uint64]struct{}
	mu      sync.Mutex
	emitted bool
}

func contextWithTraceActivation(ctx context.Context) context.Context {
	return context.WithValue(ctx, traceActivationKey{}, &traceActivationState{
		parents: map[uint64]struct{}{},
	})
}

func traceActivationFromContext(ctx context.Context) *traceActivationState {
	state, ok := ctx.Value(traceActivationKey{}).(*traceActivationState)
	if !ok {
		return nil
	}
	return state
}

func recordTraceReceive(ctx context.Context, ordered OrderedMsg) {
	state := traceActivationFromContext(ctx)
	if state == nil {
		return
	}

	state.mu.Lock()
	defer state.mu.Unlock()

	if state.emitted {
		clear(state.parents)
		state.emitted = false
	}

	if ordered.index != 0 {
		state.parents[ordered.index] = struct{}{}
	}
}

func currentTraceParents(ctx context.Context, msg Msg) []uint64 {
	state := traceActivationFromContext(ctx)
	if state == nil {
		parentTraceIDs := parentTraceIDsFromMsg(msg)
		if len(parentTraceIDs) == 0 {
			if traceID, ok := TraceIDFromMsg(msg); ok && traceID != 0 {
				return []uint64{traceID}
			}
		}
		return parentTraceIDs
	}

	state.mu.Lock()
	defer state.mu.Unlock()

	if len(state.parents) != 0 {
		parentTraceIDs := make([]uint64, 0, len(state.parents))
		for traceID := range state.parents {
			parentTraceIDs = append(parentTraceIDs, traceID)
		}
		slices.Sort(parentTraceIDs)
		state.emitted = true
		return parentTraceIDs
	}

	parentTraceIDs := parentTraceIDsFromMsg(msg)
	if len(parentTraceIDs) == 0 {
		if traceID, ok := TraceIDFromMsg(msg); ok && traceID != 0 {
			parentTraceIDs = []uint64{traceID}
		}
	}
	state.emitted = true
	return parentTraceIDs
}

func parentTraceIDsForSend(ctx context.Context, msg Msg, causes []OrderedMsg) []uint64 {
	if len(causes) != 0 {
		return parentTraceIDsFromOrdered(causes)
	}
	return currentTraceParents(ctx, msg)
}

func (t *Tracer) RecordSent(sender PortSlotAddr, ordered OrderedMsg) {
	t.store.mu.Lock()
	defer t.store.mu.Unlock()

	traceID := ordered.index
	parentTraceIDs := parentTraceIDsFromMsg(ordered.Msg)
	hop := t.store.hops[traceID]
	hop.TraceID = traceID
	hop.ParentTraceIDs = slices.Clone(parentTraceIDs)
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

func (t *Tracer) traceTreeByID(traceID uint64, visited map[uint64]struct{}) (TraceTree, bool) {
	if traceID == 0 {
		return TraceTree{}, false
	}
	if _, seen := visited[traceID]; seen {
		return TraceTree{}, false
	}
	visited[traceID] = struct{}{}

	hop, ok := t.traceHopByID(traceID)
	if !ok {
		return TraceTree{}, false
	}

	tree := TraceTree{
		Hop:     hop,
		Parents: make([]TraceTree, 0, len(hop.ParentTraceIDs)),
	}
	for _, parentTraceID := range hop.ParentTraceIDs {
		parentTree, ok := t.traceTreeByID(parentTraceID, visited)
		if !ok {
			continue
		}
		tree.Parents = append(tree.Parents, parentTree)
	}

	delete(visited, traceID)
	return tree, true
}

// TraceCauseTreeByID reconstructs full multi-parent ancestry for the given trace id.
func (t *Tracer) TraceCauseTreeByID(traceID uint64) (TraceTree, bool) {
	return t.traceTreeByID(traceID, map[uint64]struct{}{})
}

// TraceCauseTree reconstructs full multi-parent ancestry for the given message.
func (t *Tracer) TraceCauseTree(msg Msg) (TraceTree, bool) {
	traceID, hasTrace := TraceIDFromMsg(msg)
	if !hasTrace {
		return TraceTree{}, false
	}
	return t.TraceCauseTreeByID(traceID)
}

// FormatDataflowTrace renders panic-focused Dataflow Trace in a readable flow format.
func (t *Tracer) FormatDataflowTrace(msg Msg) string {
	tree, ok := t.TraceCauseTree(msg)
	if !ok {
		return ""
	}

	var builder strings.Builder
	panicReceiver := "<?>"
	if tree.Hop.Receiver != nil {
		panicReceiver = formatPortSlotAddr(*tree.Hop.Receiver)
	}
	panicComponent := componentNameFromReceiver(tree.Hop.Receiver)

	builder.WriteString("panic cause dataflow trace\n")
	builder.WriteString("direction: newest -> oldest (top -> bottom)\n")
	_, _ = fmt.Fprintf(&builder, "panic sink: %s\n", panicReceiver)
	if panicComponent != "" {
		_, _ = fmt.Fprintf(&builder, "panic component: %s\n", panicComponent)
	}
	builder.WriteString("events:\n")
	formatTraceTree(&builder, &tree, "  ")

	return strings.TrimRight(builder.String(), "\n")
}

func formatTraceTree(builder *strings.Builder, tree *TraceTree, indent string) {
	_, _ = fmt.Fprintf(builder, "%s- %s\n", indent, formatHopFlow(tree.Hop))
	for _, parent := range tree.Parents {
		formatTraceTree(builder, &parent, indent+"  ")
	}
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

func TraceCauseTreeByID(traceID uint64) (TraceTree, bool) {
	return globalTracer.TraceCauseTreeByID(traceID)
}

func TraceCauseTree(msg Msg) (TraceTree, bool) {
	return globalTracer.TraceCauseTree(msg)
}

func FormatDataflowTrace(msg Msg) string {
	return globalTracer.FormatDataflowTrace(msg)
}

func resetTraceStoreForTests() {
	globalTracer.store.mu.Lock()
	defer globalTracer.store.mu.Unlock()
	globalTracer.store.hops = make(map[uint64]TraceHop)
}
