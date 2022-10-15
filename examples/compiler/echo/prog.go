package main

import (
	"github.com/emil14/neva/internal/compiler/src"
	iohead "github.com/emil14/neva/internal/funcfx/io/header"
)

func helloWorld() src.Program {
	return src.Program{
		Packages: map[src.PkgRef]src.Package{
			{Name: "echo"}: {
				Imports: map[string]src.PkgRef{
					"utils": src.NewLocalPkgRef("utils"),
				},
				Types: map[string]src.Type{
					"users": {
						Expr: src.TypeExpr{
							Ref: src.NewLocalTypeRef("list"),
							RefArgs: []src.TypeExpr{
								src.NewRefExpr(
									src.NewLocalTypeRef("user"),
								),
							},
						},
						StructExpr: map[string]src.TypeExpr{},
					},
					"user": {
						Expr: src.NewRefExpr(
							src.NewLocalTypeRef("struct"),
						),
						StructExpr: map[string]src.TypeExpr{
							"name": src.NewRefExpr(
								src.NewLocalTypeRef("int"),
							),
							"age": src.NewRefExpr(
								src.NewLocalTypeRef("str"),
							),
							"friends": {
								Ref: src.NewLocalTypeRef("users"),
							},
						},
					},
				},
				Messages:      map[string]src.MsgDef{},
				Components:    map[string]src.Component{},
				RootComponent: "root",
			},
			{Name: "utils"}: {
				Components: map[string]src.Component{
					"wrapper": {
						Header: iohead.Headers()["print"],
						Nodes: src.Nodes{
							Effects: src.EffectNodes{
								Func: src.NewFuncRefSet("print"),
							},
						},
						Net: src.Network{
							{
								Sender: src.ConnectionSide{
									PortRef:    src.PortRef{},
									ActionType: src.ReadStruct,
									Payload: src.ActionPayload{
										StructPath: []string{""},
									},
								},
								Receivers: []src.ConnectionSide{},
							},
						},
					},
				},
				Exports: map[string]src.ExportRef{
					"print": {Type: src.ComponentExport, LocalName: "wrapper"},
				},
			},
		},
		RootPkg: src.PkgRef{Name: "echo"},
	}
}
