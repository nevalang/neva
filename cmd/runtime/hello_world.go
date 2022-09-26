package main

import "github.com/emil14/neva/pkg/runtimesdk"

func helloWorld() *runtimesdk.Program {
	return prog(
		portAddr("in", "sig"),
		ports(
			port("in", "sig", 0),
			port("const", "greeting", 0),
			port("lock.in", "sig", 0),
			port("lock.in", "data", 0),
			port("lock.out", "data", 0),
			port("print.in", "data", 0),
			port("print.out", "data", 0),
		),
		ops(
			op(
				opref("flow", "Lock"),
				portsAddrs(
					portAddr("lock.in", "sig"),
					portAddr("lock.in", "data"),
				),
				portsAddrs(portAddr("lock.out", "data")),
			),
			op(
				opref("io", "Print"),
				portsAddrs(portAddr("print.in", "data")),
				portsAddrs(portAddr("print.out", "data")),
			),
		),
		consts(
			cnst(
				portAddr("const", "greeting"),
				strmsg("hello world!\n"),
			),
		),
		conns(
			conn(
				portAddr("in", "sig"),
				points(
					point(portAddr("lock.in", "sig")),
				),
			),
			conn(
				portAddr("const", "greeting"),
				points(
					point(portAddr("lock.in", "data")),
				),
			),
			conn(
				portAddr("lock.out", "data"),
				points(
					point(portAddr("print.in", "data")),
				),
			),
		),
	)
}
