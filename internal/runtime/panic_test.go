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
		FailProgram(ctx)
	}, nil
}

func TestCall_ReturnsProgramPanicError(t *testing.T) {
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

	_, err := Call(context.Background(), prog, registry, NewStringMsg("boom"))
	if !IsProgramPanicError(err) {
		t.Fatalf("expected ProgramPanicError, got %T (%v)", err, err)
	}
}
