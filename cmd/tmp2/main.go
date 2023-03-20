package main

import (
	"bytes"
	"io/fs"
	"os"

	"github.com/emil14/neva/internal/compiler/backend/golang"
	"github.com/emil14/neva/internal/compiler/irgen"
	"github.com/emil14/neva/internal/compiler/src"
)

var efs = golang.Efs
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

	h := src.Helper{}

	prog := src.Program{
		Pkgs: map[string]src.Pkg{
			"main": {
				Imports: h.Imports("io", "flow"),
				Entities: map[string]src.Entity{
					"code": h.IntMsg(false, 0),
					"main": h.MainComponent(map[string]src.Node{
						"print": h.Node(
							h.NodeInstance("io", "Print", h.Inst("str")),
						),
						"trigger": h.NodeWithStaticPorts(
							h.NodeInstance("flow", "Trigger", h.Rec(nil)),
							map[src.RelPortAddr]src.EntityRef{
								{
									Name: "v",
									Idx:  0,
								}: {
									Pkg:  "",
									Name: "code",
								},
							},
						),
					}, []src.Connection{
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
					}),
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
