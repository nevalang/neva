package main

// import (
// 	h "github.com/emil14/neva/internal/pkg/runtimehelpers"
// 	"github.com/emil14/neva/pkg/runtimesdk"
// )

// func helloWorld() *runtimesdk.Program {
// 	return h.Prog(
// 		h.PortAddr("in", "sig"),
// 		h.Ports(
// 			h.SinglePort("in", "sig", 0),
// 			h.SinglePort("trigger.in", "sig", 0),
// 			h.SinglePort("trigger.out", "greeting", 0),
// 			h.SinglePort("print.in", "data", 0),
// 			h.SinglePort("print.out", "data", 0),
// 		),
// 		h.Effects(
// 			h.Operators(
// 				h.Operator(
// 					h.OperatorRef("io", "Print"),
// 					h.PortAddrs(h.PortAddr("print.in", "data")),
// 					h.PortAddrs(h.PortAddr("print.out", "data")),
// 				),
// 			),
// 			h.Constants(),
// 			h.Triggers(
// 				h.Trigger(
// 					h.PortAddr("trigger.in", "sig"),
// 					h.PortAddr("trigger.out", "greeting"),
// 					h.StrMsg("hello world!\n"),
// 				),
// 			),
// 		),
// 		h.Connections(
// 			h.Connection(
// 				h.PortAddr("in", "sig"),
// 				h.Points(
// 					h.Point(
// 						h.PortAddr("trigger.in", "sig"),
// 					),
// 				),
// 			),
// 			h.Connection(
// 				h.PortAddr("trigger.out", "greeting"),
// 				h.Points(
// 					h.Point(h.PortAddr("print.in", "data")),
// 				),
// 			),
// 		),
// 	)
// }
