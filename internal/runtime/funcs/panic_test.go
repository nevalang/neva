package funcs

import (
	"strings"
	"testing"

	"github.com/nevalang/neva/internal/runtime"
)

func TestFormatTraceHopFlow_NormalizesSendToReceive(t *testing.T) {
	hop := runtime.TraceHop{
		Sender:   &runtime.PortSlotAddr{PortAddr: runtime.PortAddr{Path: "http/in", Port: "req"}},
		Receiver: &runtime.PortSlotAddr{PortAddr: runtime.PortAddr{Path: "parse/in", Port: "data"}},
	}

	stats := collectTraceRenderStats(&traceTree{Hop: hop})
	got := formatTraceHopFlow(hop, stats, false)
	if got != "parse:data <- http" {
		t.Fatalf("unexpected hop flow: %s", got)
	}
}

func TestFormatTraceTree_FanInRendersAllParents(t *testing.T) {
	tree := traceTree{
		Hop: runtime.TraceHop{
			Sender:   &runtime.PortSlotAddr{PortAddr: runtime.PortAddr{Path: "fanin/out", Port: "res"}},
			Receiver: &runtime.PortSlotAddr{PortAddr: runtime.PortAddr{Path: "prog/out", Port: "stop"}},
		},
		Parents: []traceTree{
			{
				Hop: runtime.TraceHop{
					Sender:   &runtime.PortSlotAddr{PortAddr: runtime.PortAddr{Path: "first/out", Port: "res"}},
					Receiver: &runtime.PortSlotAddr{PortAddr: runtime.PortAddr{Path: "fanin/in", Port: "first"}},
				},
			},
			{
				Hop: runtime.TraceHop{
					Sender:   &runtime.PortSlotAddr{PortAddr: runtime.PortAddr{Path: "second/out", Port: "res"}},
					Receiver: &runtime.PortSlotAddr{PortAddr: runtime.PortAddr{Path: "fanin/in", Port: "second"}},
				},
			},
			{
				Hop: runtime.TraceHop{
					Sender:   &runtime.PortSlotAddr{PortAddr: runtime.PortAddr{Path: "third/out", Port: "res"}},
					Receiver: &runtime.PortSlotAddr{PortAddr: runtime.PortAddr{Path: "fanin/in", Port: "third"}},
				},
			},
		},
	}

	formatted := buildFormattedPanicTraceForTest(&tree)
	if !strings.Contains(formatted, "sink: prog:stop") {
		t.Fatalf("expected normalized sink, got:\n%s", formatted)
	}
	if !strings.Contains(formatted, "component: prog") {
		t.Fatalf("expected component, got:\n%s", formatted)
	}
	if !strings.Contains(formatted, "prog:stop <- fanin") {
		t.Fatalf("expected fan-in output hop, got:\n%s", formatted)
	}
	if strings.Count(formatted, "fanin:first <- first")+strings.Count(formatted, "fanin:second <- second")+strings.Count(formatted, "fanin:third <- third") < 3 {
		t.Fatalf("expected all fan-in parents, got:\n%s", formatted)
	}
}

func buildFormattedPanicTraceForTest(tree *traceTree) string {
	var builder strings.Builder
	panicReceiver := formatTracePortSlotAddr(*tree.Hop.Receiver)
	panicComponent := traceComponentName(tree.Hop.Receiver)
	builder.WriteString("panic cause dataflow trace\n")
	builder.WriteString("direction: newest -> oldest (top -> bottom)\n")
	builder.WriteString("sink: " + panicReceiver + "\n")
	if panicComponent != "" {
		builder.WriteString("component: " + panicComponent + "\n")
	}
	stats := collectTraceRenderStats(tree)
	builder.WriteString(formatTraceHopFlow(tree.Hop, stats, true) + "\n")
	for i := range tree.Parents {
		formatTraceTree(&builder, &tree.Parents[i], "", i == len(tree.Parents)-1, stats)
	}
	return strings.TrimRight(builder.String(), "\n")
}
