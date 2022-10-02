package main

import (
	h "github.com/emil14/neva/internal/pkg/runtimehelpers"
	"github.com/emil14/neva/pkg/runtimesdk"
)

func helloWorld() *runtimesdk.Program {
	return h.Prog(
		h.PortAddr("in", "sig"),
		h.Ports(
			h.Port("in", "sig", 0),
			h.Port("read.in", "sig", 0),
			h.Port("read.out", "data", 0),
			h.Port("print.in", "data", 0),
			h.Port("print.out", "data", 0),
		),
		h.Effects(
			h.Operators(
				h.Operator(
					h.OperatorRef("io", "Println"),
					h.PortAddrs(h.PortAddr("print.in", "data")),
					h.PortAddrs(h.PortAddr("print.out", "data")),
				),
				h.Operator(
					h.OperatorRef("io", "Readln"),
					h.PortAddrs(h.PortAddr("read.in", "sig")),
					h.PortAddrs(h.PortAddr("read.out", "data")),
				),
			),
			h.Constants(),
			h.Triggers(),
		),
		h.Connections(
			h.Connection(
				h.PortAddr("in", "sig"),
				h.Points(
					h.Point(
						h.PortAddr("read.in", "sig"),
					),
				),
			),
			h.Connection(
				h.PortAddr("read.out", "data"),
				h.Points(
					h.Point(h.PortAddr("print.in", "data")),
				),
			),
			h.Connection(
				h.PortAddr("print.out", "data"),
				h.Points(
					h.Point(h.PortAddr("read.in", "sig")),
				),
			),
		),
	)
}
