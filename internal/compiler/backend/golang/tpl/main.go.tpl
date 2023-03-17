package main

import (
    "context"
    "os"

    "github.com/emil14/neva/internal/runtime"
    "github.com/emil14/neva/internal/runtime/std/flow"
    "github.com/emil14/neva/internal/runtime/std/io"
)

func main() {
    // component refs
    {{range .ComponentRefs}}
    {{.VarName}} := runtime.ComponentRef{
        Pkg:  "{{.Pkg}}",
        Name: "{{.Name}}",
    }
    {{end}}

    // routine runner
    repo := map[runtime.ComponentRef]runtime.ComponentFunc{
        {{range .ComponentRefs}}
        {{.VarName}}: {{.Pkg}}.{{.FuncName}},
        {{end}}
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
    {{range .Ports}}
    {{.VarName}} := make(chan runtime.Msg)
    {{.VarName}}Addr := runtime.PortAddr{Name: "{{.Name}}"{{if .Path}}, Path:"{{.Path}}"{{end}}}

    {{end}}

    // Messages
    {{range .Messages}}
    {{.VarName}} := runtime.New{{.Type}}Msg({{.Value}})
    {{end}}

    prog := runtime.Program{
        Ports: map[runtime.PortAddr]chan runtime.Msg{
            {{range .PortConnections}}
            {{.SenderMeta.VarName}}: {{.Sender.VarName}},
            {{- range .Receivers }}
            {{- .Meta.VarName}}: {{.VarName}},
            {{- end }}
            {{end}}
        },
        Connections: []runtime.Connection{
            {{range .Connections}}
            {
                Sender: runtime.ConnectionSide{
                    Port: {{.Sender.VarName}},
                    Meta: runtime.ConnectionSideMeta{
                        PortAddr: {{.SenderMeta.VarName}},
                    },
                },
                Receivers: []runtime.ConnectionSide{
                    {{- range .Receivers }}
                    {
                        Port: {{.VarName}},
                        Meta: runtime.ConnectionSideMeta{
                            PortAddr: {{.Meta.VarName}},
                        },
                    },
                    {{- end }}
                },
            },
            {{end}}
        },
        Routines: runtime.Routines{
            Giver: []runtime.GiverRoutine{
                {{range .Givers}}
                {
                    OutPort: {{.OutPort.VarName}},
                    Msg:     {{.Msg.VarName}},
                },
                {{end}}
            },
            Component: []runtime.ComponentRoutine{
                {{range .ComponentRoutines}}
                {
                    Ref: {{.Ref.VarName}},
                    IO: runtime.ComponentIO{
                        In: map[string][]chan runtime.Msg{
                            {{range $key, $value := .In}}
                            "{{$key}}": { {{range $idx, $chan := $value}}{{if ne $idx 0}}, {{end}}{{$chan.VarName}}{{end}} },
                            {{end}}
                        },
                        Out: map[string][]chan runtime.Msg{
                            {{range $key, $value := .Out}}
                            "{{$key}}": { {{range $idx, $chan := $value}}{{if ne $idx 0}}, {{end}}{{$chan.VarName}}{{end}} },
                            {{end}}
                        },
                    },
                },
                {{end}}
            },
        },
    }

    exitCode, err := r.Run(context.Background(), prog)
    if err != nil {
        panic(err)
    }

    os.Exit(exitCode)
}
