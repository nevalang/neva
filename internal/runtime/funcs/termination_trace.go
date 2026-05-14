package funcs

import (
	"fmt"
	"os"
	"strings"

	"github.com/nevalang/neva/internal/runtime"
)

// formatTerminationDataflowTrace renders termination-oriented dataflow ancestry
// for the current message. Missing trace data is treated as invariant violation.
func formatTerminationDataflowTrace(title string, tracer *runtime.Tracer, msg runtime.OrderedMsg) string {
	tree, ok := tracer.TraceCauseTree(msg)
	if !ok {
		panic("runtime invariant: missing dataflow trace for termination message")
	}

	var builder strings.Builder
	receiver := "<?>"
	if tree.Hop.Receiver != nil {
		receiver = formatPanicPortSlotAddr(*tree.Hop.Receiver)
	}
	component := panicComponentName(tree.Hop.Receiver)

	builder.WriteString(title + "\n")
	builder.WriteString("direction: newest -> oldest (top -> bottom)\n")
	builder.WriteString("sink: " + receiver + "\n")
	if component != "" {
		builder.WriteString("component: " + component + "\n")
	}
	builder.WriteString("events:\n")
	formatPanicTraceTree(&builder, &tree, "  ")

	return strings.TrimRight(builder.String(), "\n")
}

func writeTerminationTrace(title string, tracer *runtime.Tracer, msg runtime.OrderedMsg) {
	trace := formatTerminationDataflowTrace(title, tracer, msg)
	if _, err := fmt.Fprintln(os.Stderr, trace); err != nil {
		panic(err)
	}
}
