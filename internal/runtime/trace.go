package runtime

import (
	"fmt"
	"slices"
	"sync"
)

type TraceHop struct {
	Sender   *PortSlotAddr
	Receiver *PortSlotAddr
	Message  string
	// CauseIndexes is the single source of truth for causal edges in traceStore.
	// Tree views are reconstructed from these indexes on demand.
	CauseIndexes []uint64
	// Index is materialized on read from traceStore map key.
	// We do not persist it in storage separately to avoid duplicated state.
	Index uint64
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
	// NOTE: this lock can become a bottleneck under high message throughput
	// because every send/receive hop touches shared storage. Keep in mind for
	// future sharding / lock-reduction work.
	mu sync.RWMutex
}

func NewTracer() *Tracer {
	return &Tracer{
		store: traceStore{
			hops: make(map[uint64]TraceHop),
		},
	}
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

func (t *Tracer) recordSent(
	sender PortSlotAddr,
	ordered OrderedMsg,
	causes []OrderedMsg,
) TraceHop {
	t.store.mu.Lock()
	defer t.store.mu.Unlock()

	hop := t.store.hops[ordered.index]
	// Runtime funcs are responsible for passing explicit OrderedMsg causes.
	// Tracer does not infer fallback causes from receive-side state.
	hop.CauseIndexes = t.causeIndexesFromOrdered(causes)
	senderCopy := sender
	hop.Sender = &senderCopy
	hop.Message = fmt.Sprint(ordered.Msg)
	t.store.hops[ordered.index] = hop

	return hop
}

func (t *Tracer) recordReceived(receiver PortSlotAddr, ordered OrderedMsg) {
	t.store.mu.Lock()
	defer t.store.mu.Unlock()

	hop := t.store.hops[ordered.index]
	receiverCopy := receiver
	hop.Receiver = &receiverCopy
	t.store.hops[ordered.index] = hop
}

func (t *Tracer) traceHopByIndex(index uint64) (TraceHop, bool) {
	t.store.mu.RLock()
	defer t.store.mu.RUnlock()

	hop, ok := t.store.hops[index]
	hop.Index = index
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
