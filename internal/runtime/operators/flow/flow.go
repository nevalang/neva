package main

import (
	"context"

	"github.com/emil14/neva/internal/runtime/core"
)

func Locker(ctx context.Context, io core.IO) error {
	slots, _ := io.In.ArrPortSlots("sigs")
	v, _ := io.In.SinglePort("v") // not arrport because 1) one-data-many-lockers 2) use many lockers
	out, _ := io.Out.SinglePort("v")

	for {
		for _, sig := range slots {
			<-sig
		}
		msg := <-v
		out <- msg
	}
}

// [x,y,z] -> (x, y, z)
func ArrPortStreamer(ctx context.Context, io core.IO) error {
	in, _ := io.Out.ArrPortSlots("in")
	out, _ := io.Out.SinglePort("out")

	for {
		out <- core.NewDictMsg(map[string]core.Msg{
			"open": core.NewBoolMsg(true),
		})

		for _, slot := range in {
			msg := <-slot
			out <- core.NewDictMsg(map[string]core.Msg{
				"open": core.NewBoolMsg(false),
				"v":    msg,
			})
		}

		out <- core.NewDictMsg(map[string]core.Msg{
			"open": core.NewBoolMsg(false),
		})
	}
}

// 0: [1,2,3]; 1: [4,5,6] -> (6 5 4) (3 2 1)
func FlatListStreamer(ctx context.Context, io core.IO) error {
	in, _ := io.Out.ArrPortSlots("in")
	out, _ := io.Out.SinglePort("out")

	for {
		out <- core.NewDictMsg(map[string]core.Msg{
			"open": core.NewBoolMsg(true),
		})

		for _, slot := range in {
			msg := <-slot

			for _, v := range msg.List() {
				out <- core.NewDictMsg(map[string]core.Msg{
					"open": core.NewBoolMsg(false),
					"v":    v,
				})
			}
		}

		out <- core.NewDictMsg(map[string]core.Msg{
			"open": core.NewBoolMsg(false),
		})
	}
}

// func StreamToListAdapter(ctx context.Context, io core.IO) error {
// 	in := io.Out.ArrPortSlots("in")
// 	out := io.Out.SinglePort("out")
// }

// 0: [1,2,3]; 1: [4,5,6] -> ( (6 5 4) (3 2 1) )
func NestedListStreamer(ctx context.Context, io core.IO) error {
	in, _ := io.Out.ArrPortSlots("in")
	out, _ := io.Out.SinglePort("out")

	for {
		out <- core.NewDictMsg(map[string]core.Msg{
			"open": core.NewBoolMsg(true),
		})

		for _, slot := range in {
			out <- core.NewDictMsg(map[string]core.Msg{
				"open": core.NewBoolMsg(true),
			})

			msg := <-slot

			for _, v := range msg.List() {
				out <- core.NewDictMsg(map[string]core.Msg{
					"open": core.NewBoolMsg(false),
					"v":    v,
				})
			}

			out <- core.NewDictMsg(map[string]core.Msg{
				"open": core.NewBoolMsg(false),
			})
		}

		out <- core.NewDictMsg(map[string]core.Msg{
			"open": core.NewBoolMsg(false),
		})
	}
}

// [ {a:3,y:4}, {a:1,y:2} ] -> ({a,3}, {y,4}) ({a,1}, {y,2})
func FlatDictStreamer(ctx context.Context, io core.IO) error {
	in, _ := io.Out.ArrPortSlots("in")
	out, _ := io.Out.SinglePort("out")

	for {
		out <- core.NewDictMsg(map[string]core.Msg{
			"open": core.NewBoolMsg(true),
		})

		for _, slot := range in {
			msg := <-slot

			for k, v := range msg.Dict() {
				out <- core.NewDictMsg(map[string]core.Msg{
					"k": core.NewStrMsg(k),
					"v": v,
				})
			}
		}

		out <- core.NewDictMsg(map[string]core.Msg{
			"open": core.NewBoolMsg(false),
		})
	}
}

// [ {a:3,y:4}, {a:1,y:2} ] -> ( ({a,3}, {y,4}) ({a,1}, {y,2}) )
func NestedDictStreamer(ctx context.Context, io core.IO) error {
	in, _ := io.Out.ArrPortSlots("in")
	out, _ := io.Out.SinglePort("out")

	for {
		out <- core.NewDictMsg(map[string]core.Msg{
			"open": core.NewBoolMsg(true),
		})

		for _, slot := range in {
			out <- core.NewDictMsg(map[string]core.Msg{
				"open": core.NewBoolMsg(true),
			})

			msg := <-slot

			for k, v := range msg.Dict() {
				out <- core.NewDictMsg(map[string]core.Msg{
					"k": core.NewStrMsg(k),
					"v": v,
				})
			}

			out <- core.NewDictMsg(map[string]core.Msg{
				"open": core.NewBoolMsg(false),
			})
		}

		out <- core.NewDictMsg(map[string]core.Msg{
			"open": core.NewBoolMsg(false),
		})
	}
}
