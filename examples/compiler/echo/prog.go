package main

// import (
// 	"github.com/emil14/neva/internal/compiler/src"
// 	iohead "github.com/emil14/neva/internal/funcfx/io/header"
// )

// func helloWorld() src.Program {
// 	return src.Program{
// 		Pkgs: map[src.PkgRef]src.Pkg{
// 			{Name: "echo"}: {
// 				Imports: map[string]src.PkgRef{
// 					"utils": src.NewLocalPkgRef("utils"),
// 				},
// 				Types: map[string]src.Type{
// 					"users": {
// 						Body: src.TypeExpr{
// 							Ref: src.NewLocalTypeRef("list"),
// 							RefArgs: []src.TypeExpr{
// 								src.NewRefExpr(
// 									src.NewLocalTypeRef("user"),
// 								),
// 							},
// 						},
// 						Struct: map[string]src.TypeExpr{},
// 					},
// 					"user": {
// 						Body: src.NewRefExpr(
// 							src.NewLocalTypeRef("struct"),
// 						),
// 						Struct: map[string]src.TypeExpr{
// 							"name": src.NewRefExpr(
// 								src.NewLocalTypeRef("int"),
// 							),
// 							"age": src.NewRefExpr(
// 								src.NewLocalTypeRef("str"),
// 							),
// 							"friends": {
// 								Ref: src.NewLocalTypeRef("users"),
// 							},
// 						},
// 					},
// 				},
// 				Messages:   map[string]src.Msg{},
// 				Components: map[string]src.Component{},
// 				Root:       "root",
// 			},
// 			{Name: "utils"}: {
// 				Components: map[string]src.Component{
// 					"wrapper": {
// 						Interface: iohead.Headers()["print"],
// 						Nodes: src.Node{
// 							Effects: src.EffectNodes{
// 								Func: src.NewFuncRefSet("print"),
// 							},
// 						},
// 						Network: src.Network{
// 							{
// 								Sender: src.ConnectionSide{
// 									PortRef:    src.ConnectionPortRef{},
// 									ActionType: src.ReadStruct,
// 									Payload: src.ActionPayload{
// 										StructPath: []string{""},
// 									},
// 								},
// 								Receivers: []src.ConnectionSide{},
// 							},
// 						},
// 					},
// 				},
// 				Exports: map[string]src.PkgEntityRef{
// 					"print": {Type: src.ComponentEntity, LocalName: "wrapper"},
// 				},
// 			},
// 		},
// 		RootPkg: src.PkgRef{Name: "echo"},
// 	}
// }
