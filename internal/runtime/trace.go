package runtime

import (
	"context"
	"fmt"
	"slices"
	"strings"
	"sync"
)

type TraceHop struct {
	Sender   *PortSlotAddr
	Receiver *PortSlotAddr
	Message  string
	// CauseIndexes is the single source of truth for causal edges in traceStore.
	// Tree views are reconstructed from these indexes on demand.
	CauseIndexes []uint64
	Index        uint64
}

type TraceTree struct {
	// Parents is a derived, read-only projection rebuilt from traceStore hop links.
	// It is intentionally denormalized for traversal/formatting APIs.
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

func NewTracer() *Tracer {
	return &Tracer{
		store: traceStore{
			hops: make(map[uint64]TraceHop),
		},
	}
}

func causeIndexesFromOrdered(causes []OrderedMsg) []uint64 {
	if len(causes) == 0 {
		return nil
	}

	indexes := make([]uint64, 0, len(causes))
	seen := make(map[uint64]struct{}, len(causes))
	for _, cause := range causes {
		if cause.index == 0 {
			continue
		}
		if _, ok := seen[cause.index]; ok {
			continue
		}
		seen[cause.index] = struct{}{}
		indexes = append(indexes, cause.index)
	}
	slices.Sort(indexes)
	return indexes
}

//nolint:ireturn // Msg is runtime contract type.
func orderedPayload(msg Msg) Msg {
	if ordered, ok := msg.(OrderedMsg); ok {
		return ordered.Msg
	}
	return msg
}

type traceActivationKey struct{}
type tracerKey struct{}

type traceActivationState struct {
	causes  map[uint64]struct{}
	mu      sync.Mutex
	emitted bool
}

func contextWithTraceActivation(ctx context.Context) context.Context {
	return context.WithValue(ctx, traceActivationKey{}, &traceActivationState{
		causes: map[uint64]struct{}{},
	})
}

func contextWithTracer(ctx context.Context, tracer *Tracer) context.Context {
	return context.WithValue(ctx, tracerKey{}, tracer)
}

func WithTracer(ctx context.Context) context.Context {
	return contextWithTracer(ctx, NewTracer())
}

func tracerFromContext(ctx context.Context) (*Tracer, bool) {
	tracer, ok := ctx.Value(tracerKey{}).(*Tracer)
	if !ok || tracer == nil {
		return nil, false
	}
	return tracer, true
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
		clear(state.causes)
		state.emitted = false
	}

	if ordered.index != 0 {
		state.causes[ordered.index] = struct{}{}
	}
}

func currentCauseIndexes(ctx context.Context) []uint64 {
	state := traceActivationFromContext(ctx)
	if state != nil {
		state.mu.Lock()
		defer state.mu.Unlock()

		if len(state.causes) != 0 {
			indexes := make([]uint64, 0, len(state.causes))
			for index := range state.causes {
				indexes = append(indexes, index)
			}
			slices.Sort(indexes)
			state.emitted = true
			return indexes
		}

		state.emitted = true
	}

	return nil
}

func causeIndexesForSend(ctx context.Context, causes []OrderedMsg) []uint64 {
	if len(causes) != 0 {
		return causeIndexesFromOrdered(causes)
	}
	return currentCauseIndexes(ctx)
}

func (t *Tracer) RecordSent(
	ctx context.Context,
	sender PortSlotAddr,
	ordered OrderedMsg,
	causes []OrderedMsg,
) TraceHop {
	t.store.mu.Lock()
	defer t.store.mu.Unlock()

	hop := t.store.hops[ordered.index]
	hop.Index = ordered.index
	hop.CauseIndexes = causeIndexesForSend(ctx, causes)
	senderCopy := sender
	hop.Sender = &senderCopy
	hop.Message = fmt.Sprint(ordered.Msg)
	t.store.hops[ordered.index] = hop

	return hop
}

func (t *Tracer) RecordReceived(receiver PortSlotAddr, ordered OrderedMsg) {
	t.store.mu.Lock()
	defer t.store.mu.Unlock()

	hop := t.store.hops[ordered.index]
	hop.Index = ordered.index
	receiverCopy := receiver
	hop.Receiver = &receiverCopy
	t.store.hops[ordered.index] = hop
}

func (t *Tracer) traceHopByIndex(index uint64) (TraceHop, bool) {
	t.store.mu.RLock()
	defer t.store.mu.RUnlock()

	hop, ok := t.store.hops[index]
	return hop, ok
}

func (t *Tracer) traceTreeByIndex(index uint64, visited map[uint64]struct{}) (TraceTree, bool) {
	if index == 0 {
		return TraceTree{}, false
	}
	if _, seen := visited[index]; seen {
		return TraceTree{}, false
	}
	visited[index] = struct{}{}

	hop, ok := t.traceHopByIndex(index)
	if !ok {
		return TraceTree{}, false
	}

	// Rebuild tree view from normalized hop links; no second persisted source of truth.
	tree := TraceTree{
		Hop:     hop,
		Parents: make([]TraceTree, 0, len(hop.CauseIndexes)),
	}
	for _, causeIndex := range hop.CauseIndexes {
		parentTree, ok := t.traceTreeByIndex(causeIndex, visited)
		if !ok {
			continue
		}
		tree.Parents = append(tree.Parents, parentTree)
	}

	delete(visited, index)
	return tree, true
}

// TraceCauseTreeByIndex reconstructs full multi-parent ancestry for the given message index.
func (t *Tracer) TraceCauseTreeByIndex(index uint64) (TraceTree, bool) {
	return t.traceTreeByIndex(index, map[uint64]struct{}{})
}

// TraceCauseTree reconstructs full multi-parent ancestry for the given ordered message.
func (t *Tracer) TraceCauseTree(ordered OrderedMsg) (TraceTree, bool) {
	if ordered.index == 0 {
		return TraceTree{}, false
	}
	return t.TraceCauseTreeByIndex(ordered.index)
}

func AsUnion(msg Msg) (UnionMsg, bool) {
	unionMsg, ok := orderedPayload(msg).(UnionMsg)
	return unionMsg, ok
}

func TraceCauseTreeByIndex(ctx context.Context, index uint64) (TraceTree, bool) {
	tracer, ok := tracerFromContext(ctx)
	if !ok {
		return TraceTree{}, false
	}
	return tracer.TraceCauseTreeByIndex(index)
}

func TraceCauseTree(ctx context.Context, ordered OrderedMsg) (TraceTree, bool) {
	tracer, ok := tracerFromContext(ctx)
	if !ok {
		return TraceTree{}, false
	}
	return tracer.TraceCauseTree(ordered)
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
