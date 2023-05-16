package main

import (
	"bytes"
	"io/fs"
	"os"

	"github.com/nevalang/neva/internal/compiler/backend/golang"
	"github.com/nevalang/neva/internal/compiler/ir"
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

	bb, err := (golang.Backend{}).GenerateTarget(nil, ir.Program{
		Ports: map[ir.PortAddr]uint8{
			{Name: "start"}:                            0,
			{Name: "exit"}:                             0,
			{Path: "printer.in", Name: "v", Idx: 0}:    0,
			{Path: "printer.out", Name: "v", Idx: 0}:   0,
			{Path: "trigger.in", Name: "sigs", Idx: 0}: 0,
			{Path: "trigger.in", Name: "v", Idx: 0}:    0,
			{Path: "trigger.out", Name: "v", Idx: 0}:   0,
			{Path: "giver.out", Name: "code", Idx: 0}:  0,
		},
		// Routines: ir.Routines{
		// 	Giver: map[ir.PortAddr]ir.Msg{
		// 		{
		// 			Path: "giver.out",
		// 			Name: "code",
		// 			Idx:  0,
		// 		}: {
		// 			Type: ir.IntMsg,
		// 			Int:  0,
		// 		},
		// 	},
		// 	Func: []ir.Func{
		// 		{
		// 			Ref: ir.FuncRef{
		// 				Pkg:  "io",
		// 				Name: "Print",
		// 			},
		// 			IO: ir.FuncIO{
		// 				In: []ir.PortAddr{
		// 					{Path: "printer.in", Name: "v", Idx: 0},
		// 				},
		// 				Out: []ir.PortAddr{
		// 					{Path: "printer.out", Name: "v", Idx: 0},
		// 				},
		// 			},
		// 		},
		// 		{
		// 			Ref: ir.FuncRef{
		// 				Pkg:  "flow",
		// 				Name: "Trigger",
		// 			},
		// 			IO: ir.FuncIO{
		// 				In: []ir.PortAddr{
		// 					{Path: "trigger.in", Name: "sigs", Idx: 0},
		// 					{Path: "trigger.in", Name: "v", Idx: 0},
		// 				},
		// 				Out: []ir.PortAddr{
		// 					{Path: "trigger.out", Name: "v", Idx: 0},
		// 				},
		// 			},
		// 		},
		// 	},
		// },
		Net: []ir.Connection{
			{
				SenderSide: ir.ConnectionSide{
					PortAddr: ir.PortAddr{
						Path: "",
						Name: "start",
						Idx:  0,
					},
					Selectors: []ir.Selector{},
				},
				ReceiverSides: []ir.ConnectionSide{
					{
						PortAddr: ir.PortAddr{
							Path: "printer.in",
							Name: "v",
							Idx:  0,
						},
						Selectors: []ir.Selector{},
					},
				},
			},
			{
				SenderSide: ir.ConnectionSide{
					PortAddr: ir.PortAddr{
						Path: "printer.out",
						Name: "v",
						Idx:  0,
					},
					Selectors: []ir.Selector{},
				},
				ReceiverSides: []ir.ConnectionSide{
					{
						PortAddr: ir.PortAddr{
							Path: "trigger.in",
							Name: "sigs",
							Idx:  0,
						},
						Selectors: []ir.Selector{},
					},
				},
			},
			{
				SenderSide: ir.ConnectionSide{
					PortAddr: ir.PortAddr{
						Path: "giver.out",
						Name: "code",
						Idx:  0,
					},
					Selectors: []ir.Selector{},
				},
				ReceiverSides: []ir.ConnectionSide{
					{
						PortAddr: ir.PortAddr{
							Path: "trigger.in",
							Name: "v",
							Idx:  0,
						},
						Selectors: []ir.Selector{},
					},
				},
			},
			{
				SenderSide: ir.ConnectionSide{
					PortAddr: ir.PortAddr{
						Path: "trigger.out",
						Name: "v",
						Idx:  0,
					},
					Selectors: []ir.Selector{},
				},
				ReceiverSides: []ir.ConnectionSide{
					{
						PortAddr: ir.PortAddr{
							Path: "",
							Name: "exit",
							Idx:  0,
						},
						Selectors: []ir.Selector{},
					},
				},
			},
		},
	})
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
		fullPath := basePath + "/internal/compiler/backend/golang/" + path
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

	_, err = f.WriteString("module github.com/nevalang/neva")
	if err != nil {
		panic(err)
	}
}
