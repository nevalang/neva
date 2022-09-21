package main

import "github.com/emil14/neva/pkg/runtimesdk"

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
