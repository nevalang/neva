package main

import (
	"context"
	"fmt"

	"github.com/emil14/neva/internal/runtime"
	"github.com/emil14/neva/internal/runtime/std/flow"
	"github.com/emil14/neva/internal/runtime/std/io"
)

func main() {
	printerRef := runtime.ComponentRef{
		Pkg:  "io",
		Name: "printer",
	}
	voidRef := runtime.ComponentRef{
		Pkg:  "io",
		Name: "void",
	}

	repo := map[runtime.ComponentRef]runtime.ComponentFunc{
		printerRef: io.Print,
		voidRef:    flow.Void,
	}
	componentRunner := runtime.NewComponentRunner(repo)

	giverRunner := runtime.GiverRunnerImlp{}

	routineRunner := runtime.NewRoutineRunner(giverRunner, componentRunner)

	interceptor := runtime.InterceptorImlp{}
	connector := runtime.NewConnector(interceptor)

	r := runtime.NewRuntime(connector, routineRunner)

	startPort := make(chan runtime.Msg)
	startPortAddr := runtime.PortAddr{
		Path: "root",
		Name: "sig",
	}

	printerInPort := make(chan runtime.Msg)
	printerInPortAddr := runtime.PortAddr{
		Path: "printer.in",
		Name: "v",
	}

	printerOutPort := make(chan runtime.Msg)
	printerOutPortAddr := runtime.PortAddr{
		Path: "printer.out",
		Name: "v",
	}

	voidInPort := make(chan runtime.Msg)
	voidInPortAddr := runtime.PortAddr{
		Path: "void.in",
		Name: "v",
	}

	prog := runtime.Program{
		StartPortAddr: startPortAddr,
		Ports: map[runtime.PortAddr]chan runtime.Msg{
			startPortAddr:      printerInPort,
			printerInPortAddr:  startPort,
			printerInPortAddr:  printerInPort,
			printerOutPortAddr: printerOutPort,
		},
		Connections: []runtime.Connection{
			{
				Sender: runtime.ConnectionSide{
					Port: startPort,
					Meta: runtime.ConnectionSideMeta{
						PortAddr: startPortAddr,
					},
				},
				Receivers: []runtime.ConnectionSide{
					{
						Port: printerInPort,
						Meta: runtime.ConnectionSideMeta{
							PortAddr: printerInPortAddr,
						},
					},
				},
			},
			{
				Sender: runtime.ConnectionSide{
					Port: printerOutPort,
					Meta: runtime.ConnectionSideMeta{
						PortAddr: printerOutPortAddr,
					},
				},
				Receivers: []runtime.ConnectionSide{
					{
						Port: voidInPort,
						Meta: runtime.ConnectionSideMeta{
							PortAddr: voidInPortAddr,
						},
					},
				},
			},
		},
		Routines: runtime.Routines{
			Component: []runtime.ComponentRoutine{
				{
					Ref: printerRef,
					IO: runtime.IO{
						In: map[string][]chan runtime.Msg{
							"v": {printerInPort},
						},
						Out: map[string][]chan runtime.Msg{
							"v": {printerOutPort},
						},
					},
				},
				{
					Ref: voidRef,
					IO: runtime.IO{
						In: map[string][]chan runtime.Msg{
							"v": {voidInPort},
						},
					},
				},
			},
		},
	}

	fmt.Println(
		r.Run(context.Background(), prog),
	)
}
