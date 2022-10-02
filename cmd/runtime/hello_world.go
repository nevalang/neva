package main

import "github.com/emil14/neva/pkg/runtimesdk"

func helloWorld() *runtimesdk.Program {
	return prog(
		portAddr("in", "sig"),
		ports(
			port("in", "sig", 0),
			port("trigger.in", "sig", 0),
			port("trigger.out", "greeting", 0),
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
					portAddr("trigger.in", "sig"),
					portAddr("trigger.out", "greeting"),
					strMsg("hello world!\n"),
				),
			),
		),
		conns(
			conn(
				portAddr("in", "sig"),
				points(
					point(portAddr("trigger.in", "sig")),
				),
			),
			conn(
				portAddr("trigger.out", "greeting"),
				points(
					point(portAddr("print.in", "data")),
				),
			),
		),
	)
}
