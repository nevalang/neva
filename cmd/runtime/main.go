package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/emil14/neva/internal/core"
	"github.com/emil14/neva/internal/runtime"
	"github.com/emil14/neva/internal/runtime/connector"
	"github.com/emil14/neva/internal/runtime/constspawner"
	"github.com/emil14/neva/internal/runtime/decoder"
	"github.com/emil14/neva/internal/runtime/opspawner"
	"github.com/emil14/neva/internal/runtime/opspawner/repo"
	"github.com/emil14/neva/internal/runtime/portgen"
	"github.com/emil14/neva/pkg/runtimesdk"
	"github.com/golang/protobuf/proto"
)

func main() {
	fmt.Println(
		core.NewListMsg([]core.Msg{
			core.NewIntMsg(42),
			core.NewDictMsg(map[string]core.Msg{
				"name": core.NewStrMsg("John"),
				"age":  core.NewIntMsg(42),
				"friends": core.NewListMsg([]core.Msg{
					core.NewDictMsg(map[string]core.Msg{
						"name":    core.NewStrMsg("John"),
						"age":     core.NewIntMsg(42),
						"friends": core.NewListMsg([]core.Msg{}),
					}),
				}),
			}),
		}),
	)

	r := runtime.MustNew(
		decoder.MustNewProto(
			decoder.NewCaster(),
			decoder.NewUnmarshaler(),
		),
		portgen.New(),
		opspawner.MustNew(
			repo.NewPlugin(map[string]repo.PluginData{
				"io": {
					Filepath: "/home/evaleev/projects/neva/plugins/print.so",
					Exports:  []string{"Print"},
				},
			}),
			opspawner.Searcher{},
		),
		constspawner.Spawner{},
		connector.MustNew(
			connector.LoggingInterceptor{},
		),
	)

	helloWorld := runtimesdk.Program{
		StartPort: &runtimesdk.PortAddr{
			Path: "in",
			Port: "start",
		},
		Ports: []*runtimesdk.PortAddr{
			{Path: "in", Port: "start"},
			{Path: "const", Port: "greeting"},
			{Path: "write", Port: "kick"},
			{Path: "write", Port: "msg"},
		},
		Operators: []*runtimesdk.Operator{
			{
				Ref: &runtimesdk.OperatorRef{
					Pkg:  "io",
					Name: "Print",
				},
				InPortAddrs: []*runtimesdk.PortAddr{
					{Path: "write", Port: "kick"},
					{Path: "write", Port: "msg"},
				},
			},
		},
		Constants: []*runtimesdk.Constant{
			{
				OutPortAddr: &runtimesdk.PortAddr{
					Path: "const", Port: "greeting",
				},
				Msg: &runtimesdk.Msg{
					Str:  "hello world!",
					Type: runtimesdk.MsgType_VALUE_TYPE_STR,
				},
			},
		},
		Connections: []*runtimesdk.Connection{
			{
				SenderOutPortAddr: &runtimesdk.PortAddr{
					Path: "in", Port: "start",
				},
				ReceiverConnectionPoints: []*runtimesdk.ConnectionPoint{
					{
						InPortAddr: &runtimesdk.PortAddr{Path: "write", Port: "kick"},
					},
				},
			},
			{
				SenderOutPortAddr: &runtimesdk.PortAddr{
					Path: "const", Port: "greeting",
				},
				ReceiverConnectionPoints: []*runtimesdk.ConnectionPoint{
					{
						InPortAddr: &runtimesdk.PortAddr{Path: "write", Port: "msg"},
					},
				},
			},
		},
	}

	bb, err := proto.Marshal(&helloWorld)
	if err != nil {
		panic(err)
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	if err := r.Run(ctx, bb); err != nil {
		log.Println(err)
	}
}
