package main

import (
	"context"
	"fmt"

	"github.com/emil14/neva/internal/runtime"
	"github.com/emil14/neva/internal/runtime/std/flow"
	"github.com/emil14/neva/internal/runtime/std/io"
)

func main() {
	// Component refs
	printerRef := runtime.ComponentRef{
		Pkg:  "io",
		Name: "printer",
	}
	triggerRef := runtime.ComponentRef{
		Pkg:  "flow",
		Name: "trigger",
	}

	// Routine runner
	repo := map[runtime.ComponentRef]runtime.ComponentFunc{
		printerRef: io.Print,
		triggerRef: flow.Trigger,
	}
	componentRunner := runtime.NewComponentRunner(repo)
	giverRunner := runtime.GiverRunnerImlp{}
	routineRunner := runtime.NewRoutineRunner(giverRunner, componentRunner)

	// Connector
	interceptor := runtime.InterceptorImlp{}
	connector := runtime.NewConnector(interceptor)

	// Runtime
	r := runtime.NewRuntime(connector, routineRunner)

	// Ports
	rootInStartPort := make(chan runtime.Msg)
	rootInStartPortAddr := runtime.PortAddr{Path: "root", Name: "sig"}
	rootOutExitPort := make(chan runtime.Msg)
	rootOutExitPortAddr := runtime.PortAddr{Path: "root", Name: "exit"}

	printerInPort := make(chan runtime.Msg)
	printerInPortAddr := runtime.PortAddr{Path: "printer.in", Name: "v"}

	printerOutPort := make(chan runtime.Msg)
	printerOutPortAddr := runtime.PortAddr{Path: "printer.out", Name: "v"}

	triggerInSigsPort := make(chan runtime.Msg)
	triggerInSigsAddr := runtime.PortAddr{Path: "trigger.in", Name: "sigs"}
	triggerInVPort := make(chan runtime.Msg)
	triggerInVAddr := runtime.PortAddr{Path: "trigger.in", Name: "v"}
	triggerOutVPort := make(chan runtime.Msg)
	triggerOutVPortAddr := runtime.PortAddr{
		Path: "trigger.out",
		Name: "v",
	}

	giverOutPort := make(chan runtime.Msg)
	giverOutPortAddr := runtime.PortAddr{
		Path: "giver.out",
		Name: "code",
	}

	// Messages
	exitCodeOneMsg := runtime.NewIntMsg(0)

	prog := runtime.Program{
		StartPortAddr: rootInStartPortAddr,
		Ports: map[runtime.PortAddr]chan runtime.Msg{
			// root
			rootInStartPortAddr: rootInStartPort,
			rootOutExitPortAddr: rootOutExitPort,
			// printer
			printerInPortAddr:  printerInPort,
			printerOutPortAddr: printerOutPort,
			// trigger
			triggerInSigsAddr:   triggerInSigsPort,
			triggerInVAddr:      triggerInVPort,
			triggerOutVPortAddr: triggerOutVPort,
			// giver
			giverOutPortAddr: giverOutPort,
		},
		Connections: []runtime.Connection{
			// root.start -> printer.in.v
			{
				Sender: runtime.ConnectionSide{
					Port: rootInStartPort,
					Meta: runtime.ConnectionSideMeta{
						PortAddr: rootInStartPortAddr,
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
			// printer.out.v -> trigger.in.sig
			{
				Sender: runtime.ConnectionSide{
					Port: printerOutPort,
					Meta: runtime.ConnectionSideMeta{
						PortAddr: printerOutPortAddr,
					},
				},
				Receivers: []runtime.ConnectionSide{
					{
						Port: triggerInSigsPort,
						Meta: runtime.ConnectionSideMeta{
							PortAddr: triggerInSigsAddr,
						},
					},
				},
			},
			// giver.out.code -> trigger.in.v
			{
				Sender: runtime.ConnectionSide{
					Port: giverOutPort,
					Meta: runtime.ConnectionSideMeta{
						PortAddr: giverOutPortAddr,
					},
				},
				Receivers: []runtime.ConnectionSide{
					{
						Port: triggerInVPort,
						Meta: runtime.ConnectionSideMeta{
							PortAddr: triggerInVAddr,
						},
					},
				},
			},
			// trigger.out.v -> root.out.exit
			{
				Sender: runtime.ConnectionSide{
					Port: triggerOutVPort,
					Meta: runtime.ConnectionSideMeta{
						PortAddr: triggerOutVPortAddr,
					},
				},
				Receivers: []runtime.ConnectionSide{
					{
						Port: rootOutExitPort,
						Meta: runtime.ConnectionSideMeta{
							PortAddr: rootOutExitPortAddr,
						},
					},
				},
			},
		},
		Routines: runtime.Routines{
			Giver: []runtime.GiverRoutine{
				{
					OutPort: giverOutPort,
					Msg:     exitCodeOneMsg,
				},
			},
			Component: []runtime.ComponentRoutine{
				// printer
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
				// trigger
				{
					Ref: triggerRef,
					IO: runtime.IO{
						In: map[string][]chan runtime.Msg{
							"sigs": {triggerInSigsPort},
							"v":   {triggerInVPort},
						},
						Out: map[string][]chan runtime.Msg{
							"v": {triggerOutVPort},
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
