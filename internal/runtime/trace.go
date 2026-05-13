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
	pendingMap map[string]*pendingCausesState
	store      traceStore
	pendingMu  sync.Mutex
}

type traceStore struct {
	hops map[uint64]TraceHop
	mu   sync.RWMutex
}

type tracerKey struct{}

func NewTracer() *Tracer {
	return &Tracer{
		store: traceStore{
			hops: make(map[uint64]TraceHop),
		},
		pendingMap: map[string]*pendingCausesState{},
	}
}

func contextWithTracer(ctx context.Context, tracer *Tracer) context.Context {
	return context.WithValue(ctx, tracerKey{}, tracer)
}

func tracerFromContext(ctx context.Context) (*Tracer, bool) {
	tracer, ok := ctx.Value(tracerKey{}).(*Tracer)
	if !ok || tracer == nil {
		return nil, false
	}
	return tracer, true
}

func WithTracer(ctx context.Context) context.Context {
	return contextWithTracer(ctx, NewTracer())
}

func (t *Tracer) causeIndexesFromOrdered(causes []OrderedMsg) []uint64 {
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

type pendingCausesState struct {
	causes  map[uint64]struct{}
	emitted bool
}

func nodePathKey(path string) string {
	return strings.TrimSuffix(strings.TrimSuffix(path, "/in"), "/out")
}

func (t *Tracer) onReceived(path string, ordered OrderedMsg) {
	if ordered.index == 0 {
		return
	}
	key := nodePathKey(path)
	t.pendingMu.Lock()
	defer t.pendingMu.Unlock()
	state, ok := t.pendingMap[key]
	if !ok {
		state = &pendingCausesState{causes: map[uint64]struct{}{}}
		t.pendingMap[key] = state
	}
	if state.emitted {
		clear(state.causes)
		state.emitted = false
	}
	state.causes[ordered.index] = struct{}{}
}

func (t *Tracer) currentCauseIndexes(path string) []uint64 {
	key := nodePathKey(path)
	t.pendingMu.Lock()
	defer t.pendingMu.Unlock()
	state, ok := t.pendingMap[key]
	if !ok || len(state.causes) == 0 {
		return nil
	}
	indexes := make([]uint64, 0, len(state.causes))
	for index := range state.causes {
		indexes = append(indexes, index)
	}
	slices.Sort(indexes)
	state.emitted = true
	return indexes
}

func (t *Tracer) causeIndexesForSend(path string, causes []OrderedMsg) []uint64 {
	if len(causes) != 0 {
		return t.causeIndexesFromOrdered(causes)
	}
	return t.currentCauseIndexes(path)
}

func (t *Tracer) RecordSent(
	sender PortSlotAddr,
	ordered OrderedMsg,
	causes []OrderedMsg,
) TraceHop {
	t.store.mu.Lock()
	defer t.store.mu.Unlock()

	hop := t.store.hops[ordered.index]
	hop.Index = ordered.index
	hop.CauseIndexes = t.causeIndexesForSend(sender.Path, causes)
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
	t.onReceived(receiver.Path, ordered)
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
