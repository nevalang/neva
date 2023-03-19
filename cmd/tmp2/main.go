package main

import (
	"bytes"
	"io/fs"
	"os"

	"github.com/emil14/neva/internal"
	"github.com/emil14/neva/internal/compiler/backend/golang"
	"github.com/emil14/neva/internal/compiler/irgen"
	"github.com/emil14/neva/internal/compiler/src"
	ts "github.com/emil14/neva/pkg/types"
)

var efs = internal.RuntimeFiles
var basePath = "/home/evaleev/projects/tmp"

func main() {
	if err := os.RemoveAll(basePath); err != nil {
		panic(err)
	}

	if err := os.MkdirAll(basePath, os.ModePerm); err != nil {
		panic(err)
	}

	putGoMod()

	putRuntime()

	prog := src.Program{
		Pkgs: map[string]src.Pkg{
			"main": {
				Imports: map[string]string{
					"io":   "io",
					"flow": "flow",
				},
				Entities: map[string]src.Entity{
					"code": {
						Kind: src.MsgEntity,
						Msg: src.Msg{
							Value: src.MsgValue{
								Type: ts.Expr{
									Inst: ts.InstExpr{
										Ref: "int",
									},
								},
								Int: 0,
							},
						},
					},
					"main": {
						Kind: src.ComponentEntity,
						Component: src.Component{
							IO: src.IO{
								In: map[string]src.Port{
									"start": {
										Type: ts.Expr{
											Lit: ts.LitExpr{
												Rec: map[string]ts.Expr{},
											},
										},
										IsArr: false,
									},
								},
								Out: map[string]src.Port{
									"exit": {
										Type: ts.Expr{
											Inst: ts.InstExpr{
												Ref: "int",
											},
										},
										IsArr: false,
									},
								},
							},
							Nodes: map[string]src.Node{
								"print": {
									Instance: src.NodeInstance{
										Ref: src.EntityRef{
											Pkg:  "io",
											Name: "Print",
										},
										TypeArgs: []ts.Expr{
											{
												Inst: ts.InstExpr{
													Ref: "str",
												},
											},
										},
									},
								},
								"trigger": {
									Instance: src.NodeInstance{
										Ref: src.EntityRef{
											Pkg:  "flow",
											Name: "Trigger",
										},
										TypeArgs: []ts.Expr{
											{
												Lit: ts.LitExpr{
													Rec: map[string]ts.Expr{},
												},
											},
										},
										DIArgs: map[string]src.NodeInstance{},
									},
									StaticInports: map[src.RelPortAddr]src.EntityRef{
										{
											Name: "v",
											Idx:  0,
										}: {
											Pkg:  "",
											Name: "code",
										},
									},
								},
							},
							Net: []src.Connection{
								{
									SenderSide: src.ConnectionSide{
										PortAddr: src.ConnPortAddr{
											Node: "in",
											RelPortAddr: src.RelPortAddr{
												Name: "start",
												Idx:  0,
											},
										},
										Selectors: []src.Selector{},
									},
									ReceiverSides: []src.ConnectionSide{
										{
											PortAddr: src.ConnPortAddr{
												Node: "print",
												RelPortAddr: src.RelPortAddr{
													Name: "v",
													Idx:  0,
												},
											},
											Selectors: []src.Selector{},
										},
									},
								},
								{
									SenderSide: src.ConnectionSide{
										PortAddr: src.ConnPortAddr{
											Node: "print",
											RelPortAddr: src.RelPortAddr{
												Name: "v",
												Idx:  0,
											},
										},
										Selectors: []src.Selector{},
									},
									ReceiverSides: []src.ConnectionSide{
										{
											PortAddr: src.ConnPortAddr{
												Node: "trigger.",
												RelPortAddr: src.RelPortAddr{
													Name: "sig",
													Idx:  0,
												},
											},
											Selectors: []src.Selector{},
										},
									},
								},
								{
									SenderSide: src.ConnectionSide{
										PortAddr: src.ConnPortAddr{
											Node: "trigger",
											RelPortAddr: src.RelPortAddr{
												Name: "v",
												Idx:  0,
											},
										},
										Selectors: []src.Selector{},
									},
									ReceiverSides: []src.ConnectionSide{
										{
											PortAddr: src.ConnPortAddr{
												Node: "out.exit.",
												RelPortAddr: src.RelPortAddr{
													Name: "sig",
													Idx:  0,
												},
											},
											Selectors: []src.Selector{},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}

	irProg, err := irgen.Generator{}.Generate(nil, prog)
	if err != nil {
		panic(err)
	}

	bb, err := golang.Backend{}.GenerateTarget(nil, irProg)
	if err != nil {
		panic(err)
	}

	// write main.go
	var buf bytes.Buffer
	if _, err := buf.Write(bb); err != nil {
		panic(err)
	}
	if err := os.WriteFile(basePath+"/"+"main.go", buf.Bytes(), os.ModePerm); err != nil {
		panic(err)
	}
}

func putRuntime() {
	// prepare directory structure and collect files to create
	files := map[string][]byte{}
	if err := fs.WalkDir(efs, "runtime", func(path string, d fs.DirEntry, err error) error {
		fullPath := basePath + "/internal/" + path
		if d.IsDir() {
			if err := os.MkdirAll(fullPath, os.ModePerm); err != nil {
				return err
			}
			return nil
		}

		bb, err := efs.ReadFile(path)
		if err != nil {
			return err
		}

		files[fullPath] = bb
		return nil
	}); err != nil {
		panic(err)
	}
	// create files
	for path, bb := range files {
		if err := os.WriteFile(path, bb, os.ModePerm); err != nil {
			panic(err)
		}
	}
}

func putGoMod() {
	f, err := os.Create(basePath + "/go.mod")
	if err != nil {
		panic(err)
	}

	defer func() {
		if err := f.Close(); err != nil {
			panic(err)
		}
	}()

	_, err = f.WriteString("module github.com/emil14/neva")
	if err != nil {
		panic(err)
	}
}
