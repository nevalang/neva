package funcs

import (
	"context"
	"testing"

	"github.com/nevalang/neva/internal/runtime"
)

func TestCall_ReturnsCustomExitCodeFromOSExit(t *testing.T) {
	startToExit := make(chan runtime.OrderedMsg, 1)
	stopChan := make(chan runtime.OrderedMsg)
	tracer := runtime.NewTracer()
	interceptor := runtime.NoEffectInterceptor{}

	prog := runtime.Program{
		Start: runtime.NewSingleOutport(
			tracer,
			runtime.PortAddr{Path: "prog", Port: "start"},
			interceptor,
			startToExit,
		),
		Stop: runtime.NewSingleInport(
			tracer,
			stopChan,
			runtime.PortAddr{Path: "prog", Port: "stop"},
			interceptor,
		),
		FuncCalls: []runtime.FuncCall{
			{
				Ref: "os_exit",
				IO: runtime.IO{
					In: runtime.NewInports(map[string]runtime.Inport{
						"code": runtime.NewInport(
							nil,
							runtime.NewSingleInport(
								tracer,
								startToExit,
								runtime.PortAddr{Path: "os/exit/in", Port: "code"},
								interceptor,
							),
						),
					}),
					Out: runtime.NewOutports(map[string]runtime.Outport{}),
				},
			},
		},
	}

	_, exitCode, err := runtime.Call(context.Background(), prog, NewRegistry(), runtime.NewIntMsg(7))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if exitCode != 7 {
		t.Fatalf("expected exit code 7, got %d", exitCode)
	}
}
