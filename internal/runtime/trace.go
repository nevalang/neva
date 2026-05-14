package runtime

import (
	"fmt"
	"slices"
	"sync"
)

type Tracer struct {
	hops map[uint64]TraceHop
	// NOTE: this lock can become a bottleneck under high message throughput
	// because every send/receive hop touches shared storage. Keep in mind for
	// future sharding / lock-reduction work.
	mu sync.RWMutex
}

// Runtime funcs are responsible for passing explicit OrderedMsg causes.
// Tracer does not infer fallback causes from receive-side state.
func (t *Tracer) recordSent(
	sender PortSlotAddr,
	ordered OrderedMsg,
	causes []OrderedMsg,
) TraceHop {
	t.mu.Lock()
	defer t.mu.Unlock()

	hop := t.hops[ordered.index]
	hop.CauseIndexes = t.causeIndexesFromOrdered(causes)
	senderCopy := sender
	hop.Sender = &senderCopy
	hop.Message = fmt.Sprint(ordered.Msg)
	t.hops[ordered.index] = hop

	return hop
}

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

func (t *Tracer) recordReceived(receiver PortSlotAddr, ordered OrderedMsg) {
	t.mu.Lock()
	defer t.mu.Unlock()

	hop := t.hops[ordered.index]
	receiverCopy := receiver
	hop.Receiver = &receiverCopy
	t.hops[ordered.index] = hop
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

func (t *Tracer) hopByIndex(index uint64) (TraceHop, bool) {
	t.mu.RLock()
	defer t.mu.RUnlock()
	hop, ok := t.hops[index]
	hop.Index = index
	return hop, ok
}

// HopByOrderedMsg returns a normalized hop for a concrete ordered message.
func (t *Tracer) HopByOrderedMsg(ordered OrderedMsg) (TraceHop, bool) {
	if ordered.index == 0 {
		return TraceHop{}, false
	}
	return t.hopByIndex(ordered.index)
}

// HopsByCauseIndexes resolves parent hops by stored cause indexes.
// Missing indexes are ignored to keep trace-read path resilient.
func (t *Tracer) HopsByCauseIndexes(causeIndexes []uint64) []TraceHop {
	if len(causeIndexes) == 0 {
		return nil
	}
	parents := make([]TraceHop, 0, len(causeIndexes))
	for _, causeIndex := range causeIndexes {
		hop, ok := t.hopByIndex(causeIndex)
		if !ok {
			continue
		}
		parents = append(parents, hop)
	}
	return parents
}

func NewTracer() *Tracer {
	return &Tracer{
		hops: make(map[uint64]TraceHop),
	}
}
