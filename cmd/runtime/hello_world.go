package main

import "github.com/emil14/neva/pkg/runtimesdk"

func helloWorld() *runtimesdk.Program {
	return prog(
		portAddr("in", "sig"),
		ports(
			port("in", "sig", 0),
			port("const.in", "sig", 0),
			port("const.out", "greeting", 0),
			port("print.in", "data", 0),
			port("print.out", "data", 0),
		),
		effects(
			ops(
				op(
					opref("io", "Print"),
					portsAddrs(portAddr("print.in", "data")),
					portsAddrs(portAddr("print.out", "data")),
				),
			),
			consts(),
			triggers(
				trigger(
					portAddr("const.in", "sig"),
					portAddr("const.out", "greeting"),
					strMsg("hello world!\n"),
				),
			),
		),
		conns(
			conn(
				portAddr("in", "sig"),
				points(
					point(portAddr("const.in", "sig")),
				),
			),
			conn(
				portAddr("const.out", "greeting"),
				points(
					point(portAddr("print.in", "data")),
				),
			),
		),
	)
}
