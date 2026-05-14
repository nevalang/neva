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
