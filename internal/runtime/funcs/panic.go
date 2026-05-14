package funcs

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/nevalang/neva/internal/runtime"
)

type panicker struct{}

// formatPanicDataflowTrace renders panic-oriented dataflow ancestry for the current message.
// Missing trace data is treated as runtime invariant violation and causes Go panic.
func formatPanicDataflowTrace(tracer *runtime.Tracer, msg runtime.OrderedMsg) string {
	tree, ok := tracer.TraceCauseTree(msg)
	if !ok {
		panic("runtime invariant: missing dataflow trace for panic message")
	}

	var builder strings.Builder
	panicReceiver := "<?>"
	if tree.Hop.Receiver != nil {
		panicReceiver = formatPanicPortSlotAddr(*tree.Hop.Receiver)
	}
	panicComponent := panicComponentName(tree.Hop.Receiver)

	builder.WriteString("panic cause dataflow trace\n")
	builder.WriteString("direction: newest -> oldest (top -> bottom)\n")
	builder.WriteString("panic sink: " + panicReceiver + "\n")
	if panicComponent != "" {
		builder.WriteString("panic component: " + panicComponent + "\n")
	}
	builder.WriteString("events:\n")
	formatPanicTraceTree(&builder, &tree, "  ")

	return strings.TrimRight(builder.String(), "\n")
}

func formatPanicTraceTree(builder *strings.Builder, tree *runtime.TraceTree, indent string) {
	builder.WriteString(indent + "- " + formatPanicHopFlow(tree.Hop) + "\n")
	for _, parent := range tree.Parents {
		formatPanicTraceTree(builder, &parent, indent+"  ")
	}
}

func formatPanicHopFlow(hop runtime.TraceHop) string {
	recv := "<?>"
	send := "<?>"
	if hop.Receiver != nil {
		recv = formatPanicPortSlotAddr(*hop.Receiver)
	}
	if hop.Sender != nil {
		send = formatPanicPortSlotAddr(*hop.Sender)
	}
	return fmt.Sprintf("%s -> %s", send, recv)
}

func panicComponentName(receiver *runtime.PortSlotAddr) string {
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

func formatPanicPortSlotAddr(slot runtime.PortSlotAddr) string {
	slot.Path = normalizePanicPortPath(slot.Path)
	s := fmt.Sprintf("%s:%s", slot.Path, slot.Port)
	if slot.Index != nil {
		s = fmt.Sprintf("%s[%d]", s, *slot.Index)
	}
	return s
}

func normalizePanicPortPath(path string) string {
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

func (p panicker) Create(
	runtimeIO runtime.IO,
	_ runtime.Msg,
) (func(ctx context.Context), error) {
	msgIn, err := runtimeIO.In.Single("data")
	if err != nil {
		//nolint:wrapcheck // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
		return nil, err
	}

	return func(ctx context.Context) {
		panicMsg, ok := msgIn.Receive(ctx)
		if !ok {
			return
		}

		if _, err := fmt.Fprintln(os.Stderr, "panic:", panicMsg); err != nil {
			panic(err)
		}

		trace := formatPanicDataflowTrace(runtimeIO.Tracer(), panicMsg)
		if _, err := fmt.Fprintln(os.Stderr, trace); err != nil {
			panic(err)
		}

		runtime.TerminateProgram(ctx, 1)
	}, nil
}
