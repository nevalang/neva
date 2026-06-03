package runtime

import (
	"context"
	"testing"
)

type testPanicCreator struct{}

func (testPanicCreator) Create(io IO, _ Msg) (func(context.Context), error) {
	inport, err := io.In.Single("data")
	if err != nil {
		return nil, err
	}

	return func(ctx context.Context) {
		_, ok := inport.Receive(ctx)
		if !ok {
			return
		}
		Terminate(ctx, 1)
	}, nil
}

func TestCall_ReturnsExitCodeOnProgramPanic(t *testing.T) {
	resetRuntimeTraceStateForTests()

	startToPanic := make(chan OrderedMsg, 1)
	stopChan := make(chan OrderedMsg)
	tracer := NewTracer()
	noEffectInterceptor := NoEffectInterceptor{}

	prog := Program{
		Start: NewSingleOutport(
			tracer,
			PortAddr{Path: "prog", Port: "start"},
			noEffectInterceptor,
			startToPanic,
		),
		Stop: NewSingleInport(
			tracer,
			stopChan,
			PortAddr{Path: "prog", Port: "stop"},
			noEffectInterceptor,
		),
		FuncCalls: []FuncCall{
			{
				Ref: "panic",
				IO: IO{
					In: NewInports(map[string]Inport{
						"data": NewInport(
							nil,
							NewSingleInport(
								tracer,
								startToPanic,
								PortAddr{Path: "panic/in", Port: "data"},
								noEffectInterceptor,
							),
						),
					}),
					Out: NewOutports(map[string]Outport{}),
				},
			},
		},
	}

	registry := map[string]FuncCreator{
		"panic": testPanicCreator{},
	}

	_, exitCode, err := Call(context.Background(), prog, registry, NewStringMsg("boom"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if exitCode != 1 {
		t.Fatalf("expected exit code 1, got %d", exitCode)
	}
}
