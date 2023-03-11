// generator/main.go
package main

import (
	"context"
	"fmt"

	"github.com/emil14/neva/internal/runtime"
	"github.com/emil14/neva/internal/runtime/std/io"
)

func main() {
	repo := map[runtime.ComponentRef]func(context.Context, runtime.IO) error{
		{
			Pkg:  "io",
			Name: "Printer",
		}: io.Print,
		{
			Pkg:  "io",
			Name: "void",
		}: io.Void,
	}

	connector := runtime.NewConnector(runtime.InterceptorImlp{})

	componentRunner := runtime.NewComponentRunner(repo)
	routineRunner := runtime.NewRoutineRunner(runtime.GiverRunnerImlp{}, componentRunner)

	r := runtime.NewRuntime(connector, routineRunner)

	startPortAddr := runtime.PortAddr{
		Path: "root",
		Name: "sig",
	}
	startPortCh := make(chan runtime.Msg)

	printerInPort := runtime.PortAddr{
		Path: "printer.in",
		Name: "v",
	}
	printerInCh := make(chan runtime.Msg)

	printerOutPort := runtime.PortAddr{
		Path: "printer.out",
		Name: "v",
	}
	printerOutCh := make(chan runtime.Msg)

	voidInCh := make(chan runtime.Msg)
	voidInPort := runtime.PortAddr{
		Path: "void.in",
		Name: "v",
	}

	prog := runtime.Program{
		StartPortAddr: startPortAddr,
		Ports: map[runtime.PortAddr]chan runtime.Msg{
			startPortAddr:  printerInCh,
			printerInPort:  startPortCh,
			printerInPort:  printerInCh,
			printerOutPort: printerOutCh,
		},
		Connections: []runtime.Connection{
			{
				Sender: runtime.ConnectionSide{
					Port: startPortCh,
					Meta: runtime.ConnectionSideMeta{
						PortAddr: startPortAddr,
					},
				},
				Receivers: []runtime.ConnectionSide{
					{
						Port: printerInCh,
						Meta: runtime.ConnectionSideMeta{
							PortAddr: startPortAddr,
						},
					},
				},
			},
			{
				Sender: runtime.ConnectionSide{
					Port: printerOutCh,
					Meta: runtime.ConnectionSideMeta{
						PortAddr: printerOutPort,
					},
				},
				Receivers: []runtime.ConnectionSide{
					{
						Port: voidInCh,
						Meta: runtime.ConnectionSideMeta{
							PortAddr: voidInPort,
						},
					},
				},
			},
		},
		Routines: runtime.Routines{
			Component: []runtime.ComponentRoutine{
				{
					Ref: runtime.ComponentRef{},
					IO: runtime.IO{
						In: map[string][]chan runtime.Msg{
							"v": {printerInCh},
						},
						Out: map[string][]chan runtime.Msg{
							"v": {printerOutCh},
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
