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
		msg, ok := inport.Receive(ctx)
		if !ok {
			return
		}
		ReportProgramPanic(ctx, msg)
		if cancel, ok := CancelFuncFromContext(ctx); ok {
			cancel()
		}
	}, nil
}

func TestCall_ReturnsProgramPanicError(t *testing.T) {
	resetRuntimeTraceStateForTests()

	startToPanic := make(chan OrderedMsg, 1)
	stopChan := make(chan OrderedMsg)

	prog := Program{
		Start: NewSingleOutport(
			PortAddr{Path: "prog", Port: "start"},
			ProdInterceptor{},
			startToPanic,
		),
		Stop: NewSingleInport(
			stopChan,
			PortAddr{Path: "prog", Port: "stop"},
			ProdInterceptor{},
		),
		FuncCalls: []FuncCall{
			{
				Ref: "panic",
				IO: IO{
					In: NewInports(map[string]Inport{
						"data": NewInport(
							nil,
							NewSingleInport(
								startToPanic,
								PortAddr{Path: "panic/in", Port: "data"},
								ProdInterceptor{},
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
	if err == nil {
		t.Fatalf("expected program panic error")
	}
	if !IsProgramPanic(err) {
		t.Fatalf("expected program panic error type, got %T (%v)", err, err)
	}
}
