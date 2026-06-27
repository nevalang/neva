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

// --- Message utils ---

// Runtime style note: keep OrderedMsg access explicit (selected.OrderedMsg.Msg),
// do not rely on promoted fields from embedded structs in hot paths.
//
//nolint:ireturn // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
func tryToUnboxIfUnion(msg runtime.Msg) runtime.Msg {
	unionMsg, ok := runtime.AsUnion(msg)
	if !ok {
		return msg
	}

	if unionMsg.Data() == nil {
		return msg
	}

	return unionMsg.Data()
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
	stats := collectTraceRenderStats(&tree)

	builder.WriteString(title + "\n")
	builder.WriteString("direction: newest <- oldest (top -> bottom)\n")
	builder.WriteString("sink: " + receiver + "\n")
	if component != "" {
		builder.WriteString("component: " + component + "\n")
	}
	builder.WriteString(formatTraceHopFlow(tree.Hop, stats, true) + "\n")
	for i := range tree.Parents {
		formatTraceTree(&builder, &tree.Parents[i], "", i == len(tree.Parents)-1, stats)
	}

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

type traceRenderStats struct {
	senderPortsByPath   map[string]map[string]struct{}
	receiverPortsByPath map[string]map[string]struct{}
}

func collectTraceRenderStats(tree *traceTree) traceRenderStats {
	stats := traceRenderStats{
		senderPortsByPath:   map[string]map[string]struct{}{},
		receiverPortsByPath: map[string]map[string]struct{}{},
	}

	var visit func(*traceTree)
	visit = func(node *traceTree) {
		if node.Hop.Sender != nil {
			addPort(stats.senderPortsByPath, normalizeTracePortPath(node.Hop.Sender.Path), node.Hop.Sender.Port)
		}
		if node.Hop.Receiver != nil {
			addPort(stats.receiverPortsByPath, normalizeTracePortPath(node.Hop.Receiver.Path), node.Hop.Receiver.Port)
		}
		for i := range node.Parents {
			visit(&node.Parents[i])
		}
	}
	visit(tree)

	return stats
}

func addPort(portsByPath map[string]map[string]struct{}, path, port string) {
	ports, ok := portsByPath[path]
	if !ok {
		ports = map[string]struct{}{}
		portsByPath[path] = ports
	}
	ports[port] = struct{}{}
}

func formatTraceTree(
	builder *strings.Builder,
	tree *traceTree,
	prefix string,
	isLast bool,
	stats traceRenderStats,
) {
	connector := "├─ "
	nextPrefix := prefix + "│  "
	if isLast {
		connector = "└─ "
		nextPrefix = prefix + "   "
	}
	builder.WriteString(prefix + connector + formatTraceHopFlow(tree.Hop, stats, false) + "\n")
	for i := range tree.Parents {
		formatTraceTree(builder, &tree.Parents[i], nextPrefix, i == len(tree.Parents)-1, stats)
	}
}

func formatTraceHopFlow(hop runtime.TraceHop, stats traceRenderStats, forceReceiverPort bool) string {
	recv := "<?>"
	send := "<?>"
	if hop.Receiver != nil {
		recv = formatTraceEndpoint(*hop.Receiver, false, stats, forceReceiverPort)
	}
	if hop.Sender != nil {
		send = formatTraceEndpoint(*hop.Sender, true, stats, false)
	}
	return fmt.Sprintf("%s <- %s", recv, send)
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

func formatTraceEndpoint(
	slot runtime.PortSlotAddr,
	isSender bool,
	stats traceRenderStats,
	forcePort bool,
) string {
	path := normalizeTracePortPath(slot.Path)
	if path == "" && slot.Port == "start" {
		return ":start"
	}

	includePort := forcePort || slot.Index != nil || shouldIncludePort(path, slot.Port, isSender, stats)
	formatted := path
	if includePort {
		formatted = fmt.Sprintf("%s:%s", path, slot.Port)
	}
	if slot.Index != nil {
		formatted = fmt.Sprintf("%s[%d]", formatted, *slot.Index)
	}
	return formatted
}

func shouldIncludePort(path, port string, isSender bool, stats traceRenderStats) bool {
	// Keep sink and panic port explicit for readability in termination traces.
	if path == "panic" || port == "data" {
		return true
	}
	portsByPath := stats.receiverPortsByPath
	if isSender {
		portsByPath = stats.senderPortsByPath
	}
	ports := portsByPath[path]
	return len(ports) > 1
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

const (
	streamTagOpen  = "Open"
	streamTagData  = "Data"
	streamTagClose = "Close"
)

func streamOpen() runtime.UnionMsg {
	return runtime.NewUnionMsg(streamTagOpen, nil)
}

func streamData(data runtime.Msg) runtime.UnionMsg {
	return runtime.NewUnionMsg(streamTagData, data)
}

func streamClose() runtime.UnionMsg {
	return runtime.NewUnionMsg(streamTagClose, nil)
}

func streamUnion(msg runtime.Msg) runtime.UnionMsg {
	defer func() {
		if recover() != nil {
			panic(fmt.Sprintf("runtime: expected stream union message, got %T", msg))
		}
	}()

	return msg.Union()
}

func isStreamOpen(msg runtime.Msg) bool {
	return streamUnion(msg).Tag() == streamTagOpen
}

func isStreamData(msg runtime.Msg) bool {
	return streamUnion(msg).Tag() == streamTagData
}

func isStreamClose(msg runtime.Msg) bool {
	return streamUnion(msg).Tag() == streamTagClose
}

//nolint:ireturn // Stream payloads are runtime.Msg values by contract.
func streamDataValue(msg runtime.Msg) runtime.Msg {
	u := streamUnion(msg)
	if u.Tag() != streamTagData {
		panic("runtime: expected stream Data message")
	}
	return u.Data()
}
