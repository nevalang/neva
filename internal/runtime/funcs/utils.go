package funcs

import (
	"context"
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/nevalang/neva/internal/runtime"
	"github.com/nevalang/neva/internal/runtime/messages"
)

// --- Errors ---

func errFromErr(err error) messages.StructMsg {
	return errFromString(err.Error())
}

func errFromString(s string) messages.StructMsg {
	return messages.NewStructMsg([]messages.StructField{
		messages.NewStructField("text", messages.NewStringMsg(s)),
		messages.NewStructField("child", messages.NewUnionMsg("None", nil)),
	})
}

// --- Structs ---

func emptyStruct() messages.StructMsg {
	return messages.NewStructMsg(nil)
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
// tryToUnboxIfUnion returns union payload if message is a data-carrying union;
// otherwise it returns the original message unchanged.
//
//nolint:ireturn // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
func tryToUnboxIfUnion(msg messages.Msg) messages.Msg {
	unionMsg, ok := messages.AsUnion(msg)
	if !ok {
		return msg
	}

	if unionMsg.Data() == nil {
		return msg
	}

	return unionMsg.Data()
}

// listToMsgs converts any supported typed list view to untyped []messages.Msg.
// Typed scalar paths avoid panicking Untyped() calls on typed list implementations.
func listToMsgs(list messages.ListMsg) []messages.Msg {
	if values, ok := messages.AsListInts(list); ok {
		msgs := make([]messages.Msg, len(values))
		for i := range values {
			msgs[i] = messages.NewIntMsg(values[i])
		}
		return msgs
	}
	if values, ok := messages.AsListStrings(list); ok {
		msgs := make([]messages.Msg, len(values))
		for i := range values {
			msgs[i] = messages.NewStringMsg(values[i])
		}
		return msgs
	}
	if values, ok := messages.AsListBools(list); ok {
		msgs := make([]messages.Msg, len(values))
		for i := range values {
			msgs[i] = messages.NewBoolMsg(values[i])
		}
		return msgs
	}
	if values, ok := messages.AsListFloats(list); ok {
		msgs := make([]messages.Msg, len(values))
		for i := range values {
			msgs[i] = messages.NewFloatMsg(values[i])
		}
		return msgs
	}
	return list.Untyped()
}

// dictToMsgs converts any supported typed dict view to untyped map[string]messages.Msg.
// Typed scalar paths avoid panicking Untyped() calls on typed dict implementations.
func dictToMsgs(dict messages.DictMsg) map[string]messages.Msg {
	if values, ok := messages.AsDictInts(dict); ok {
		msgs := make(map[string]messages.Msg, len(values))
		for key, value := range values {
			msgs[key] = messages.NewIntMsg(value)
		}
		return msgs
	}
	if values, ok := messages.AsDictStrings(dict); ok {
		msgs := make(map[string]messages.Msg, len(values))
		for key, value := range values {
			msgs[key] = messages.NewStringMsg(value)
		}
		return msgs
	}
	if values, ok := messages.AsDictBools(dict); ok {
		msgs := make(map[string]messages.Msg, len(values))
		for key, value := range values {
			msgs[key] = messages.NewBoolMsg(value)
		}
		return msgs
	}
	if values, ok := messages.AsDictFloats(dict); ok {
		msgs := make(map[string]messages.Msg, len(values))
		for key, value := range values {
			msgs[key] = messages.NewFloatMsg(value)
		}
		return msgs
	}
	return dict.Untyped()
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
