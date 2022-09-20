package main

import (
	"context"
	"log"

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
	r := mustCreateRuntime()
	hw := helloWorld()

	bb, err := proto.Marshal(hw)
	if err != nil {
		panic(err)
	}

	log.SetFlags(log.Flags() &^ (log.Ldate | log.Ltime))

	if err := r.Run(context.Background(), bb); err != nil {
		log.Println(err)
	}
}

func helloWorld() *runtimesdk.Program {
	return prog(
		port("in", "sig"),
		ports(
			port("in", "sig"),
			port("const", "greeting"),
			port("lock.in", "sig"),
			port("lock.in", "data"),
			port("lock.out", "data"),
			port("print.in", "data"),
			port("print.out", "data"),
		),
		ops(
			op(
				opref("flow", "Lock"),
				ports(
					port("lock.in", "sig"),
					port("lock.in", "data"),
				),
				ports(port("lock.out", "data")),
			),
			op(
				opref("io", "Print"),
				ports(port("print.in", "data")),
				ports(port("print.out", "data")),
			),
		),
		consts(
			cnst(
				port("const", "greeting"),
				strmsg("hello world!\n"),
			),
		),
		conns(
			conn(
				port("in", "sig"),
				points(
					point(port("lock.in", "sig")),
				),
			),
			conn(
				port("const", "greeting"),
				points(
					point(port("lock.in", "data")),
				),
			),
			conn(
				port("lock.out", "data"),
				points(
					point(port("print.in", "data")),
				),
			),
		),
	)
}

// sdk helpers

func prog(
	start *runtimesdk.PortAddr,
	ports []*runtimesdk.PortAddr,
	ops []*runtimesdk.Operator,
	consts []*runtimesdk.Constant,
	conns []*runtimesdk.Connection,
) *runtimesdk.Program {
	return &runtimesdk.Program{
		StartPort:   start,
		Ports:       ports,
		Operators:   ops,
		Constants:   consts,
		Connections: conns,
	}
}

func conns(cc ...*runtimesdk.Connection) []*runtimesdk.Connection {
	return cc
}

func conn(sender *runtimesdk.PortAddr, receivers []*runtimesdk.ConnectionPoint) *runtimesdk.Connection {
	return &runtimesdk.Connection{
		SenderOutPortAddr:        sender,
		ReceiverConnectionPoints: receivers,
	}
}

func points(pp ...*runtimesdk.ConnectionPoint) []*runtimesdk.ConnectionPoint {
	return pp
}

func point(in *runtimesdk.PortAddr) *runtimesdk.ConnectionPoint {
	return &runtimesdk.ConnectionPoint{
		InPortAddr:      in,
		Type:            0,
		StructFieldPath: []string{},
	}
}

func consts(cc ...*runtimesdk.Constant) []*runtimesdk.Constant {
	return cc
}

func cnst(out *runtimesdk.PortAddr, msg *runtimesdk.Msg) *runtimesdk.Constant {
	return &runtimesdk.Constant{
		OutPortAddr: out,
		Msg:         msg,
	}
}

func strmsg(s string) *runtimesdk.Msg {
	return &runtimesdk.Msg{
		Str:  "hello world!\n",
		Type: runtimesdk.MsgType_VALUE_TYPE_STR, //nolint
	}
}

func ops(oo ...*runtimesdk.Operator) []*runtimesdk.Operator {
	return oo
}

func op(ref *runtimesdk.OperatorRef, in, out []*runtimesdk.PortAddr) *runtimesdk.Operator {
	return &runtimesdk.Operator{
		Ref:          ref,
		InPortAddrs:  in,
		OutPortAddrs: out,
	}
}

func opref(pkg, name string) *runtimesdk.OperatorRef {
	return &runtimesdk.OperatorRef{
		Pkg: pkg, Name: name,
	}
}

func ports(pp ...*runtimesdk.PortAddr) []*runtimesdk.PortAddr {
	return pp
}

func port(path, name string) *runtimesdk.PortAddr {
	return slot(path, name, 0)
}

func slot(path, name string, idx uint32) *runtimesdk.PortAddr {
	return &runtimesdk.PortAddr{
		Path: path, Port: name, Idx: idx,
	}
}

// runtime

func mustCreateRuntime() runtime.Runtime {
	r := runtime.MustNew(
		decoder.MustNewProto(
			decoder.NewCaster(),
			decoder.NewUnmarshaler(),
		),
		portgen.New(),
		opspawner.MustNew(
			repo.NewPlugin(map[string]repo.Package{
				"flow": {
					Filepath: "/home/evaleev/projects/neva/plugins/lock.so",
					Exports:  []string{"Lock"},
				},
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
	return r
}
