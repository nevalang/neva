package funcs

import (
	"context"
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/nevalang/neva/internal/runtime"
)

// --- Errors ---

func errFromErr(err error) runtime.StructMsg {
	return errFromString(err.Error())
}

func errFromString(s string) runtime.StructMsg {
	return runtime.NewStructMsg([]runtime.StructField{
		runtime.NewStructField("text", runtime.NewStringMsg(s)),
		runtime.NewStructField("child", runtime.NewUnionMsg("None", nil)),
	})
}

func streamItem(data runtime.Msg, idx int64, last bool) runtime.StructMsg {
	return runtime.NewStructMsg([]runtime.StructField{
		runtime.NewStructField("data", data),
		runtime.NewStructField("idx", runtime.NewIntMsg(idx)),
		runtime.NewStructField("last", runtime.NewBoolMsg(last)),
	})
}

// --- Structs ---

func emptyStruct() runtime.StructMsg {
	return runtime.NewStructMsg(nil)
}

// --- Receives ---

func receive2(
	ctx context.Context,
	firstIn runtime.SingleInport,
	secondIn runtime.SingleInport,
) (runtime.OrderedMsg, runtime.OrderedMsg, bool) {
	var firstMsg, secondMsg runtime.OrderedMsg
	var firstOK, secondOK bool

	var waitGroup sync.WaitGroup
	waitGroup.Go(func() {
		firstMsg, firstOK = firstIn.Receive(ctx)
	})
	waitGroup.Go(func() {
		secondMsg, secondOK = secondIn.Receive(ctx)
	})
	waitGroup.Wait()

	return firstMsg, secondMsg, firstOK && secondOK
}
func receive3(
	ctx context.Context,
	firstIn runtime.SingleInport,
	secondIn runtime.SingleInport,
	thirdIn runtime.SingleInport,
) (runtime.OrderedMsg, runtime.OrderedMsg, runtime.OrderedMsg, bool) {
	var firstMsg, secondMsg, thirdMsg runtime.OrderedMsg
	var firstOK, secondOK, thirdOK bool

	var waitGroup sync.WaitGroup
	waitGroup.Go(func() {
		firstMsg, firstOK = firstIn.Receive(ctx)
	})
	waitGroup.Go(func() {
		secondMsg, secondOK = secondIn.Receive(ctx)
	})
	waitGroup.Go(func() {
		thirdMsg, thirdOK = thirdIn.Receive(ctx)
	})
	waitGroup.Wait()

	return firstMsg, secondMsg, thirdMsg, firstOK && secondOK && thirdOK
}

func receive4(
	ctx context.Context,
	firstIn runtime.SingleInport,
	secondIn runtime.SingleInport,
	thirdIn runtime.SingleInport,
	fourthIn runtime.SingleInport,
) (runtime.OrderedMsg, runtime.OrderedMsg, runtime.OrderedMsg, runtime.OrderedMsg, bool) {
	var firstMsg, secondMsg, thirdMsg, fourthMsg runtime.OrderedMsg
	var firstOK, secondOK, thirdOK, fourthOK bool

	var waitGroup sync.WaitGroup
	waitGroup.Go(func() {
		firstMsg, firstOK = firstIn.Receive(ctx)
	})
	waitGroup.Go(func() {
		secondMsg, secondOK = secondIn.Receive(ctx)
	})
	waitGroup.Go(func() {
		thirdMsg, thirdOK = thirdIn.Receive(ctx)
	})
	waitGroup.Go(func() {
		fourthMsg, fourthOK = fourthIn.Receive(ctx)
	})
	waitGroup.Wait()

	return firstMsg, secondMsg, thirdMsg, fourthMsg, firstOK && secondOK && thirdOK && fourthOK
}

// --- Trace ---

// formatTerminationDataflowTrace renders termination-oriented dataflow ancestry
// for the current message. Missing trace data is treated as invariant violation.
func formatTerminationDataflowTrace(title string, tracer *runtime.Tracer, msg runtime.OrderedMsg) string {
	tree, ok := traceFromOrderedMsg(tracer, msg)
	if !ok {
		panic("runtime invariant: missing dataflow trace for termination message")
	}

	var builder strings.Builder
	receiver := "<?>"
	if tree.Hop.Receiver != nil {
		receiver = formatTracePortSlotAddr(*tree.Hop.Receiver)
	}
	component := traceComponentName(tree.Hop.Receiver)

	builder.WriteString(title + "\n")
	builder.WriteString("direction: newest -> oldest (top -> bottom)\n")
	builder.WriteString("sink: " + receiver + "\n")
	if component != "" {
		builder.WriteString("component: " + component + "\n")
	}
	builder.WriteString("events:\n")
	formatTraceTree(&builder, &tree, "  ")

	return strings.TrimRight(builder.String(), "\n")
}

func writeTerminationTrace(title string, io runtime.IO, msg runtime.OrderedMsg) {
	tracer := runtime.TracerFromIO(io)
	trace := formatTerminationDataflowTrace(title, tracer, msg)
	if _, err := fmt.Fprintln(os.Stderr, trace); err != nil {
		panic(err)
	}
}

// traceTree is a derived, read-only projection rebuilt from traceStore hop links.
// It is intentionally denormalized for traversal/formatting APIs.
type traceTree struct {
	Parents []traceTree
	Hop     runtime.TraceHop
}

func formatTraceTree(builder *strings.Builder, tree *traceTree, indent string) {
	builder.WriteString(indent + "- " + formatTraceHopFlow(tree.Hop) + "\n")
	for _, parent := range tree.Parents {
		formatTraceTree(builder, &parent, indent+"  ")
	}
}

func formatTraceHopFlow(hop runtime.TraceHop) string {
	recv := "<?>"
	send := "<?>"
	if hop.Receiver != nil {
		recv = formatTracePortSlotAddr(*hop.Receiver)
	}
	if hop.Sender != nil {
		send = formatTracePortSlotAddr(*hop.Sender)
	}
	return fmt.Sprintf("%s -> %s", send, recv)
}

func traceComponentName(receiver *runtime.PortSlotAddr) string {
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

func formatTracePortSlotAddr(slot runtime.PortSlotAddr) string {
	slot.Path = normalizeTracePortPath(slot.Path)
	s := fmt.Sprintf("%s:%s", slot.Path, slot.Port)
	if slot.Index != nil {
		s = fmt.Sprintf("%s[%d]", s, *slot.Index)
	}
	return s
}

func normalizeTracePortPath(path string) string {
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

func traceFromOrderedMsg(tracer *runtime.Tracer, ordered runtime.OrderedMsg) (traceTree, bool) {
	rootHop, ok := tracer.HopByOrderedMsg(ordered)
	if !ok {
		return traceTree{}, false
	}
	return traceTreeFromHop(tracer, rootHop, map[uint64]struct{}{})
}

// traceTreeFromHop is just a recursive helper for traceFromOrderedMsg.
func traceTreeFromHop(tracer *runtime.Tracer, hop runtime.TraceHop, visited map[uint64]struct{}) (traceTree, bool) {
	if hop.Index == 0 {
		return traceTree{}, false
	}
	if _, seen := visited[hop.Index]; seen {
		return traceTree{}, false
	}
	visited[hop.Index] = struct{}{}

	// Rebuild tree view from normalized hop links; no second persisted source of truth.
	tree := traceTree{
		Hop:     hop,
		Parents: make([]traceTree, 0, len(hop.CauseIndexes)),
	}
	for _, parentHop := range tracer.HopsByCauseIndexes(hop.CauseIndexes) {
		parentTree, ok := traceTreeFromHop(tracer, parentHop, visited)
		if !ok {
			continue
		}
		tree.Parents = append(tree.Parents, parentTree)
	}

	delete(visited, hop.Index)
	return tree, true
}
