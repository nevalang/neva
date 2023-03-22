package irgen_test

import (
	"context"
	"testing"

	"github.com/emil14/neva/internal/compiler"
	"github.com/emil14/neva/internal/compiler/helper"
	"github.com/emil14/neva/internal/compiler/ir"
	"github.com/emil14/neva/internal/compiler/irgen"
	"github.com/stretchr/testify/assert"
)

var h helper.Helper

func TestGenerator_Generate(t *testing.T) {
	tests := []struct {
		name    string
		prog    compiler.Program
		native  map[compiler.EntityRef]ir.FuncRef
		want    ir.Program
		wantErr bool
	}{
		// - pkgs==nil test
		{
			name: "program that does nothing",
			prog: compiler.Program{
				Pkgs: map[string]compiler.Pkg{
					"main": {
						Imports: h.Imports("flow"),
						Entities: map[string]compiler.Entity{
							"code": h.IntMsg(false, 0),
							"main": h.MainComponent(map[string]compiler.Node{
								"trigger": h.NodeWithStaticPorts(
									h.NodeInstance("flow", "Trigger", h.Rec(nil)),
									map[compiler.RelPortAddr]compiler.EntityRef{
										{Name: "v"}: {Name: "code"},
									},
								),
							}, []compiler.Connection{
								{
									SenderSide: compiler.ConnectionSide{
										PortAddr: compiler.ConnPortAddr{
											Node: "in",
											RelPortAddr: compiler.RelPortAddr{
												Name: "start",
											},
										},
									},
									ReceiverSides: []compiler.ConnectionSide{
										{
											PortAddr: compiler.ConnPortAddr{
												Node: "trigger.",
												RelPortAddr: compiler.RelPortAddr{
													Name: "sig",
												},
											},
										},
									},
								},
								{
									SenderSide: compiler.ConnectionSide{
										PortAddr: compiler.ConnPortAddr{
											Node: "trigger",
											RelPortAddr: compiler.RelPortAddr{
												Name: "v",
											},
										},
									},
									ReceiverSides: []compiler.ConnectionSide{
										{
											PortAddr: compiler.ConnPortAddr{
												Node: "out.exit.",
												RelPortAddr: compiler.RelPortAddr{
													Name: "sig",
												},
											},
										},
									},
								},
							}),
						},
					},
				},
			},
			native: map[compiler.EntityRef]ir.FuncRef{
				{Pkg: "flow", Name: "Trigger"}: {Pkg: "flow", Name: "Trigger"},
			},
			want: ir.Program{
				Ports: map[ir.PortAddr]uint8{
					{Name: "start"}: 0,
					{Name: "exit"}:  0,
					{Path: "trigger.in", Name: "sigs", Idx: 0}: 0,
					{Path: "trigger.in", Name: "v", Idx: 0}:    0,
					{Path: "trigger.out", Name: "v", Idx: 0}:   0,
					{Path: "giver.out", Name: "code", Idx: 0}:  0,
				},
				Routines: ir.Routines{
					Giver: map[ir.PortAddr]ir.Msg{
						{Path: "giver.out", Name: "code"}: {Type: ir.IntMsg, Int: 0},
					},
					Func: []ir.FuncRoutine{
						{
							Ref: ir.FuncRef{Pkg: "flow", Name: "Trigger"},
							IO: ir.FuncIO{
								In: []ir.PortAddr{
									{Path: "trigger.in", Name: "sigs"},
									{Path: "trigger.in", Name: "v"},
								},
								Out: []ir.PortAddr{
									{Path: "trigger.out", Name: "v"},
								},
							},
						},
					},
				},
				Connections: []ir.Connection{
					{
						SenderSide: ir.ConnectionSide{
							PortAddr: ir.PortAddr{Name: "start"},
						},
						ReceiverSides: []ir.ConnectionSide{
							{
								PortAddr: ir.PortAddr{Path: "trigger.in", Name: "sigs"},
							},
						},
					},
					{
						SenderSide: ir.ConnectionSide{
							PortAddr: ir.PortAddr{Path: "giver.out", Name: "code"},
						},
						ReceiverSides: []ir.ConnectionSide{
							{
								PortAddr: ir.PortAddr{Path: "trigger.in", Name: "v"},
							},
						},
					},
					{
						SenderSide: ir.ConnectionSide{
							PortAddr: ir.PortAddr{Path: "trigger.out", Name: "v"},
						},
						ReceiverSides: []ir.ConnectionSide{
							{
								PortAddr: ir.PortAddr{Name: "exit"},
							},
						},
					},
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := irgen.New(tt.native)
			got, err := g.Generate(context.Background(), tt.prog)

			assert.Equal(t, tt.want, got)
			assert.Equal(t, tt.wantErr, err)
		})
	}
}
