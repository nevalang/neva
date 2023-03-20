//go:build integration
// +build integration

package analyzer_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/emil14/neva/internal/compiler"
	"github.com/emil14/neva/internal/compiler/analyzer"
	"github.com/emil14/neva/internal/compiler/helper"
	ts "github.com/emil14/neva/pkg/types"
)

var h helper.Helper

func TestAnalyzer(t *testing.T) {
	t.Parallel()

	type testcase struct {
		enabled bool
		name    string
		prog    compiler.Program
		wantErr error
	}

	tests := []testcase{
		{
			name: "root_pkg_refers_to_imported_pkg",
			prog: compiler.Program{
				Pkgs: map[string]compiler.Pkg{
					"pkg2": {
						Entities: map[string]compiler.Entity{
							"t1": h.TypeEntity(
								true,
								h.Def( // type t1<a> = vec<a>
									h.Inst("vec", h.Inst("a")),
									h.ParamWithNoConstr("a"),
								),
							),
							"c1": {
								Exported: true,
								Kind:     compiler.ComponentEntity,
							},
						},
					},
					"main": {
						Imports: h.Imports("pkg2"),
						Entities: map[string]compiler.Entity{
							"t1": h.TypeEntity(
								true,
								h.Def( // type t1 = pkg2.t1<int>
									h.Inst("pkg2.t1", h.Inst("int")),
								),
							),
							"main": h.MainComponent(map[string]compiler.Node{
								"n1": h.Node(
									h.NodeInstance("pkg1", "c1"),
								),
							}, nil),
						},
					},
				},
			},
			wantErr: nil,
		},
		{
			name: "root_pkg_refers_imported_pkg_that_refers_another_imported_pkg",
			prog: compiler.Program{
				Pkgs: map[string]compiler.Pkg{
					"pkg3": {
						Entities: map[string]compiler.Entity{
							"t1": h.TypeEntity(
								true,
								h.Def( // type t1<a> = vec<a>
									h.Inst("vec", h.Inst("a")),
									h.ParamWithNoConstr("a"),
								),
							),
						},
					},
					"pkg2": {
						Imports: h.Imports("pkg3"),
						Entities: map[string]compiler.Entity{
							"t1": h.TypeEntity(
								true,
								h.Def( // type t1<a> = t1<a>
									h.Inst("pkg3.t1", h.Inst("a")),
									h.ParamWithNoConstr("a"),
								),
							),
							"c1": {
								Exported: true,
								Kind:     compiler.ComponentEntity,
							},
						},
					},
					"main": {
						Imports: h.Imports("pkg2"),
						Entities: map[string]compiler.Entity{
							"t1": h.TypeEntity(
								true,
								h.Def( // type t1 = pkg2.t1<int>
									h.Inst("pkg2.t1", h.Inst("int")),
								),
							),
							"main": h.MainComponent(map[string]compiler.Node{
								"n1": h.Node(
									h.NodeInstance("pkg1", "c1"),
								),
							}, nil),
						},
					},
				},
			},
			wantErr: nil,
		},
		{ // FIXME false-positive
			name: "inassignable_message_and_4-step_import_chain",
			prog: compiler.Program{
				Pkgs: map[string]compiler.Pkg{
					"pkg1": {
						Entities: map[string]compiler.Entity{
							"m1": h.IntMsg(true, 42),
							"c1": {
								Exported: true,
								Kind:     compiler.ComponentEntity,
							},
						},
					},
					"pkg2": {
						Imports: h.Imports("pkg1"),
						Entities: map[string]compiler.Entity{
							"m1": h.MsgWithRefEntity(true, &compiler.EntityRef{
								Pkg:  "pkg1",
								Name: "m1",
							}),
						},
					},
					"pkg3": {
						Imports: h.Imports("pkg1", "pkg2"),
						Entities: map[string]compiler.Entity{
							"m1": h.IntVecMsgEntity(
								true,
								[]compiler.Msg{
									{
										Ref: &compiler.EntityRef{
											Pkg:  "pkg1",
											Name: "m1",
										},
									},
									{
										Ref: &compiler.EntityRef{
											Pkg:  "pkg2",
											Name: "m1",
										},
									},
									{Value: h.IntMsgValue(43)},
								},
							),
						},
					},
					"main": {
						Imports: h.Imports("pkg1", "pkg2", "pkg3"),
						Entities: map[string]compiler.Entity{
							"m1": h.IntVecMsgEntity(
								true,
								[]compiler.Msg{
									{Value: h.IntMsgValue(44)},
									{
										Ref: &compiler.EntityRef{
											Pkg:  "pkg3",
											Name: "m1",
										},
									},
								},
							),
							"main": h.MainComponent(map[string]compiler.Node{
								"n1": h.Node(h.NodeInstance("pkg1", "c1")),
							}, nil),
						},
					},
				},
			},
			wantErr: analyzer.ErrVecEl,
		},
		{
			enabled: true,
			name:    "pkg1_imports_pkg2_and_pkg3_but_refers_to_only_pkg2_while_pkg2_actually_refers_pkg3",
			prog: compiler.Program{
				Pkgs: map[string]compiler.Pkg{
					"pkg1": {
						Imports: h.Imports("pkg2", "pkg3"), // pkg3 unused
						Entities: map[string]compiler.Entity{
							"main": h.MainComponent(map[string]compiler.Node{
								"n1": h.Node(h.NodeInstance("pkg1", "c1")),
							}, nil),
						},
					},
					"pkg2": {
						Imports: map[string]string{
							"pkg3": "pkg3",
						},
						Entities: map[string]compiler.Entity{
							"c1": {
								Exported: true,
								Kind:     compiler.ComponentEntity,
							},
							"m1": h.MsgWithRefEntity(true, &compiler.EntityRef{
								Pkg:  "pkg3",
								Name: "m1",
							}),
						},
					},
					"pkg3": {
						Entities: map[string]compiler.Entity{
							"m1": h.IntMsg(true, 42),
						},
					},
				},
			},
			wantErr: analyzer.ErrUnusedImport,
		},
	}

	a := analyzer.MustNew(
		ts.NewDefaultResolver(),
		ts.NewDefaultCompatChecker(),
		ts.Validator{},
	)

	for _, tt := range tests {
		if !tt.enabled {
			continue
		}

		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			_, err := a.Analyze(context.Background(), tt.prog)
			fmt.Println(err)
			assert.ErrorIs(t, err, tt.wantErr)
		})
	}
}
