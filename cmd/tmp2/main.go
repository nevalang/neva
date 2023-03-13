// generator/main.go
package main

import (
	"bytes"
	"os"

	"github.com/emil14/neva/internal"
)

func main() {
	createGoModFile()

	os.MkdirAll("", os.ModeDir)

	runtimeBb, err := internal.RuntimeFiles.ReadFile("runtime/runtime.go")
	if err != nil {
		panic(err)
	}

	var buf bytes.Buffer

	buf.Write(runtimeBb)

	f, err := os.Create("tmp/runtime/runtime.go")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	if _, err := buf.WriteTo(f); err != nil {
		panic(err)
	}

	buf.Reset()

	if _, err := buf.WriteString(prog); err != nil {
		panic(err)
	}

	f, err = os.Create("tmp/main.go")
	if err != nil {
		panic(err)
	}
	defer f.Close()
}

func createGoModFile() {
	f, err := os.Create("tmp/go.mod")
	if err != nil {
		panic(err)
	}

	defer func() {
		if err := f.Close(); err != nil {
			panic(err)
		}
	}()

	_, err = f.WriteString("module root")
	if err != nil {
		panic(err)
	}
}

var prog = `
package main

import (
	"context"
	"fmt"

	"root/runtime"
	"root/std/flow"
	"root/std/io"
)

func main() {
	// component refs
	printerRef := runtime.ComponentRef{
		Pkg:  "io",
		Name: "printer",
	}
	voidRef := runtime.ComponentRef{
		Pkg:  "io",
		Name: "void",
	}

	// component refs to std functions map
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
`
