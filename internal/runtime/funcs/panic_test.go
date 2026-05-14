package funcs

import (
	"strings"
	"testing"

	"github.com/nevalang/neva/internal/runtime"
)

func TestFormatPanicHopFlow_NormalizesSendToReceive(t *testing.T) {
	hop := runtime.TraceHop{
		Sender:   &runtime.PortSlotAddr{PortAddr: runtime.PortAddr{Path: "http/in", Port: "req"}},
		Receiver: &runtime.PortSlotAddr{PortAddr: runtime.PortAddr{Path: "parse/in", Port: "data"}},
	}

	got := formatPanicHopFlow(hop)
	if got != "http:req -> parse:data" {
		t.Fatalf("unexpected hop flow: %s", got)
	}
}

func TestFormatPanicTraceTree_FanInRendersAllParents(t *testing.T) {
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
	if !strings.Contains(formatted, "fanin:res -> prog:stop") {
		t.Fatalf("expected fan-in output hop, got:\n%s", formatted)
	}
	if strings.Count(formatted, "first:res -> fanin:first")+strings.Count(formatted, "second:res -> fanin:second")+strings.Count(formatted, "third:res -> fanin:third") < 3 {
		t.Fatalf("expected all fan-in parents, got:\n%s", formatted)
	}
}

func buildFormattedPanicTraceForTest(tree *traceTree) string {
	var builder strings.Builder
	panicReceiver := formatPanicPortSlotAddr(*tree.Hop.Receiver)
	panicComponent := panicComponentName(tree.Hop.Receiver)
	builder.WriteString("panic cause dataflow trace\n")
	builder.WriteString("direction: newest -> oldest (top -> bottom)\n")
	builder.WriteString("sink: " + panicReceiver + "\n")
	if panicComponent != "" {
		builder.WriteString("component: " + panicComponent + "\n")
	}
	builder.WriteString("events:\n")
	formatPanicTraceTree(&builder, tree, "  ")
	return strings.TrimRight(builder.String(), "\n")
}
